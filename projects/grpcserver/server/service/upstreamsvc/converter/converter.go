package converter

import (
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

type UpstreamInputConverter interface {
	ConvertInputToUpstreamSpec(input *v1.UpstreamInput) *gloov1.UpstreamSpec
}

type converter struct{}

func NewUpstreamInputConverter() UpstreamInputConverter {
	return &converter{}
}

func (*converter) ConvertInputToUpstreamSpec(input *v1.UpstreamInput) *gloov1.UpstreamSpec {
	switch input.GetSpec().(type) {
	case *v1.UpstreamInput_Aws:
		return &gloov1.UpstreamSpec{
			UpstreamType: &gloov1.UpstreamSpec_Aws{
				Aws: input.GetAws(),
			},
		}
	case *v1.UpstreamInput_Azure:
		return &gloov1.UpstreamSpec{
			UpstreamType: &gloov1.UpstreamSpec_Azure{
				Azure: input.GetAzure(),
			},
		}
	case *v1.UpstreamInput_Consul:
		return &gloov1.UpstreamSpec{
			UpstreamType: &gloov1.UpstreamSpec_Consul{
				Consul: input.GetConsul(),
			},
		}
	case *v1.UpstreamInput_Kube:
		return &gloov1.UpstreamSpec{
			UpstreamType: &gloov1.UpstreamSpec_Kube{
				Kube: input.GetKube(),
			},
		}
	case *v1.UpstreamInput_Static:
		return &gloov1.UpstreamSpec{
			UpstreamType: &gloov1.UpstreamSpec_Static{
				Static: input.GetStatic(),
			},
		}
	}
	return nil
}
