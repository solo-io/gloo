package installer

import (
	"context"
	"time"
)

type Installer interface {
	// GetContext returns the installation cluster and the installation namespace
	GetContext() (string, string)
	// Install attempts an installation, and returns an error if one occurred
	Install(ctx context.Context) error
	// Uninstall attempts an uninstallation, and returns an error if one occurred
	Uninstall(ctx context.Context) error
}

var _ Installer = new(NoOpInstaller)

type NoOpInstaller struct {
	duration time.Duration
	err      error
}

func NewNoOpInstaller(duration time.Duration, err error) *NoOpInstaller {
	return &NoOpInstaller{
		duration: duration,
		err:      err,
	}
}

func (n NoOpInstaller) GetContext() (string, string) {
	return "", ""
}

func (n NoOpInstaller) Install(_ context.Context) error {
	time.Sleep(n.duration)
	return n.err
}

func (n NoOpInstaller) Uninstall(_ context.Context) error {
	time.Sleep(n.duration)
	return n.err
}
