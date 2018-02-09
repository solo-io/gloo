package eventloop

import (
	"time"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/pkg/errors"

	"github.com/solo-io/glue/internal/configwatcher/file"
	"github.com/solo-io/glue/internal/configwatcher/kube"
	translator2 "github.com/solo-io/glue/internal/translator"
	"github.com/solo-io/glue/internal/xds"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/configwatcher"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/pkg/plugin2"
	"github.com/solo-io/glue/pkg/secretwatcher"
	"github.com/solo-io/glue/pkg/signals"
	"github.com/solo-io/glue/pkg/translator"
)

type ConfigWatcherType string

const (
	ConfigWatcherType_Kube = "kube"
	ConfigWatcherType_File = "file"
)

type Options struct {
	ConfigWatcherOptions ConfigWatcherOptions
}

type ConfigWatcherOptions struct {
	Type                     ConfigWatcherType
	SyncFrequency            time.Duration
	FileConfigWatcherOptions struct {
		Directory string
	}
	KubeConfigWatcherOptions struct {
		// if these are not set, will attempt to use in cluster config
		KubeConfig string
		MasterURL  string
	}
}

type eventLoop struct {
	configWatcher     configwatcher.Interface
	endpointDiscovery []endpointdiscovery.Interface
	secretWatcher     secretwatcher.Interface
	translator        translator.Interface
	xdsConfig         envoycache.Cache
	plugins           []plugin.TranslatorPlugin

	startFuncs []func() error
}

func Setup(opts *Options) (*eventLoop, error) {
	gracefulStop := signals.SetupSignalHandler()
	cfgWatcher, err := setupConfigWatcher(opts.ConfigWatcherOptions, gracefulStop)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up config watcher")
	}
	trans := translator2.NewTranslator()
}

func setupConfigWatcher(opts ConfigWatcherOptions, stopCh <-chan struct{}) (configwatcher.Interface, error) {
	switch opts.Type {
	case ConfigWatcherType_File:
		dir := opts.FileConfigWatcherOptions.Directory
		if dir == "" {
			return nil, errors.New("must provide directory for file config watcher")
		}
		cfgWatcher, err := file.NewFileConfigWatcher(dir, opts.SyncFrequency)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start file config watcher for directory %v", dir)
		}
		return cfgWatcher, nil
	case ConfigWatcherType_Kube:
		cfgWatcher, err := kube.NewCrdWatcher(opts.KubeConfigWatcherOptions.MasterURL, opts.KubeConfigWatcherOptions.KubeConfig, opts.SyncFrequency, stopCh)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to start kube config watcher with config %#v", opts.KubeConfigWatcherOptions)
		}
		return cfgWatcher, nil
	}
	return nil, errors.Errorf("unknown or unspecified config watcher type: %v", opts.Type)
}

func Run() error {

	for {
		select {
		case configWatcher.Config():

		}
	}
}

func updateXds(translator translator.Interface, xdsConfig envoycache.Cache,
	cfg *v1.Config, secretMap secretwatcher.SecretMap, endpoints endpointdiscovery.EndpointGroups) error {
	snapshot, status := translator.Translate(cfg, secretMap, endpoints)
	log.Printf("TODO: do something with this status eventually: %v", status)
	xdsConfig.SetSnapshot(xds.NodeKey, *snapshot)
	return nil
}
