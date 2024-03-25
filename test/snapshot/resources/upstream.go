package resources

import (
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/options/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HttpbinUpstream is an Upstream that represents the httpbin service on port 8000
var HttpbinUpstream = &gloov1.Upstream{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "httpbin-v1-httpbin-8000",
		Namespace: gloodefaults.GlooSystem,
	},
	Spec: gloov1.UpstreamSpec{
		DiscoveryMetadata: &gloov1.DiscoveryMetadata{
			Labels: map[string]string{
				"app":     "httpbin-v1",
				"service": "httpbin-v1",
			},
		},
		UpstreamType: &gloov1.UpstreamSpec_Kube{
			Kube: &kubernetes.UpstreamSpec{
				Selector: map[string]string{
					"app": "httpbin-v1",
				},
				ServiceNamespace: "httpbin-v1",
				ServiceName:      "httpbin-v1",
				ServicePort:      uint32(8000),
			},
		},
	},
}
