package xds_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
)

var _ = Describe("Cache", func() {

	It("SnapshotCacheKeys returns the keys formatted correctly", func() {
		owner, namespace1, namespace2, name1, name2 := "owner", "namespace1", "namespace2", "name1", "name2"
		p1 := v1.NewProxy(namespace1, name1)
		p1.Metadata.Labels = map[string]string{utils.ProxyTypeKey: owner}
		p2 := v1.NewProxy(namespace2, name2)
		p2.Metadata.Labels = map[string]string{utils.ProxyTypeKey: owner}
		proxies := []*v1.Proxy{p1, p2}
		expectedKeys := []string{fmt.Sprintf("%v~%v~%v", owner, namespace1, name1), fmt.Sprintf("%v~%v~%v", owner, namespace2, name2)}
		actualKeys := xds.SnapshotCacheKeys(proxies)
		Expect(actualKeys).To(BeEquivalentTo(expectedKeys))
	})
})
