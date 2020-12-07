package clientside_sharding_test

import (
	"context"

	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	. "github.com/solo-io/solo-projects/test/regressions/internal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// Regular extAuth as run in other suites
var _ = Describe("ExtAuth tests", func() {
	sharedInputs := ExtAuthTestInputs{}

	BeforeEach(func() {
		sharedInputs.TestHelper = testHelper
	})

	Context("running ExtAuth tests", func() {
		RunExtAuthTests(&sharedInputs)
	})

	AfterEach(func() {
		ctx := context.TODO()
		kubeCache := kube.NewKubeCache(ctx)
		cfg, err := kubeutils.GetConfig("", "")
		authConfigClientFactory := &factory.KubeResourceClientFactory{
			Crd:         extauthv1.AuthConfigCrd,
			Cfg:         cfg,
			SharedCache: kubeCache,
		}
		authConfigClient, err := extauthv1.NewAuthConfigClient(ctx, authConfigClientFactory)
		Expect(err).NotTo(HaveOccurred(), "should create auth config client")
		authConfigs, err := authConfigClient.List(testHelper.InstallNamespace, clients.ListOpts{})
		for _, authConfig := range authConfigs {
			err = authConfigClient.Delete(testHelper.InstallNamespace, authConfig.Metadata.Name, clients.DeleteOpts{
				Ctx:            ctx,
				IgnoreNotExist: true,
			})
			Expect(err).NotTo(HaveOccurred(), "should delete authconfigs on cleanup")
		}
	})

})
