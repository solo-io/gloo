package translator

import (
	"fmt"
	"sort"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
)

type route struct {
	path         string
	destinations []destination
}

type destination struct {
	upstreamRef core.ResourceRef
	weight      uint32
}

func RoutePath(typeName, fieldName string) string {
	return fmt.Sprintf("/%v.%v", typeName, fieldName)
}

func buildRoutes(resolverMap *v1.ResolverMap) []route {
	var routes []route
	for typeName, typeResolver := range resolverMap.Types {
		for fieldName, fieldResolver := range typeResolver.Fields {
			glooResolver, ok := fieldResolver.Resolver.(*v1.FieldResolver_GlooResolver)
			if !ok {
				continue
			}
			routes = append(routes, route{
				path:         RoutePath(typeName, fieldName),
				destinations: destinationsForFunction(glooResolver.GlooResolver),
			})
		}
	}
	sort.SliceStable(routes, func(i, j int) bool {
		return routes[i].path < routes[j].path
	})
	return routes
}

func destinationsForFunction(resolver *v1.GlooResolver) []destination {
	switch function := resolver.Destination.Function.(type) {
	case *v1.GlooResolver_SingleFunction:
		return []destination{
			{
				upstreamName: function.SingleFunction.Upstream,
				functionName: function.SingleFunction.Function,
			},
		}
	case *v1.GlooResolver_MultiFunction:
		var dests []destination
		for _, weightedFunc := range function.MultiFunction.WeightedFunctions {
			dests = append(dests, destination{
				upstreamName: weightedFunc.Function.Upstream,
				functionName: weightedFunc.Function.Function,
				weight:       weightedFunc.Weight,
			})
		}
		return dests
	}
	panic("unknown function time")
}
