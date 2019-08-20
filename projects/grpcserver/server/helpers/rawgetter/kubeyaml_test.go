package rawgetter_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
)

var getter rawgetter.RawGetter

var _ = Describe("Kube Yaml RawGetter", func() {

	BeforeEach(func() {
		getter = rawgetter.NewKubeYamlRawGetter()
	})

	Describe("GetRaw", func() {
		It("works", func() {
			resource := &gloov1.Proxy{
				Status: core.Status{},
				Metadata: core.Metadata{
					Namespace: "namespace",
					Name:      "name",
				},
			}

			expected := &v1.Raw{
				FileName: "name.yaml",
				Content:  "apiVersion: gloo.solo.io/v1\nkind: Proxy\nmetadata:\n  creationTimestamp: null\n  name: name\n  namespace: namespace\nspec: {}\nstatus: {}\n",
			}

			actual := getter.GetRaw(context.Background(), resource, gloov1.ProxyCrd)
			ExpectEqualProtoMessages(actual, expected)
		})
	})
})
