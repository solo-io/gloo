package caching

import (
	"fmt"

	envoycaching "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cache/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	envoycore "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	grpc "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/cache/grpc"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/caching"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"google.golang.org/protobuf/proto"
)

var (
	_ plugins.Plugin           = new(plugin)
	_ plugins.HttpFilterPlugin = new(plugin)
)

const (
	// FilterName is the name of the filter as registered in the filter_chain
	FilterName = "io.solo.filters.cache.grpc"
	// ExtensionName is the name of the extension as seen in gloo
	ExtensionName = "caching"

	defaultMaxSize = 65535 * 4 // default to the stream_window_size * 4
)

var (
	// currently set filter to be post authorization (in case its stale)
	// prior to ratelimit as we are not hitting the upstream if in cache
	filterStage = plugins.BeforeStage(plugins.RateLimitStage)
)

type plugin struct {
	serverSettings *caching.Settings
}

// NewPlugin instantiates the cache plugin
// It has no extra fields
func NewPlugin() *plugin {
	return &plugin{}
}

// Name returns the string name of the plugin
func (p *plugin) Name() string {
	return ExtensionName
}

// Init is used to re-initialize plugins and is executed for each translation loop
func (p *plugin) Init(params plugins.InitParams) {
	p.serverSettings = params.Settings.GetCachingServer()
}

// HttpFilters is called for every http listener.
func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {

	serverSettings := listener.GetOptions().GetCaching()

	if serverSettings == nil {
		serverSettings = p.serverSettings
	}

	upstreamRef := serverSettings.GetCachingServiceRef()
	if upstreamRef == nil {
		return []plugins.StagedHttpFilter{}, nil
	}

	_, err := params.Snapshot.Upstreams.Find(upstreamRef.GetNamespace(), upstreamRef.GetName())
	if err != nil {
		return []plugins.StagedHttpFilter{}, fmt.Errorf("caching server upstream not found %s", upstreamRef.String())
	}
	svc := &envoycore.GrpcService{TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
		EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
			ClusterName: translator.UpstreamToClusterName(upstreamRef),
		},
	}}

	maxSizePtr := serverSettings.GetMaxPayloadSize()
	maxSize := maxSizePtr.GetValue()
	if maxSizePtr == nil {
		maxSize = defaultMaxSize
	}

	grpcConfig := &grpc.GrpcCacheConfig{
		Service:        svc,
		MaxPayloadSize: maxSize, // set the max size as enforced by the grpc cache
	}

	typedConfig := utils.MustMessageToAny(grpcConfig)
	varyHeaders := serverSettings.GetAllowedVaryHeaders()

	allowedHeaders := []*v3.StringMatcher{}
	if varyHeaders != nil {
		allowedHeaders = make([]*v3.StringMatcher, len(varyHeaders))
		for i, header := range varyHeaders {
			allowedHeaders[i] = &v3.StringMatcher{}
			b, _ := proto.Marshal(header.ProtoReflect().Interface())
			proto.Unmarshal(b, allowedHeaders[i])
		}
	}

	cachingConfig := envoycaching.CacheConfig{
		AllowedVaryHeaders: allowedHeaders,
		TypedConfig:        typedConfig,
		// Max size is set at the custom layer
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, &cachingConfig, filterStage)
	if err != nil {
		return nil, err
	}
	return []plugins.StagedHttpFilter{stagedFilter}, nil
}
