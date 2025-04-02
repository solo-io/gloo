package query

import (
	"fmt"

	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	RouteOptionTargetField = "routeOption.targetRef"
)

func IterateIndices(f func(client.Object, string, client.IndexerFunc) error) error {
	return f(&solokubev1.RouteOption{}, RouteOptionTargetField, routeOptionTargetRefIndexer)
}

func routeOptionTargetRefIndexer(obj client.Object) []string {
	rtOpt, ok := obj.(*solokubev1.RouteOption)
	if !ok {
		panic(fmt.Sprintf("wrong type %T provided to indexer. expected gateway.solo.io.RouteOption", obj))
	}

	return utils.IndexTargetRefs(rtOpt.Spec.GetTargetRefs(), rtOpt.GetNamespace(), wellknown.HTTPRouteKind)
}
