package query

import (
	"fmt"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
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

	var res []string
	targetRefs := vhOpt.Spec.GetTargetRefs()
	if len(targetRefs) == 0 {
		return res
	}

	foundNns := map[string]any{}

	// TODO: fix up logic to be cleaner
	for _, targetRef := range targetRefs {
		if targetRef == nil {
			continue
		}
		if targetRef.GetGroup() != gwv1.GroupName {
			continue
		}
		if targetRef.GetKind() != wellknown.GatewayKind {
			continue
		}
		ns := targetRef.GetNamespace().GetValue()
		if ns == "" {
			ns = vhOpt.GetNamespace()
		}
		targetNN := types.NamespacedName{
			Namespace: ns,
			Name:      targetRef.GetName(),
		}

		foundNns[targetNN.String()] = struct{}{}
	}

	for k := range foundNns {
		res = append(res, k)
	}

	fmt.Printf("virtualHostOptionTargetRefIndexer, %d items in list: %+v\n", len(res), res)

	return res
}
