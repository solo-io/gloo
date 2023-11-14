package federation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Util", func() {

	Context("owner labels and annotations", func() {
		It("label and annotation should be the same for short values", func() {
			name := "fed-upstream"
			fedUpstream := &v1.FederatedUpstream{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      name,
				},
			}
			anno := federation.GetOwnerAnnotation(fedUpstream)[federation.HubOwner]
			label := federation.GetOwnerLabel(fedUpstream)[federation.HubOwner]
			fullIdentifier := federation.GetIdentifier(fedUpstream)
			Expect(anno).To(Equal(fullIdentifier))
			Expect(label).To(Equal(fullIdentifier))
		})

		It("label should be shortened for long values", func() {
			name := "fed-upstream-with-very-long-name-that-is-over-one-hundred-characters-long-00000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999990000000000"
			fedUpstream := &v1.FederatedUpstream{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      name,
				},
			}
			anno := federation.GetOwnerAnnotation(fedUpstream)[federation.HubOwner]
			label := federation.GetOwnerLabel(fedUpstream)[federation.HubOwner]
			fullIdentifier := federation.GetIdentifier(fedUpstream)
			// annotation should contain the full id, but label should be shortened
			Expect(anno).To(Equal(fullIdentifier))
			Expect(label).To(HaveLen(63))
			Expect(len(label)).To(BeNumerically("<", len(anno)))
		})
	})

	Context("merge", func() {
		It("can merge with empty values", func() {
			m1 := map[string]string{}
			m2 := map[string]string{"key": "value"}
			merged := federation.Merge(m1, m2)
			Expect(merged).To(HaveLen(1))
			Expect(merged["key"]).To(Equal("value"))
		})

		It("can merge duplicates", func() {
			m1 := map[string]string{
				"a":   "b",
				"key": "oldValue",
			}
			m2 := map[string]string{
				"key": "newValue",
				"c":   "d",
			}
			merged := federation.Merge(m1, m2)
			Expect(merged).To(HaveLen(3))
			Expect(merged["a"]).To(Equal("b"))
			Expect(merged["c"]).To(Equal("d"))
			Expect(merged["key"]).To(Equal("newValue"))
		})

	})

})
