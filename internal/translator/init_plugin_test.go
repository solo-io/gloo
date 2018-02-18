package translator

import (
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/aws"
	// . "github.com/solo-io/gloo-testing/helpers"
	// . "github.com/solo-io/gloo/internal/translator"

	"github.com/solo-io/gloo/pkg/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InitPlugin", func() {

	It("should do correct weights for functions", func() {
		funcs := []plugin.FunctionPlugin{&aws.Plugin{}}
		initPlugin := newInitializerPlugin(funcs)

		outroute := envoyroute.Route{Metadata: new(envoycore.Metadata)}
		outroute.Metadata.FilterMetadata = make(map[string]*types.Struct)
		inroute := &v1.Route{
			MultipleDestinations: []*v1.WeightedDestination{{
				Destination: &v1.Destination{
					DestinationType: &v1.Destination_Function{
						Function: &v1.FunctionDestination{
							UpstreamName: "my-upstream",
							FunctionName: "func1",
						},
					},
				},
				Weight: 10,
			}, {
				Destination: &v1.Destination{
					DestinationType: &v1.Destination_Function{
						Function: &v1.FunctionDestination{
							UpstreamName: "my-upstream",
							FunctionName: "func2",
						},
					},
				},
				Weight: 10,
			}, {
				Destination: &v1.Destination{
					DestinationType: &v1.Destination_Upstream{
						Upstream: &v1.UpstreamDestination{
							Name: "my-upstream2",
						},
					},
				},
				Weight: 10,
			}},
		}
		params := &plugin.RoutePluginParams{}
		err := initPlugin.ProcessRoute(params, inroute, &outroute)
		Expect(err).NotTo(HaveOccurred())
		clusterweight := outroute.Action.(*envoyroute.Route_Route).Route.ClusterSpecifier.(*envoyroute.RouteAction_WeightedClusters).WeightedClusters
		routeweight := clusterweight.TotalWeight.Value
		Expect(clusterweight.Clusters[0].Weight.Value).To(BeEquivalentTo(20))
		Expect(clusterweight.Clusters[1].Weight.Value).To(BeEquivalentTo(10))
		Expect(routeweight).To(BeEquivalentTo(30))

		metadata := outroute.Metadata
		clusterroutemeta := metadata.FilterMetadata[filterName].Fields["my-upstream"].Kind.(*types.Value_StructValue).StructValue.Fields
		totalfuncweight := clusterroutemeta[multiFunctionWeightDestinationKey].Kind.(*types.Value_NumberValue).NumberValue
		Expect(totalfuncweight).To(BeEquivalentTo(20))
		functionslist := clusterroutemeta[multiFunctionDestinationKey].Kind.(*types.Value_ListValue).ListValue.Values
		for _, fl := range functionslist {
			weight := fl.Kind.(*types.Value_StructValue).StructValue.Fields["weight"].Kind.(*types.Value_NumberValue).NumberValue
			Expect(weight).To(BeEquivalentTo(10))

		}

	})

})
