package eventloop

import (
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/internal/swagger"
	"github.com/solo-io/gloo-function-discovery/internal/updater"
	"github.com/solo-io/gloo-function-discovery/internal/upstreamwatcher"
	"github.com/solo-io/gloo-function-discovery/pkg/resolver"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	"github.com/solo-io/gloo-storage/file"
	"github.com/solo-io/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/secretwatcher"
	filesecrets "github.com/solo-io/gloo/pkg/secretwatcher/file"
	kubesecrets "github.com/solo-io/gloo/pkg/secretwatcher/kube"
	"github.com/solo-io/gloo/pkg/secretwatcher/vault"
	"k8s.io/client-go/kubernetes"
)

func Run(opts bootstrap.Options, autoDiscoverSwagger bool, swaggerUrisToTry []string, stop <-chan struct{}, errs chan error) error {
	store, err := createStorageClient(opts)
	if err != nil {
		return errors.Wrap(err, "failed to create config store client")
	}

	upstreams, err := upstreamwatcher.WatchUpstreams(store, stop, errs)
	if err != nil {
		return errors.Wrap(err, "failed to start monitoring upstreams")
	}

	secretWatcher, err := setupSecretWatcher(opts, stop)
	if err != nil {
		return errors.Wrap(err, "failed to set up secret watcher")
	}

	resolver := createResolver(opts)

	var cache struct {
		secrets   secretwatcher.SecretMap
		upstreams []*v1.Upstream
	}

	update := func() {
		go func() {
			// update secret refs on secret watcher
			refs := updater.GetSecretRefsToWatch(cache.upstreams)
			secretWatcher.TrackSecrets(refs)
		}()
		if autoDiscoverSwagger {
			swagger.DiscoverSwaggerUpstreams(resolver, swaggerUrisToTry, cache.upstreams)
		}
		if err := updater.UpdateFunctionalUpstreams(store, cache.upstreams, cache.secrets); err != nil {
			errs <- err
		}
	}

	tick := time.Tick(opts.ConfigWatcherOptions.SyncFrequency)

	for {
		select {
		case cache.secrets = <-secretWatcher.Secrets():
			update()
		case cache.upstreams = <-upstreams:
			update()
		case <-tick:
			update()
		case err := <-secretWatcher.Error():
			errs <- err
		case <-stop:
			return nil
		}
	}
}

func createStorageClient(opts bootstrap.Options) (storage.Interface, error) {
	switch opts.ConfigWatcherOptions.Type {
	case bootstrap.WatcherTypeFile:
		dir := opts.FileOptions.ConfigDir
		if dir == "" {
			return nil, errors.New("must provide directory for file config watcher")
		}
		client, err := file.NewStorage(dir, opts.ConfigWatcherOptions.SyncFrequency)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start file config watcher for directory %v", dir)
		}
		return client, nil
	case bootstrap.WatcherTypeKube:
		cfg, err := clientcmd.BuildConfigFromFlags(opts.KubeOptions.MasterURL, opts.KubeOptions.KubeConfig)
		if err != nil {
			return nil, errors.Wrap(err, "building kube restclient")
		}
		cfgWatcher, err := crd.NewStorage(cfg, opts.KubeOptions.Namespace, opts.ConfigWatcherOptions.SyncFrequency)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start kube config watcher with config %#v", opts.KubeOptions)
		}
		return cfgWatcher, nil
	}
	return nil, errors.Errorf("unknown or unspecified config watcher type: %v", opts.ConfigWatcherOptions.Type)
}

func createResolver(opts bootstrap.Options) *resolver.Resolver {
	kube, err := func() (kubernetes.Interface, error) {
		cfg, err := clientcmd.BuildConfigFromFlags(opts.KubeOptions.MasterURL, opts.KubeOptions.KubeConfig)
		if err != nil {
			return nil, err
		}
		kube, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, err
		}
		return kube, nil
	}()
	if err != nil {
		log.Warnf("create kube client failed: %v. swagger services running in kubernetes will not be discovered by function discovery")
	}
	return &resolver.Resolver{Kube: kube}
}

func setupSecretWatcher(opts bootstrap.Options, stop <-chan struct{}) (secretwatcher.Interface, error) {
	switch opts.SecretWatcherOptions.Type {
	case bootstrap.WatcherTypeFile:
		secretWatcher, err := filesecrets.NewSecretWatcher(opts.FileOptions.SecretDir, opts.SecretWatcherOptions.SyncFrequency)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start file secret watcher with config %#v", opts.KubeOptions)
		}
		return secretWatcher, nil
	case bootstrap.WatcherTypeKube:
		secretWatcher, err := kubesecrets.NewSecretWatcher(opts.KubeOptions.MasterURL, opts.KubeOptions.KubeConfig, opts.SecretWatcherOptions.SyncFrequency, stop)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start kube secret watcher with config %#v", opts.KubeOptions)
		}
		return secretWatcher, nil
	case bootstrap.WatcherTypeVault:
		secretWatcher, err := vault.NewVaultSecretWatcher(opts.SecretWatcherOptions.SyncFrequency, opts.VaultOptions.Retries, opts.VaultOptions.VaultAddr, opts.VaultOptions.AuthToken, stop)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start vault secret watcher with config %#v", opts.VaultOptions)
		}
		return secretWatcher, nil
	}
	return nil, errors.Errorf("unknown or unspecified secret watcher type: %v", opts.SecretWatcherOptions.Type)
}
