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

	var res []string
	targetRefs := lisOpt.Spec.GetTargetRefs()
	if len(targetRefs) == 0 {
		return res
	}

	foundNns := map[string]any{}

	for _, targetRef := range targetRefs {
		if targetRef == nil ||
			targetRef.GetGroup() != gwv1.GroupName ||
			targetRef.GetKind() != wellknown.GatewayKind {
			continue
		}

		ns := targetRef.GetNamespace().GetValue()
		if ns == "" {
			ns = lisOpt.GetNamespace()
		}

		targetNN := types.NamespacedName{
			Namespace: ns,
			Name:      targetRef.GetName(),
		}
		foundNns[targetNN.String()] = struct{}{}
	}

	for targetNN := range foundNns {
		res = append(res, targetNN)
	}

	return res
}
