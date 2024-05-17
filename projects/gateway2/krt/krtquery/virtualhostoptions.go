package krtquery

import (
	"istio.io/istio/pkg/kube/krt"
	gwapi "sigs.k8s.io/gateway-api/apis/v1beta1"

	solov1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
)

type VirtualHostOptionAttachement = Attachment[*solov1.VirtualHostOption]

func VirtualHostOptionAttachements(
	Gateways krt.Collection[*gwapi.Gateway],
	VirtualHostOptions krt.Collection[*solov1.VirtualHostOption],
) krt.Collection[VirtualHostOptionAttachement] {
	return krt.NewCollection[*solov1.VirtualHostOption, VirtualHostOptionAttachement](
		VirtualHostOptions,
		func(ctx krt.HandlerContext, vho *solov1.VirtualHostOption) *Attachment[*solov1.VirtualHostOption] {
			attachement := attachementFromTargetRef(vho, vho.Spec.TargetRef)
			if attachement.AttachedToGateway(ctx, Gateways) {
				return &attachement
			}
			// TODO else report status
			return nil
		},
	)
}
