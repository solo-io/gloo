package mutation

import (
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

//go:generate mockgen -destination mocks/mutation_factory_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/mutation Factory

var (
	EmptyInputError = errors.Errorf("Upstream input is empty")

	UnsupportedUpstreamTypeError = errors.Errorf("The provided upstream type is not yet supported")
)

type Factory interface {
	ConfigureUpstream(input *v1.UpstreamInput) Mutation
}

type factory struct{}

var _ Factory = factory{}

func NewFactory() factory {
	return factory{}
}

func (factory) ConfigureUpstream(input *v1.UpstreamInput) Mutation {
	return func(upstream *gloov1.Upstream) error {
		if input.GetSpec() == nil {
			return EmptyInputError
		}

		if upstream.GetUpstreamSpec() == nil {
			upstream.UpstreamSpec = &gloov1.UpstreamSpec{}
		}

		switch input.GetSpec().(type) {
		case *v1.UpstreamInput_Aws:
			upstream.UpstreamSpec.UpstreamType = &gloov1.UpstreamSpec_Aws{Aws: input.GetAws()}
		case *v1.UpstreamInput_Azure:
			upstream.UpstreamSpec.UpstreamType = &gloov1.UpstreamSpec_Azure{Azure: input.GetAzure()}
		case *v1.UpstreamInput_Consul:
			upstream.UpstreamSpec.UpstreamType = &gloov1.UpstreamSpec_Consul{Consul: input.GetConsul()}
		case *v1.UpstreamInput_Kube:
			upstream.UpstreamSpec.UpstreamType = &gloov1.UpstreamSpec_Kube{Kube: input.GetKube()}
		case *v1.UpstreamInput_Static:
			upstream.UpstreamSpec.UpstreamType = &gloov1.UpstreamSpec_Static{Static: input.GetStatic()}
		case *v1.UpstreamInput_AwsEc2:
			upstream.UpstreamSpec.UpstreamType = &gloov1.UpstreamSpec_AwsEc2{AwsEc2: input.GetAwsEc2()}
		default:
			return UnsupportedUpstreamTypeError
		}
		upstream.Status = core.Status{}
		return nil
	}
}
