package builders

import (
	"github.com/golang/protobuf/ptypes/duration"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	gloo_mutation_rules_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/common/mutation_rules/v3"
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
	grpcServiceBuilder     *GrpcServiceBuilder
	stage                  *filters.FilterStage
	failureModeAllow       *wrappers.BoolValue
	processingMode         *gloo_ext_proc_v3.ProcessingMode
	asyncMode              *wrappers.BoolValue
	statPrefix             *wrappers.StringValue
	mutationRules          *gloo_mutation_rules_v3.HeaderMutationRules
	disableClearRouteCache *wrappers.BoolValue
	forwardRules           *extproc.HeaderForwardingRules
	filterMetadata         *structpb.Struct
	allowModeOverride      *wrappers.BoolValue
	messageTimeout         *duration.Duration
	maxMessageTimeout      *duration.Duration
	requestAttributes      []string
	responseAttributes     []string
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

func (b *ExtProcBuilder) WithAsyncMode(asyncMode *wrappers.BoolValue) *ExtProcBuilder {
	b.asyncMode = asyncMode
	return b
}

func (b *ExtProcBuilder) WithStatPrefix(statPrefix *wrappers.StringValue) *ExtProcBuilder {
	b.statPrefix = statPrefix
	return b
}

func (b *ExtProcBuilder) WithMutationRules(mutationRules *gloo_mutation_rules_v3.HeaderMutationRules) *ExtProcBuilder {
	b.mutationRules = mutationRules
	return b
}

func (b *ExtProcBuilder) WithDisableClearRouteCache(disable *wrappers.BoolValue) *ExtProcBuilder {
	b.disableClearRouteCache = disable
	return b
}

func (b *ExtProcBuilder) WithForwardRules(rules *extproc.HeaderForwardingRules) *ExtProcBuilder {
	b.forwardRules = rules
	return b
}

func (b *ExtProcBuilder) WithFilterMetadata(md *structpb.Struct) *ExtProcBuilder {
	b.filterMetadata = md
	return b
}

func (b *ExtProcBuilder) WithAllowModeOverride(allow *wrappers.BoolValue) *ExtProcBuilder {
	b.allowModeOverride = allow
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
		FilterStage:            b.stage,
		FailureModeAllow:       b.failureModeAllow,
		ProcessingMode:         b.processingMode,
		AsyncMode:              b.asyncMode,
		MessageTimeout:         b.messageTimeout,
		MaxMessageTimeout:      b.maxMessageTimeout,
		RequestAttributes:      b.requestAttributes,
		ResponseAttributes:     b.responseAttributes,
		StatPrefix:             b.statPrefix,
		MutationRules:          b.mutationRules,
		DisableClearRouteCache: b.disableClearRouteCache,
		ForwardRules:           b.forwardRules,
		FilterMetadata:         b.filterMetadata,
		AllowModeOverride:      b.allowModeOverride,
	}

	if b.grpcServiceBuilder != nil {
		s.GrpcService = b.grpcServiceBuilder.Build()
	}

	return s
}
