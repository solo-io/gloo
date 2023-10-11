package builders

import (
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	gloo_config_core_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extproc"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

const DefaultExtProcUpstreamName = "extproc-upstream"
const OverrideExtProcUpstreamName = "override-upstream"
const DefaultExtProcUpstreamNamespace = defaults.GlooSystem

func GetDefaultGrpcServiceBuilder() *GrpcServiceBuilder {
	return NewGrpcServiceBuilder().
		WithUpstreamName(DefaultExtProcUpstreamName).
		WithUpstreamNamespace(DefaultExtProcUpstreamNamespace).
		WithAuthority(&wrappers.StringValue{Value: "xyz"}).
		WithRetryPolicy(&gloo_config_core_v3.RetryPolicy{
			RetryBackOff: &gloo_config_core_v3.BackoffStrategy{
				BaseInterval: &duration.Duration{Seconds: 5},
				MaxInterval:  &duration.Duration{Seconds: 10},
			},
			NumRetries: &wrappers.UInt32Value{Value: 7},
		}).
		WithTimeout(&duration.Duration{Seconds: 100}).
		WithInitialMetadata([]*gloo_config_core_v3.HeaderValue{
			{Key: "A", Value: "B"},
			{Key: "C", Value: "D"},
		})
}

type GrpcServiceBuilder struct {
	upstreamName      string
	upstreamNamespace string
	authority         *wrappers.StringValue
	retryPolicy       *gloo_config_core_v3.RetryPolicy
	timeout           *duration.Duration
	initialMetadata   []*gloo_config_core_v3.HeaderValue
}

func NewGrpcServiceBuilder() *GrpcServiceBuilder {
	return &GrpcServiceBuilder{}
}

func (b *GrpcServiceBuilder) WithUpstreamName(name string) *GrpcServiceBuilder {
	b.upstreamName = name
	return b
}

func (b *GrpcServiceBuilder) WithUpstreamNamespace(namespace string) *GrpcServiceBuilder {
	b.upstreamNamespace = namespace
	return b
}

func (b *GrpcServiceBuilder) WithAuthority(authority *wrappers.StringValue) *GrpcServiceBuilder {
	b.authority = authority
	return b
}

func (b *GrpcServiceBuilder) WithRetryPolicy(retryPolicy *gloo_config_core_v3.RetryPolicy) *GrpcServiceBuilder {
	b.retryPolicy = retryPolicy
	return b
}

func (b *GrpcServiceBuilder) WithTimeout(timeout *duration.Duration) *GrpcServiceBuilder {
	b.timeout = timeout
	return b
}

func (b *GrpcServiceBuilder) WithInitialMetadata(initialMetadata []*gloo_config_core_v3.HeaderValue) *GrpcServiceBuilder {
	b.initialMetadata = initialMetadata
	return b
}

func (b *GrpcServiceBuilder) Build() *extproc.GrpcService {
	svc := &extproc.GrpcService{
		Authority:       b.authority,
		RetryPolicy:     b.retryPolicy,
		Timeout:         b.timeout,
		InitialMetadata: b.initialMetadata,
	}

	if b.upstreamName != "" || b.upstreamNamespace != "" {
		svc.ExtProcServerRef = &core.ResourceRef{Name: b.upstreamName, Namespace: b.upstreamNamespace}
	}

	return svc
}
