package rawgetter_test

import (
	"context"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Describe("InitResourceFromYamlString", func() {
		It("works", func() {
			ref := &core.ResourceRef{
				Name:      "default-petstore-8080",
				Namespace: "gloo-system",
			}
			file, openErr := os.Open("./fixtures/example-upstream.yaml")
			Expect(openErr).NotTo(HaveOccurred())

			yamlString, readErr := ioutil.ReadAll(file)
			Expect(readErr).NotTo(HaveOccurred())

			emptyUpstream := &gloov1.Upstream{}
			parseErr := getter.InitResourceFromYamlString(context.TODO(), string(yamlString), ref, emptyUpstream)

			Expect(parseErr).NotTo(HaveOccurred())

			status := emptyUpstream.GetStatus()
			Expect((&status).GetReportedBy()).To(Equal("gloo"))

			metadata := emptyUpstream.GetMetadata()
			Expect((&metadata).GetResourceVersion()).To(Equal("58959"))
		})

		It("fails when the namespace is edited", func() {
			ref := &core.ResourceRef{
				Name:      "default-petstore-8080",
				Namespace: "gloo-system",
			}
			file, openErr := os.Open("./fixtures/changed-upstream-ref.yaml")
			Expect(openErr).NotTo(HaveOccurred())

			yamlString, readErr := ioutil.ReadAll(file)
			Expect(readErr).NotTo(HaveOccurred())

			emptyUpstream := &gloov1.Upstream{}
			parseErr := getter.InitResourceFromYamlString(context.TODO(), string(yamlString), ref, emptyUpstream)

			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr.Error()).To(ContainSubstring(rawgetter.EditedRefError(ref).Error()))
		})

		It("fails when the resource version is omitted", func() {
			ref := &core.ResourceRef{
				Name:      "default-petstore-8080",
				Namespace: "gloo-system",
			}
			file, openErr := os.Open("./fixtures/blank-resource-version-upstream.yaml")
			Expect(openErr).NotTo(HaveOccurred())

			yamlString, readErr := ioutil.ReadAll(file)
			Expect(readErr).NotTo(HaveOccurred())

			emptyUpstream := &gloov1.Upstream{}
			parseErr := getter.InitResourceFromYamlString(context.TODO(), string(yamlString), ref, emptyUpstream)

			Expect(parseErr).To(HaveOccurred())
			Expect(parseErr.Error()).To(ContainSubstring(rawgetter.NoResourceVersionError(ref).Error()))
		})
	})
})
