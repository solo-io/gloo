package server

import (
	apiv1 "github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/pkg/source"
	"github.com/solo-io/gloo/pkg/protoutil"
)

type handler struct {
	poller *source.Poller
}

func (h *handler) Update(u *apiv1.Upstream) {
	h.poller.Update(toUpstream(u))
}

func (h *handler) Remove(name string) {
	h.poller.Remove(name)
}

func toUpstream(u *apiv1.Upstream) source.Upstream {
	s, _ := protoutil.MarshalMap(u.Spec)

	return source.Upstream{
		Name:      u.Name,
		Type:      u.Type,
		Spec:      s,
		Functions: toFunctions(u.Functions),
	}
}

func toFunctions(functions []*apiv1.Function) []source.Function {
	result := make([]source.Function, len(functions))
	for i, f := range functions {
		s, _ := protoutil.MarshalMap(f.Spec)
		result[i] = source.Function{
			Name: f.Name,
			Spec: s,
		}
	}
	return result
}
