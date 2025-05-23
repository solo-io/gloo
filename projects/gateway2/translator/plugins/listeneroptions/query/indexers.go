package query

import (
	"fmt"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ListenerOptionTargetField = "listenerOption.targetRef"
)

func IterateIndices(f func(client.Object, string, client.IndexerFunc) error) error {
	return f(&solokubev1.ListenerOption{}, ListenerOptionTargetField, listenerOptionTargetRefIndexer)
}

func listenerOptionTargetRefIndexer(obj client.Object) []string {
	lisOpt, ok := obj.(*solokubev1.ListenerOption)
	if !ok {
		panic(fmt.Sprintf("wrong type %T provided to indexer. expected gateway.solo.io.ListenerOption", obj))
	}

	return utils.IndexTargetRefsNnk(lisOpt.Spec.GetTargetRefs(), lisOpt.GetNamespace(), utils.ListenerTargetRefGVKs)

}
