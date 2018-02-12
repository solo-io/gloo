package bootstrap

import "time"

type WatcherType string

const (
	WatcherTypeKube  = "kube"
	WatcherTypeFile  = "file"
	WatcherTypeVault = "vault"
)

type Options struct {
	WatcherOptions WatcherOptions
	XdsOptions     XdsOptions
	Extra          map[string]string
}

type WatcherOptions struct {
	Type                     WatcherType
	SyncFrequency            time.Duration
	FileConfigWatcherOptions struct {
		Directory string
	}
	KubeOptions struct {
		KubeConfig string
		MasterURL  string
	}
	VaultOptions struct {
		VaultAddr string
		AuthToken string
		Retries   int
	}
}

type XdsOptions struct {
	Port int
}
