package aws

import (
	"github.com/golang/protobuf/ptypes/any"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	envoyaws "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	ossplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/aws"
)

var (
	_ plugins.Plugin           = new(plugin)
	_ plugins.UpstreamPlugin   = new(plugin)
	_ plugins.RoutePlugin      = new(plugin)
	_ plugins.HttpFilterPlugin = new(plugin)
)

const (
	FilterName    = "io.solo.api_gateway.api_gateway_transformer"
	FilterTypeUrl = "type.googleapis.com/envoy.config.transformer.aws_lambda.v2.ApiGatewayTransformation"
)

type plugin struct {
	ossplugin.Plugin
}

func NewPlugin() plugins.Plugin {
	// The Enterprise Plugin is a replica of the Open Source Plugin,
	// except that it builds the PerRoute configuration differently
	return ossplugin.NewPlugin(GenerateEnterpriseAWSLambdaRouteConfig)
}

func GenerateEnterpriseAWSLambdaRouteConfig(destination *aws.DestinationSpec, upstream *aws.UpstreamSpec) (*envoyaws.AWSLambdaPerRoute, error) {
	lambdaPerRoute, err := ossplugin.GenerateAWSLambdaRouteConfig(destination, upstream)
	if err != nil {
		return lambdaPerRoute, err
	}

	var transformerConfig *v3.TypedExtensionConfig

	if destination.GetUnwrapAsApiGateway() && !destination.GetUnwrapAsAlb() {
		transformerConfig = &v3.TypedExtensionConfig{
			Name: FilterName,
			TypedConfig: &any.Any{
				TypeUrl: FilterTypeUrl,
			},
		}
	}

	lambdaPerRoute.TransformerConfig = transformerConfig
	return lambdaPerRoute, nil
}
