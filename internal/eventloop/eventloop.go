package eventloop

import (
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/pkg/errors"
	"github.com/solo-io/glue/internal/plugins"
	"github.com/solo-io/glue/internal/reporter"
	filesecrets "github.com/solo-io/glue/internal/secretwatcher/file"
	kubesecrets "github.com/solo-io/glue/internal/secretwatcher/kube"
	"github.com/solo-io/glue/internal/secretwatcher/vault"
	"github.com/solo-io/glue/pkg/plugin"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/solo-io/glue/internal/bootstrap"
	"github.com/solo-io/glue/internal/configwatcher/file"
	"github.com/solo-io/glue/internal/configwatcher/kube"
	"github.com/solo-io/glue/internal/translator"
	"github.com/solo-io/glue/internal/xds"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/configwatcher"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

type eventLoop struct {
	configWatcher       configwatcher.Interface
	secretWatcher       secretwatcher.Interface
	endpointDiscoveries []endpointdiscovery.Interface
	//TODO: reporter
	reporter         reporter.Interface
	translator       *translator.Translator
	xdsConfig        envoycache.Cache
	updateSecretRefs func(cfg *v1.Config) []string

	startFuncs []func() error
}

func Setup(opts bootstrap.Options, stopCh <-chan struct{}) (*eventLoop, error) {
	cfgWatcher, err := setupConfigWatcher(opts, stopCh)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up config watcher")
	}

	secretWatcher, err := setupSecretWatcher(opts, stopCh)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up secret watcher")
	}

	xdsConfig, _, err := xds.RunXDS(opts.XdsOptions.Port)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start xds server")
	}

	plugs := plugins.RegisteredPlugins()

	trans := translator.NewTranslator(plugs)

	e := &eventLoop{
		configWatcher:    cfgWatcher,
		secretWatcher:    secretWatcher,
		translator:       trans,
		xdsConfig:        xdsConfig,
		updateSecretRefs: updateSecretRefsFor(plugs),
	}

	for _, endpointDiscoveryInitializer := range plugins.EndpointDiscoveryInitializers() {
		e.startFuncs = append(e.startFuncs, func() error {
			discovery, err := endpointDiscoveryInitializer(opts, stopCh)
			if err != nil {
				return err
			}
			e.endpointDiscoveries = append(e.endpointDiscoveries, discovery)
			return nil
		})
	}
	return e, nil
}

func updateSecretRefsFor(plugins []plugin.TranslatorPlugin) func(cfg *v1.Config) []string {
	return func(cfg *v1.Config) []string {
		var secretRefs []string
		for _, plug := range plugins {
			deps := plug.GetDependencies(cfg)
			if deps != nil {
				secretRefs = append(secretRefs, deps.SecretRefs...)
			}
		}
		return secretRefs
	}
}

func setupConfigWatcher(opts bootstrap.Options, stopCh <-chan struct{}) (configwatcher.Interface, error) {
	switch opts.ConfigWatcherOptions.Type {
	case bootstrap.WatcherTypeFile:
		dir := opts.FileOptions.ConfigDir
		if dir == "" {
			return nil, errors.New("must provide directory for file config watcher")
		}
		cfgWatcher, err := file.NewFileConfigWatcher(dir, opts.ConfigWatcherOptions.SyncFrequency)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start file config watcher for directory %v", dir)
		}
		return cfgWatcher, nil
	case bootstrap.WatcherTypeKube:
		cfgWatcher, err := kube.NewCrdWatcher(opts.KubeOptions.MasterURL, opts.KubeOptions.KubeConfig, opts.ConfigWatcherOptions.SyncFrequency, stopCh)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start kube config watcher with config %#v", opts.KubeOptions)
		}
		return cfgWatcher, nil
	}
	return nil, errors.Errorf("unknown or unspecified config watcher type: %v", opts.ConfigWatcherOptions.Type)
}

func setupSecretWatcher(opts bootstrap.Options, stopCh <-chan struct{}) (secretwatcher.Interface, error) {
	switch opts.SecretWatcherOptions.Type {
	case bootstrap.WatcherTypeFile:
		secretWatcher, err := filesecrets.NewSecretWatcher(opts.FileOptions.SecretDir, opts.SecretWatcherOptions.SyncFrequency)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start file secret watcher with config %#v", opts.KubeOptions)
		}
		return secretWatcher, nil
	case bootstrap.WatcherTypeKube:
		secretWatcher, err := kubesecrets.NewSecretWatcher(opts.KubeOptions.MasterURL, opts.KubeOptions.KubeConfig, opts.SecretWatcherOptions.SyncFrequency, stopCh)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start kube secret watcher with config %#v", opts.KubeOptions)
		}
		return secretWatcher, nil
	case bootstrap.WatcherTypeVault:
		secretWatcher, err := vault.NewVaultSecretWatcher(opts.SecretWatcherOptions.SyncFrequency, opts.VaultOptions.Retries, opts.VaultOptions.VaultAddr, opts.VaultOptions.AuthToken, stopCh)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start vault secret watcher with config %#v", opts.VaultOptions)
		}
		return secretWatcher, nil
	}
	return nil, errors.Errorf("unknown or unspecified secret watcher type: %v", opts.SecretWatcherOptions.Type)
}

