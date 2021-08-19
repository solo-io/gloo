package glooinstance_handler

import (
	"sort"

	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
)

func sortGlooInstances(glooInstances []*rpc_edge_v1.GlooInstance) {
	sort.Slice(glooInstances, func(i, j int) bool {
		x := glooInstances[i]
		y := glooInstances[j]
		return x.GetMetadata().GetNamespace()+x.GetMetadata().GetName() < y.GetMetadata().GetNamespace()+y.GetMetadata().GetName()
	})
}
