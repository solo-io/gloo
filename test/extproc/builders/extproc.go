package builders

import (
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	gloo_ext_proc_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/ext_proc/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extproc"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/filters"
)

func GetDefaultExtProcBuilder() *ExtProcBuilder {
	return NewExtProcBuilder().
		WithGrpcServiceBuilder(GetDefaultGrpcServiceBuilder()).
		WithStage(&filters.FilterStage{Stage: filters.FilterStage_AcceptedStage, Predicate: filters.FilterStage_Before}).
		WithFailureModeAllow(&wrappers.BoolValue{Value: true}).
		WithProcessingMode(&gloo_ext_proc_v3.ProcessingMode{
			RequestHeaderMode:   gloo_ext_proc_v3.ProcessingMode_SEND,
			ResponseHeaderMode:  gloo_ext_proc_v3.ProcessingMode_SEND,
			RequestBodyMode:     gloo_ext_proc_v3.ProcessingMode_BUFFERED,
			ResponseBodyMode:    gloo_ext_proc_v3.ProcessingMode_BUFFERED_PARTIAL,
			RequestTrailerMode:  gloo_ext_proc_v3.ProcessingMode_SKIP,
			ResponseTrailerMode: gloo_ext_proc_v3.ProcessingMode_DEFAULT,
		}).
		WithMessageTimeout(&duration.Duration{Seconds: 1}).
		WithMaxMessageTimeout(&duration.Duration{Seconds: 5}).
		WithRequestAttributes([]string{"req1", "req2"}).
		WithResponseAttributes([]string{"resp1", "resp2", "resp3"})
}

type ExtProcBuilder struct {
	grpcServiceBuilder *GrpcServiceBuilder
	stage              *filters.FilterStage
	failureModeAllow   *wrappers.BoolValue
	processingMode     *gloo_ext_proc_v3.ProcessingMode
	messageTimeout     *duration.Duration
	maxMessageTimeout  *duration.Duration
	requestAttributes  []string
	responseAttributes []string
}

func NewExtProcBuilder() *ExtProcBuilder {
	return &ExtProcBuilder{}
}

func (b *ExtProcBuilder) WithGrpcServiceBuilder(builder *GrpcServiceBuilder) *ExtProcBuilder {
	b.grpcServiceBuilder = builder
	return b
}

func (b *ExtProcBuilder) WithStage(stage *filters.FilterStage) *ExtProcBuilder {
	b.stage = stage
	return b
}

func (b *ExtProcBuilder) WithFailureModeAllow(allow *wrappers.BoolValue) *ExtProcBuilder {
	b.failureModeAllow = allow
	return b
}

func (b *ExtProcBuilder) WithProcessingMode(mode *gloo_ext_proc_v3.ProcessingMode) *ExtProcBuilder {
	b.processingMode = mode
	return b
}

func (b *ExtProcBuilder) WithMessageTimeout(timeout *duration.Duration) *ExtProcBuilder {
	b.messageTimeout = timeout
	return b
}

func (b *ExtProcBuilder) WithMaxMessageTimeout(timeout *duration.Duration) *ExtProcBuilder {
	b.maxMessageTimeout = timeout
	return b
}

func (b *ExtProcBuilder) WithRequestAttributes(attr []string) *ExtProcBuilder {
	b.requestAttributes = attr
	return b
}

func (b *ExtProcBuilder) WithResponseAttributes(attr []string) *ExtProcBuilder {
	b.responseAttributes = attr
	return b
}

func (b *ExtProcBuilder) Build() *extproc.Settings {
	s := &extproc.Settings{
		FilterStage:        b.stage,
		FailureModeAllow:   b.failureModeAllow,
		ProcessingMode:     b.processingMode,
		MessageTimeout:     b.messageTimeout,
		MaxMessageTimeout:  b.maxMessageTimeout,
		RequestAttributes:  b.requestAttributes,
		ResponseAttributes: b.responseAttributes,
		// other fields not currently tested
		// AsyncMode:              nil,
		// StatPrefix:             nil,
		// MutationRules:          nil,
		// DisableClearRouteCache: nil,
		// ForwardRules:           nil,
		// FilterMetadata:         nil,
		// AllowModeOverride:      nil,
	}

	if b.grpcServiceBuilder != nil {
		s.GrpcService = b.grpcServiceBuilder.Build()
	}

	return s
}
