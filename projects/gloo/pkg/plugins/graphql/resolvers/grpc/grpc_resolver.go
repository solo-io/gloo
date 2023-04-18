package grpc

import (
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/types"

	resolver_utils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/utils"

	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

const (
	grpcResolverTypedExtensionConfigName = "io.solo.graphql.resolver.grpc"
	GrpcRegistryExtensionName            = "grpc_extension"
)

func TranslateGrpcResolver(upstreams types.UpstreamList, r *v1beta1.GrpcResolver) (*v3.TypedExtensionConfig, error) {
	requestTransform, err := translateGrpcRequestTransform(r.RequestTransform)
	if err != nil {
		return nil, err
	}
	us, err := upstreams.Find(r.UpstreamRef.GetNamespace(), r.UpstreamRef.GetName())
	if err != nil {
		return nil, eris.Wrapf(err, "unable to find upstream `%s` in namespace `%s` to resolve schema", r.UpstreamRef.GetName(), r.UpstreamRef.GetNamespace())
	}

	timeout := resolver_utils.DefaultTimeout

	if connTimeout := us.GetConnectionConfig().GetConnectTimeout(); connTimeout != nil {
		// This is a decent solution for initial config, but the GRPC resolver config will not update if the
		// upstream connection config changes
		timeout = connTimeout
	}
	// Use per-resolution timeout if we have one
	if resolverTimeout := r.GetTimeout(); resolverTimeout != nil {
		timeout = resolverTimeout
	}
	grpcResolver := &v2.GrpcResolver{
		ServerUri: &v3.HttpUri{
			Uri: "ignored", // ignored by graphql filter
			HttpUpstreamType: &v3.HttpUri_Cluster{
				Cluster: translator.UpstreamToClusterName(us.GetMetadata().Ref()),
			},
			Timeout: timeout,
		},
		RequestTransform: requestTransform,
		SpanName:         r.SpanName,
	}
	marshalledGRPCResolver, err := utils.MessageToAny(grpcResolver)
	if err != nil {
		return nil, eris.Wrapf(err, "unable to marshal grpcResolver")
	}
	return &v3.TypedExtensionConfig{
		Name:        grpcResolverTypedExtensionConfigName,
		TypedConfig: marshalledGRPCResolver,
	}, nil
}

func translateGrpcRequestTransform(transform *v1beta1.GrpcRequestTemplate) (*v2.GrpcRequestTemplate, error) {
	if transform == nil {
		return nil, nil
	}
	rt := &v2.GrpcRequestTemplate{
		OutgoingMessageJson: nil, // filled in later
		ServiceName:         transform.ServiceName,
		MethodName:          transform.MethodName,
		RequestMetadata:     transform.RequestMetadata,
	}
	jv, err := resolver_utils.TranslateJsonValue(transform.OutgoingMessageJson)
	if err != nil {
		return nil, err
	}
	rt.OutgoingMessageJson = jv
	return rt, nil
}
