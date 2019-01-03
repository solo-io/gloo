package syncer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
)

var _ = Describe("Vcs Sync", func() {

	It("should transition all resources", func() {
		// at this point, we want all resources to by synced, in the future we may change that plan
		Expect(transitionSettings(&gloov1.Settings{}, &gloov1.Settings{})).To(BeTrue())
		Expect(transitionProxies(&gloov1.Proxy{}, &gloov1.Proxy{})).To(BeTrue())
		Expect(transitionGateways(&gatewayv1.Gateway{}, &gatewayv1.Gateway{})).To(BeTrue())
		Expect(transitionVirtualServices(&gatewayv1.VirtualService{}, &gatewayv1.VirtualService{})).To(BeTrue())
		Expect(transitionResolverMaps(&sqoopv1.ResolverMap{}, &sqoopv1.ResolverMap{})).To(BeTrue())
		Expect(transitionSchemas(&sqoopv1.Schema{}, &sqoopv1.Schema{})).To(BeTrue())
	})
})
