package krtquery

import (
	gwapi "sigs.k8s.io/gateway-api/apis/v1beta1"

	"istio.io/istio/pkg/kube/krt"
)

type HTTPRouteAttachement = Attachment[*gwapi.HTTPRoute]

// TODO implement allowed types/namespaces
// TOOD delegation
func HTTPRouteAttachements(
	Gateways krt.Collection[*gwapi.Gateway],
	HTTPRoutes krt.Collection[*gwapi.HTTPRoute],
) krt.Collection[HTTPRouteAttachement] {
	return krt.NewManyCollection[*gwapi.HTTPRoute, HTTPRouteAttachement](
		HTTPRoutes,
		func(ctx krt.HandlerContext, hr *gwapi.HTTPRoute) []HTTPRouteAttachement {
			var out []HTTPRouteAttachement
			for _, ref := range hr.Spec.ParentRefs {
				attachment := attachementFromParentRef(hr, ref)
				if attachment.AttachedToGateway(ctx, Gateways) {
					out = append(out, attachment)
				}
				// TODO else report status
			}
			return out
		},
	)
}
