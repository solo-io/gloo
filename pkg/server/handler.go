package server

import (
	"fmt"

	"github.com/solo-io/glue-discovery/pkg/source"
	apiv1 "github.com/solo-io/glue/pkg/api/types/v1"
	solov1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
)

type handler struct {
	poller *source.Poller
}

func (h *handler) Update(u *solov1.Upstream) {
	h.poller.Update(toUpstream(u))
}

func (h *handler) Remove(ID string) {
	h.poller.Remove(ID)
}

func toUpstream(u *solov1.Upstream) source.Upstream {
	return source.Upstream{
		ID:        toID(u),
		Namespace: u.Namespace,
		Name:      u.Name,
		Type:      toType(u.Spec.Type),
		Spec:      u.Spec.Spec,
		Functions: toFunctions(u.Spec.Functions),
	}
}

func toFunctions(functions []apiv1.Function) []source.Function {
	result := make([]source.Function, len(functions))
	for i, f := range functions {
		result[i] = source.Function{
			Name: f.Name,
			Spec: f.Spec,
		}
	}
	return result
}

func toID(u *solov1.Upstream) string {
	return fmt.Sprintf("%s/%s", u.Namespace, u.Name)
}

func toType(t apiv1.UpstreamType) string {
	return fmt.Sprintf("%v", t)
}
