package bootstrap

import (
	"time"
)

const (
	WatcherTypeKube  = "kube"
	WatcherTypeFile  = "file"
	WatcherTypeVault = "vault"
)

var (
	SupportedCwTypes = []string{
		WatcherTypeFile,
		WatcherTypeKube,
	}
	SupportedFwTypes = []string{
		WatcherTypeFile,
		WatcherTypeKube,
	}
	SupportedSwTypes = []string{
		WatcherTypeVault,
		WatcherTypeKube,
	}
)

type Options struct {
	// these 3 get copied around. fun, i know
	KubeOptions          KubeOptions
	ConfigWatcherOptions WatcherOptions
	SecretWatcherOptions WatcherOptions
	FileWatcherOptions   WatcherOptions
	FileOptions          FileOptions
	VaultOptions         VaultOptions
	XdsOptions           XdsOptions
	// may be needed by plugins
	Extra map[string]string
}

type WatcherOptions struct {
	Type          string
	SyncFrequency time.Duration
}

type KubeOptions struct {
	KubeConfig string
	MasterURL  string
	Namespace  string // where to watch for storage
}

type VaultOptions struct {
	VaultAddr string
	AuthToken string
	Retries   int
}

type FileOptions struct {
	ConfigDir string
	SecretDir string
	FilesDir  string
}

type XdsOptions struct {
	Port int
}
