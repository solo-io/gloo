package extauth_test

import (
	"context"
	"time"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth/test_fixtures"

	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/util/rand"

	. "github.com/onsi/ginkgo/v2"
)

// syncer.TranslatorSyncerExtension Sync methods are assumed to be thread-safe.
// This is a critical assumption, because they are utilized both by our translation and
// validation pipelines. These tests exist to ensure that we don't break this assumption.
// In the past, when we have not respected this, we have seen race conditions:
//
//	https://github.com/solo-io/gloo/pull/7207 is one example
var _ = Describe("ExtAuth Translation - Concurrency Tests", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		extension syncer.TranslatorSyncerExtension
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		extension = extauth.NewTranslatorSyncerExtension(ctx, syncer.TranslatorSyncerExtensionParams{
			Hasher: translator.EnvoyCacheResourcesListToFnvHash,
		})
	})

	AfterEach(func() {
		cancel()
	})

	It("Can perform Sync concurrently without a panic", func() {
		syncerFunc := func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			reports := make(reporter.ResourceReports)
			settings := &v1.Settings{}
			snapshotSetter := &syncer.NoOpSnapshotSetter{}
			snapshot := &gloosnapshot.ApiSnapshot{
				AuthConfigs: test_fixtures.AuthConfigSlice(writeNamespace, 10),
				Proxies: test_fixtures.ProxySlice(writeNamespace, test_fixtures.ResourceFrequency{
					Proxies:              2,
					VirtualHostsPerProxy: 5,
					RoutesPerVirtualHost: 10,
				}),
			}

			extension.Sync(ctx, snapshot, settings, snapshotSetter, reports)
		}

		errorGroup := errgroup.Group{}
		for i := 0; i < 100; i++ {
			errorGroup.Go(func() error {
				defer GinkgoRecover()
				syncerFunc()
				return nil
			})
		}

		Expect(errorGroup.Wait()).NotTo(HaveOccurred(), "Should not panic while executing sync from multiple goroutines")
	})
})
