package graphql

import (
	"time"

	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	grpcResolverTypedExtensionConfigName = "io.solo.graphql.resolver.grpc"
	grpcRegistryExtensionName            = "grpc_extension"
)

func translateGrpcResolver(params plugins.RouteParams, r *v1alpha1.GrpcResolver) (*v3.TypedExtensionConfig, error) {
	requestTransform, err := translateGrpcRequestTransform(r.RequestTransform)
	if err != nil {
		return nil, err
	}
	us, err := params.Snapshot.Upstreams.Find(r.UpstreamRef.GetNamespace(), r.UpstreamRef.GetName())
	if err != nil {
		return nil, eris.Wrapf(err, "unable to find upstream `%s` in namespace `%s` to resolve schema", r.UpstreamRef.GetName(), r.UpstreamRef.GetNamespace())
	}
	grpcResolver := &v2.GrpcResolver{
		ServerUri: &v3.HttpUri{
			Uri: "ignored", // ignored by graphql filter
			HttpUpstreamType: &v3.HttpUri_Cluster{
				Cluster: translator.UpstreamToClusterName(us.GetMetadata().Ref()),
			},
			Timeout: durationpb.New(1 * time.Second),
		},
		RequestTransform: requestTransform,
		SpanName:         r.SpanName,
	}
	return &v3.TypedExtensionConfig{
		Name:        grpcResolverTypedExtensionConfigName,
		TypedConfig: utils.MustMessageToAny(grpcResolver),
	}, nil
}

func translateGrpcRequestTransform(transform *v1alpha1.GrpcRequestTemplate) (*v2.GrpcRequestTemplate, error) {
	if transform == nil {
		return nil, nil
	}
	rt := &v2.GrpcRequestTemplate{
		OutgoingMessageJson: nil, // filled in later
		ServiceName:         transform.ServiceName,
		MethodName:          transform.MethodName,
		RequestMetadata:     transform.RequestMetadata,
	}
	jv, err := TranslateJsonValue(transform.OutgoingMessageJson)
	if err != nil {
		return nil, err
	}
	rt.OutgoingMessageJson = jv
	return rt, nil
}
