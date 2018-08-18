package aws_test

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	awsapi "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	. "github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/aws"
)

var _ = Describe("Plugin", func() {
	var (
		params   plugins.Params
		plugin   plugins.Plugin
		upstream *v1.Upstream
		route    *v1.Route
		out      *envoyapi.Cluster
		outroute *envoyroute.Route
	)
	BeforeEach(func() {
		plugin = NewAwsPlugin()
		plugin.Init(plugins.InitParams{})
		upstreamName := "up"
		clusterName := upstreamName
		funcname := "foo"
		upstream = &v1.Upstream{
			Metadata: core.Metadata{
				Name: upstreamName,
				// TODO(yuval-k): namespace
				Namespace: "",
			},
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Aws{
					Aws: &aws.UpstreamSpec{
						LambdaFunctions: []*aws.LambdaFunctionSpec{{
							LogicalName:        funcname,
							LambdaFunctionName: "foo",
							Qualifier:          "v1",
						}},
						Region:    "us-east1",
						SecretRef: "secretref",
					},
				},
			},
		}
		route = &v1.Route{
			Action: &v1.Route_RouteAction{
				RouteAction: &v1.RouteAction{
					Destination: &v1.RouteAction_Single{
						Single: &v1.Destination{
							UpstreamName: upstreamName,
							DestinationSpec: &v1.DestinationSpec{
								DestinationType: &v1.DestinationSpec_Aws{
									Aws: &awsapi.DestinationSpec{
										LogicalName: funcname,
									},
								},
							},
						},
					},
				},
			},
		}

		out = &envoyapi.Cluster{}
		outroute = &envoyroute.Route{
			Action: &envoyroute.Route_Route{
				Route: &envoyroute.RouteAction{
					ClusterSpecifier: &envoyroute.RouteAction_Cluster{
						Cluster: clusterName,
					},
				},
			},
		}

		params.Snapshot = &v1.Snapshot{
			SecretList: v1.SecretList{{
				Metadata: core.Metadata{
					Name: "secretref",
					// TODO(yuval-k): namespace
					Namespace: "",
				},
				Data: map[string]string{
					"access_key": "access_key",
					"secret_key": "secret_key",
				},
			}},
		}
	})
	Context("upstreams", func() {

		It("should process upstream with secrets", func() {
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).NotTo(HaveOccurred())
			Expect(out.ExtensionProtocolOptions).NotTo(BeEmpty())
			Expect(out.ExtensionProtocolOptions).To(HaveKey(FilterName))
		})
		It("should error upstream with no secrets", func() {
			params.Snapshot.SecretList = nil
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(HaveOccurred())
		})

		It("should error upstream with no access_key", func() {
			delete(params.Snapshot.SecretList[0].Data, "access_key")
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(HaveOccurred())
		})
		It("should error upstream with no secret_key", func() {
			delete(params.Snapshot.SecretList[0].Data, "secret_key")
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(HaveOccurred())
		})
		It("should error upstream with invalid access_key", func() {
			params.Snapshot.SecretList[0].Data["access_key"] = "\xffbinary!"
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(HaveOccurred())
		})
		It("should error upstream with invalid secret_key", func() {
			params.Snapshot.SecretList[0].Data["secret_key"] = "\xffbinary!"
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).To(HaveOccurred())
		})

		It("should not process and not error with non aws upstream", func() {
			upstream.UpstreamSpec.UpstreamType = nil
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).NotTo(HaveOccurred())
			Expect(out.ExtensionProtocolOptions).To(BeEmpty())
			Expect(outroute.PerFilterConfig).NotTo(HaveKey(FilterName))

		})
	})

	Context("routes", func() {

		var destination *v1.Destination

		BeforeEach(func() {
			err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
			Expect(err).NotTo(HaveOccurred())
			destination = route.Action.(*v1.Route_RouteAction).RouteAction.Destination.(*v1.RouteAction_Single).Single
		})

		It("should process route", func() {
			err := plugin.(plugins.RoutePlugin).ProcessRoute(params, route, outroute)
			Expect(err).NotTo(HaveOccurred())
			Expect(outroute.PerFilterConfig).To(HaveKey(FilterName))
		})

		It("should not process with no spec", func() {
			destination.DestinationSpec = nil

			err := plugin.(plugins.RoutePlugin).ProcessRoute(params, route, outroute)
			Expect(err).NotTo(HaveOccurred())
			Expect(outroute.PerFilterConfig).NotTo(HaveKey(FilterName))
		})

		It("should not process with no spec", func() {
			Skip("redo this when we have more destination type")
			// destination.DestinationSpec.DestinationType =

			err := plugin.(plugins.RoutePlugin).ProcessRoute(params, route, outroute)
			Expect(err).NotTo(HaveOccurred())
			Expect(outroute.PerFilterConfig).NotTo(HaveKey(FilterName))
		})

	})
})
