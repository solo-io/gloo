package query

import (
	"fmt"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	VirtualHostOptionTargetField = "virtualHostOption.targetRef"
)

func IterateIndices(f func(client.Object, string, client.IndexerFunc) error) error {
	return f(&solokubev1.VirtualHostOption{}, VirtualHostOptionTargetField, virtualHostOptionTargetRefIndexer)
}

func virtualHostOptionTargetRefIndexer(obj client.Object) []string {
	vhOpt, ok := obj.(*solokubev1.VirtualHostOption)
	if !ok {
		panic(fmt.Sprintf("wrong type %T provided to indexer. expected gateway.solo.io.VirtualHostOption", obj))
	}

	return utils.IndexTargetRefsNnk(vhOpt.Spec.GetTargetRefs(), vhOpt.GetNamespace(), utils.ListenerTargetRefGVKs)
}
