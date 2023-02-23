package installer_test

import (
	"context"
	"time"

	"github.com/solo-io/solo-projects/test/kubeutils"

	"github.com/solo-io/solo-projects/test/kubeutils/installer"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Install MultiClusterInstaller", func() {

	var (
		ctx     context.Context
		cancel  context.CancelFunc
		manager *installer.MultiClusterInstaller
	)

	Context("Non-Kube", func() {
		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())

			manager = installer.NewMultiCluster(kubeutils.NewNoOpOrchestrator())
		})

		AfterEach(func() {
			cancel()
		})

		It("can execute multiple installs in parallel", func() {
			expensiveInstaller := installer.NewNoOpInstaller(time.Second*3, nil)

			for i := 0; i < 10; i++ {
				manager.RegisterInstaller(expensiveInstaller)
			}

			beforeInstall := time.Now()
			err := manager.Install(ctx)
			Expect(err).NotTo(HaveOccurred())

			// We executed 10 installs, each taking 3 seconds, with a 0 second buffer by the manager
			// 30 seconds: if entirely serial
			// 3 seconds: if entirely parallel
			installDuration := time.Since(beforeInstall)
			Expect(installDuration.Seconds()).To(BeNumerically(">", 2))
			Expect(installDuration.Seconds()).To(BeNumerically("<", 4))
		})

		It("can execute multiple uninstalls in parallel", func() {
			expensiveInstaller := installer.NewNoOpInstaller(time.Second*3, nil)

			for i := 0; i < 10; i++ {
				manager.RegisterInstaller(expensiveInstaller)
			}

			beforeUninstall := time.Now()
			err := manager.Uninstall(ctx)
			Expect(err).NotTo(HaveOccurred())

			// We executed 10 uninstalls, each taking 3 seconds, with a 0 second buffer by the manager
			// 30 seconds: if entirely serial
			// 3 seconds: if entirely parallel
			uninstallDuration := time.Since(beforeUninstall)
			Expect(uninstallDuration.Seconds()).To(BeNumerically(">", 2))
			Expect(uninstallDuration.Seconds()).To(BeNumerically("<", 4))
		})

	})

})
