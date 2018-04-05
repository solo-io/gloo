package crd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-storage/crd"
	"github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo/pkg/protoutil"
)

var _ = Describe("Conversion", func() {
	Describe("UpstreamToCrd", func() {
		It("Converts a gloo upstream to crd", func() {
			us := helpers.NewTestUpstream1()
			annotations := map[string]string{"foo": "bar"}
			us.Metadata = &v1.Metadata{
				Annotations: annotations,
			}
			upCrd, err := UpstreamToCrd("foo", us)
			Expect(err).NotTo(HaveOccurred())
			Expect(upCrd.Name).To(Equal(us.Name))
			Expect(upCrd.Namespace).To(Equal("foo"))
			Expect(upCrd.Annotations).To(Equal(annotations))
			Expect(upCrd.Spec).NotTo(BeNil())
			spec := *upCrd.Spec
			// removed parts
			Expect(spec["name"]).To(BeNil())
			Expect(spec["metadata"]).To(BeNil())
			Expect(spec["status"]).To(BeNil())
			Expect(spec["annotations"]).To(BeNil())

			// shifted parts
			Expect(spec["type"]).To(Equal(us.Type))
			m, err := protoutil.MarshalMap(us.Spec)
			Expect(err).To(BeNil())
			Expect(spec["spec"]).To(Equal(m))
			var fnsInSpec []interface{}
			for _, fn := range us.Functions {
				fnSpec, err := protoutil.MarshalMap(fn.Spec)
				Expect(err).To(BeNil())
				fnsInSpec = append(fnsInSpec, map[string]interface{}{
					"name": fn.Name,
					"spec": fnSpec,
				})
			}
			Expect(spec["functions"]).To(Equal(fnsInSpec))
		})
	})
	Describe("VirtualhostToCrd", func() {
		It("Converts a gloo virtualhost to crd", func() {
			vHost := helpers.NewTestVirtualHost("foo", helpers.NewTestRoute1())
			annotations := map[string]string{"foo": "bar"}
			vHost.Metadata = &v1.Metadata{
				Annotations: annotations,
			}
			vhCrd, err := VirtualHostToCrd("foo", vHost)
			Expect(err).NotTo(HaveOccurred())
			Expect(vhCrd.Name).To(Equal(vHost.Name))
			Expect(vhCrd.Namespace).To(Equal("foo"))
			Expect(vhCrd.Spec).NotTo(BeNil())
			Expect(vhCrd.Annotations).To(Equal(annotations))
			spec := *vhCrd.Spec
			// removed parts
			Expect(spec["name"]).To(BeNil())
			Expect(spec["metadata"]).To(BeNil())
			Expect(spec["status"]).To(BeNil())
			Expect(spec["annotations"]).To(BeNil())

			//// shifted parts
			vhostMap, err := protoutil.MarshalMap(vHost)
			Expect(err).To(BeNil())
			Expect(spec["routes"]).To(Equal(vhostMap["routes"]))
		})
	})
	Describe("VirtualhostFromCrd", func() {
		It("Converts a gloo virtualhost to crd", func() {
			vHost := helpers.NewTestVirtualHost("foo", helpers.NewTestRoute1())
			annotations := map[string]string{"foo": "bar"}
			vHost.Metadata = &v1.Metadata{
				Annotations: annotations,
			}
			vhCrd, err := VirtualHostToCrd("foo", vHost)
			Expect(err).NotTo(HaveOccurred())
			Expect(vhCrd.Name).To(Equal(vHost.Name))
			Expect(vhCrd.Namespace).To(Equal("foo"))
			Expect(vhCrd.Spec).NotTo(BeNil())
			spec := *vhCrd.Spec
			// removed parts
			Expect(spec["name"]).To(BeNil())
			Expect(spec["metadata"]).To(BeNil())
			Expect(spec["status"]).To(BeNil())

			//// shifted parts
			vhostMap, err := protoutil.MarshalMap(vHost)
			Expect(err).To(BeNil())
			Expect(spec["routes"]).To(Equal(vhostMap["routes"]))

			// bring it back now
			outVhost, err := VirtualHostFromCrd(vhCrd)
			vHost.Metadata = &v1.Metadata{
				ResourceVersion: vhCrd.ResourceVersion,
				Namespace:       vhCrd.Namespace,
				Annotations:     annotations,
			}
			Expect(err).To(BeNil())
			Expect(outVhost).To(Equal(vHost))
		})
	})
})
