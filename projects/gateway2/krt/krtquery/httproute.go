package krtquery

import (
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"istio.io/istio/pkg/kube/krt"
)

type HTTPRouteAttachement = Attachment[*gwv1.HTTPRoute]

// TODO implement allowed types/namespaces
// TOOD delegation
func HTTPRouteAttachements(
	Gateways krt.Collection[*gwv1.Gateway],
	HTTPRoutes krt.Collection[*gwv1.HTTPRoute],
) krt.Collection[HTTPRouteAttachement] {
	return krt.NewManyCollection[*gwv1.HTTPRoute, HTTPRouteAttachement](
		HTTPRoutes,
		func(ctx krt.HandlerContext, hr *gwv1.HTTPRoute) []HTTPRouteAttachement {
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