func (e *eventLoop) Run() error {
	for _, fn := range e.startFuncs {
		if err := fn(); err != nil {
			return err
		}
	}

	endpointDiscovery := e.endpointDiscovery()
	workerErrors := e.errors()

	// cache the most recent read for any of these
	current := newCache()
	for {
		select {
		case cfg := <-e.configWatcher.Config():
			current.cfg = cfg
			secretRefs := e.updateSecretRefs(cfg)
			go e.secretWatcher.TrackSecrets(secretRefs)
			for _, discovery := range e.endpointDiscoveries {
				go func() {
					discovery.TrackUpstreams(cfg.Upstreams)
				}()
			}
			e.updateXds(current)
		case secrets := <-e.secretWatcher.Secrets():
			current.secrets = secrets
			e.updateXds(current)
		case endpointTuple := <-endpointDiscovery:
			current.endpoints[endpointTuple.discoveredBy] = endpointTuple.endpoints
			e.updateXds(current)
		case err := <-workerErrors:
			runtime.HandleError(err)
		}
	}
}

func (e *eventLoop) updateXds(cache *cache) {
	if !cache.ready() {
		log.Debugf("cache is not fully constructed to produce a first snapshot yet")
		return
	}
	aggregatedEndpoints := make(endpointdiscovery.EndpointGroups)
	for _, endpointGroups := range cache.endpoints {
		for upstreamName, endpointSet := range endpointGroups {
			aggregatedEndpoints[upstreamName] = endpointSet
		}
	}
	snapshot, status, err := e.translator.Translate(cache.cfg, cache.secrets, aggregatedEndpoints)
	if err != nil {
		runtime.HandleError(errors.Wrap(err, "failed to translate based on the latest config"))
	}
	log.Printf("TODO: do something with this status eventually: %v", status)
	log.Debugf("FINAL: XDS Snapshot: %v", snapshot)
	e.xdsConfig.SetSnapshot(xds.NodeKey, *snapshot)
}

// fan out to cover all endpoint discovery services
func (e *eventLoop) endpointDiscovery() <-chan endpointTuple {
	aggregatedEndpointsChan := make(chan endpointTuple)
	for _, ed := range e.endpointDiscoveries {
		go func() {
			for endpoints := range ed.Endpoints() {
				aggregatedEndpointsChan <- endpointTuple{
					endpoints:    endpoints,
					discoveredBy: ed,
				}
			}
		}()
	}
	return aggregatedEndpointsChan
}

// fan out to cover all channels that return errors
func (e *eventLoop) errors() <-chan error {
	aggregatedErrorsChan := make(chan error)
	go func() {
		for err := range e.configWatcher.Error() {
			aggregatedErrorsChan <- errors.Wrap(err, "config watcher encountered an error")
		}
	}()
	go func() {
		for err := range e.secretWatcher.Error() {
			aggregatedErrorsChan <- errors.Wrap(err, "secret watcher encountered an error")
		}
	}()
	for _, ed := range e.endpointDiscoveries {
		go func() {
			for err := range ed.Error() {
				aggregatedErrorsChan <- err
			}
		}()
	}
	return aggregatedErrorsChan
}

// cache contains the latest "glue snapshot"
type cache struct {
	cfg     *v1.Config
	secrets secretwatcher.SecretMap
	// need to separate endpoints by the service who discovered them
	endpoints map[endpointdiscovery.Interface]endpointdiscovery.EndpointGroups
}

func newCache() *cache {
	return &cache{
		endpoints: make(map[endpointdiscovery.Interface]endpointdiscovery.EndpointGroups),
	}
}

// ready doesn't necessarily tell us whetehr endpoints have been discovered yet
// but that's okay. envoy won't mind
func (c *cache) ready() bool {
	return c.cfg != nil
}

type endpointTuple struct {
	discoveredBy endpointdiscovery.Interface
	endpoints    endpointdiscovery.EndpointGroups
}
