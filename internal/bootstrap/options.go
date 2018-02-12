package bootstrap

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

const (
	WatcherTypeKube  = "kube"
	WatcherTypeFile  = "file"
	WatcherTypeVault = "vault"
)

var (
	supportedCwTypes = []string{
		WatcherTypeFile,
		WatcherTypeKube,
	}
	supportedSwTypes = []string{
		WatcherTypeVault,
		WatcherTypeKube,
	}
)

type Options struct {
	// these 3 get copied around. fun, i know
	kubeOptions          KubeOptions
	CommonWatcherOptions WatcherOptions

	ConfigWatcherOptions   WatcherOptions
	SecretWatcherOptions   WatcherOptions
	EndpointWatcherOptions WatcherOptions
	XdsOptions             XdsOptions
	Extra                  map[string]string
}

type WatcherOptions struct {
	Type          string
	SyncFrequency time.Duration
	FileOptions   FileOptions
	KubeOptions   KubeOptions
	VaultOptions  VaultOptions
}

type KubeOptions struct {
	KubeConfig string
	MasterURL  string
}

type VaultOptions struct {
	VaultAddr string
	AuthToken string
	Retries   int
}

type FileOptions struct {
	Path string
}

type XdsOptions struct {
	Port int
}

func (o *Options) InitFlags() {
	// core glue options
	flag.StringVar(&o.ConfigWatcherOptions.Type, "storage.type", WatcherTypeFile, fmt.Sprintf("storage backend for config objects. supported: [%s]", strings.Join(supportedCwTypes, " | ")))
	flag.StringVar(&o.SecretWatcherOptions.Type, "secrets.type", WatcherTypeFile, fmt.Sprintf("storage backend for secret objects. supported: [%s]", strings.Join(supportedSwTypes, " | ")))
	flag.IntVar(&o.XdsOptions.Port, "xds.port", 8081, "auth token for reading vault secrets")
	flag.IntVar(&o.CommonWatcherOptions.SyncFrequency, "refreshrate", WatcherTypeFile, fmt.Sprintf("storage backend for secret objects. supported: [%s]", strings.Join(supportedSwTypes, " | ")))

	// file
	flag.StringVar(&o.ConfigWatcherOptions.FileOptions.Path, "storage.file.dir", "_glue_config", "root directory to use for storing glue config files")
	flag.StringVar(&o.SecretWatcherOptions.FileOptions.Path, "secret.file.dir", "_glue_secrets", "root directory to use for storing glue secret files")

	// kube
	flag.StringVar(&o.kubeOptions.MasterURL, "kube.master", "", "url of the kubernetes apiserver. not needed if running in-cluster")
	flag.StringVar(&o.kubeOptions.KubeConfig, "kube.config", "", "path to kubeconfig file. not needed if running in-cluster")

	// vault
	flag.StringVar(&o.SecretWatcherOptions.VaultOptions.VaultAddr, "vault.addr", "", "url for vault server")
	flag.StringVar(&o.SecretWatcherOptions.VaultOptions.AuthToken, "vault.token", "", "auth token for reading vault secrets")
	flag.IntVar(&o.SecretWatcherOptions.VaultOptions.Retries, "vault.retries", 3, "number of times to retry failed requests to vault")
}

func (o *Options) ParseFlags() {
	flag.Parse()
	if o.ConfigWatcherOptions.Type == WatcherTypeKube {
		o.ConfigWatcherOptions.KubeOptions = o.kubeOptions
	}
	if o.SecretWatcherOptions.Type == WatcherTypeKube {
		o.SecretWatcherOptions.KubeOptions = o.kubeOptions
	}
}
