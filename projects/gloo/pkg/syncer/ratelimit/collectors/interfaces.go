package collectors

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"go.uber.org/zap"
)

//go:generate mockgen -destination ./mocks/mock_interfaces.go -source ./interfaces.go

type CollectorType int

const (
	Global CollectorType = iota
	Basic
	Crd
)

type ConfigCollector interface {
	// Processes rate limit config on the given virtual host.
	ProcessVirtualHost(virtualHost *v1.VirtualHost, proxy *v1.Proxy, reports reporter.ResourceReports)

	// Processes rate limit config on the given route.
	ProcessRoute(route *v1.Route, virtualHost *v1.VirtualHost, proxy *v1.Proxy, reports reporter.ResourceReports)

	// Returns the current state of the collector as an xDS rate limit config.
	ToXdsConfiguration() (*enterprise.RateLimitConfig, error)
}

type ConfigCollectorFactory interface {
	MakeInstance(typ CollectorType, snapshot *v1snap.ApiSnapshot, logger *zap.SugaredLogger) (ConfigCollector, error)
}
