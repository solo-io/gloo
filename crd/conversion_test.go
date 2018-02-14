package crd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/glue-storage/crd"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/protoutil"
	"github.com/solo-io/glue/test/helpers"
)

var _ = Describe("Conversion", func() {
	Describe("UpstreamToCrd", func() {
		It("Converts a glue upstream to crd", func() {
			us := helpers.NewTestUpstream1()
			upCrd, err := UpstreamToCrd("foo", us)
			Expect(err).NotTo(HaveOccurred())
			Expect(upCrd.Name).To(Equal(us.Name))
			Expect(upCrd.Namespace).To(Equal("foo"))
			Expect(upCrd.Spec).NotTo(BeNil())
			spec := *upCrd.Spec
			// removed parts
			Expect(spec["name"]).To(BeNil())
			Expect(spec["metadata"]).To(BeNil())
			Expect(spec["status"]).To(BeNil())

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
		It("Converts a glue virtualhost to crd", func() {
			vHost := helpers.NewTestVirtualHost("foo", helpers.NewTestRoute1())
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
		})
	})
	Describe("VirtualhostFromCrd", func() {
		It("Converts a glue virtualhost to crd", func() {
			vHost := helpers.NewTestVirtualHost("foo", helpers.NewTestRoute1())
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
			}
			Expect(err).To(BeNil())
			Expect(outVhost).To(Equal(vHost))
		})
	})
})
