package builders

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	gloo_config_core_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	gloo_ext_proc_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/ext_proc/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extproc"
)

func GetDefaultExtProcRouteBuilder() *ExtProcRouteBuilder {
	return NewExtProcRouteBuilder().
		WithGrpcServiceBuilder(NewGrpcServiceBuilder().
			WithUpstreamName(OverrideExtProcUpstreamName).
			WithUpstreamNamespace(DefaultExtProcUpstreamNamespace).
			WithInitialMetadata([]*gloo_config_core_v3.HeaderValue{{Key: "aaa", Value: "bbb"}}),
		).
		WithProcessingMode(&gloo_ext_proc_v3.ProcessingMode{
			RequestHeaderMode:  gloo_ext_proc_v3.ProcessingMode_SEND,
			ResponseHeaderMode: gloo_ext_proc_v3.ProcessingMode_SKIP,
			RequestBodyMode:    gloo_ext_proc_v3.ProcessingMode_STREAMED,
		}).
		WithAsyncMode(&wrappers.BoolValue{Value: true}).
		WithRequestAttributes([]string{"x"}).
		WithResponseAttributes([]string{"y"})
}

type ExtProcRouteBuilder struct {
	disabled *wrappers.BoolValue
	// overrides
	processingMode     *gloo_ext_proc_v3.ProcessingMode
	asyncMode          *wrappers.BoolValue
	requestAttributes  []string
	responseAttributes []string
	grpcServiceBuilder *GrpcServiceBuilder
}

func NewExtProcRouteBuilder() *ExtProcRouteBuilder {
	return &ExtProcRouteBuilder{}
}

// WithDisabled sets the disabled value to the given BoolValue, and sets all other fields to nil.
// The envoy ExtProcPerRoute only allows one of `disabled` or `overrides` to be set, and this is
// just a convenience so we don't leave unused fields set on the builder.
// Likewise, setting any of the overrides using the With* functions below will set `disabled` to nil.
func (b *ExtProcRouteBuilder) WithDisabled(disabled *wrappers.BoolValue) *ExtProcRouteBuilder {
	b.disabled = disabled

	// set everything else to nil
	b.processingMode = nil
	b.asyncMode = nil
	b.requestAttributes = nil
	b.responseAttributes = nil
	b.grpcServiceBuilder = nil

	return b
}

func (b *ExtProcRouteBuilder) WithProcessingMode(mode *gloo_ext_proc_v3.ProcessingMode) *ExtProcRouteBuilder {
	b.processingMode = mode
	b.disabled = nil
	return b
}

func (b *ExtProcRouteBuilder) WithAsyncMode(asyncMode *wrappers.BoolValue) *ExtProcRouteBuilder {
	b.asyncMode = asyncMode
	b.disabled = nil
	return b
}

func (b *ExtProcRouteBuilder) WithRequestAttributes(attr []string) *ExtProcRouteBuilder {
	b.requestAttributes = attr
	b.disabled = nil
	return b
}

func (b *ExtProcRouteBuilder) WithResponseAttributes(attr []string) *ExtProcRouteBuilder {
	b.responseAttributes = attr
	b.disabled = nil
	return b
}

func (b *ExtProcRouteBuilder) WithGrpcServiceBuilder(builder *GrpcServiceBuilder) *ExtProcRouteBuilder {
	b.grpcServiceBuilder = builder
	b.disabled = nil
	return b
}

func (b *ExtProcRouteBuilder) Build() *extproc.RouteSettings {
	if b.disabled != nil {
		return &extproc.RouteSettings{
			Override: &extproc.RouteSettings_Disabled{
				Disabled: b.disabled,
			},
		}
	}

	s := &extproc.RouteSettings{
		Override: &extproc.RouteSettings_Overrides{
			Overrides: &extproc.Overrides{
				ProcessingMode:     b.processingMode,
				AsyncMode:          b.asyncMode,
				RequestAttributes:  b.requestAttributes,
				ResponseAttributes: b.responseAttributes,
			},
		},
	}

	if b.grpcServiceBuilder != nil {
		s.Override.(*extproc.RouteSettings_Overrides).Overrides.GrpcService = b.grpcServiceBuilder.Build()
	}

	return s
}
