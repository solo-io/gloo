package aws_test

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	. "github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/aws"
)

var _ = Describe("Plugin", func() {
	var (
		params   plugins.Params
		plugin   plugins.Plugin
		upstream *v1.Upstream
		out      *envoyapi.Cluster
	)
	BeforeEach(func() {
		plugin = NewAwsPlugin()
		plugin.Init(plugins.InitParams{})
		upstream = &v1.Upstream{
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Aws{
					Aws: &aws.UpstreamSpec{
						LambdaFunctions: []*aws.LambdaFunctionSpec{{
							LogicalName:        "foo",
							LambdaFunctionName: "foo",
							Qualifier:          "v1",
						}},
						Region:    "us-east1",
						SecretRef: "secretref",
					},
				},
			},
		}
		out = &envoyapi.Cluster{}
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

	It("should process upstream with secrets", func() {
		err := plugin.(plugins.UpstreamPlugin).ProcessUpstream(params, upstream, out)
		Expect(err).NotTo(HaveOccurred())
	})

})
