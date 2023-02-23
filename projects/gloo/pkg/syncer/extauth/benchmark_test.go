package extauth_test

import (
	"context"
	"os"
	"sort"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	extauthsyncer "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth/test_fixtures"
)

// syncer.TranslatorSyncerExtension Sync methods are a frequently executed code path.
// We are progressively adding micro-benchmarking to this area of the code to ensure
// we don't introduce unintended latency to this space.
//
// At the moment, these tests (and our benchmarking approach) are not stable, so while
// we define the tests here, they are skipped in CI. The hope is that running locally
// we can see the measured difference now, and in the future we except to flush out
// our benchmarking tests.
var _ = Describe("ExtAuth Translation - Benchmarking Tests", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		if os.Getenv("GCLOUD_BUILD_ID") != "" {
			Skip("Skip Benchmark as they are not yet stable")
		}
	})

	AfterEach(func() {
		cancel()
	})

	// https://onsi.github.io/ginkgo/#the-spectimeout-and-nodetimeout-decorators
	// It would be preferable to apply a SpecTimeout to these tests to allow them to run for longer than the default
	// I opted to adjust them to fit within the 1-minute window for now, though am leaving this as a note
	// in case we want to expand the test in the future
	DescribeTable("xDS Snapshot Producer",
		func(producer extauthsyncer.XdsSnapshotProducer, resourceFrequency test_fixtures.ResourceFrequency, durationAssertion types.GomegaMatcher) {
			By("Seed Snapshot with resources")
			snapshot := &gloov1snap.ApiSnapshot{
				AuthConfigs: test_fixtures.AuthConfigSlice(writeNamespace, resourceFrequency.AuthConfigs),
				Proxies:     test_fixtures.ProxySlice(writeNamespace, resourceFrequency),
			}
			// Settings are only used by the XdsSnapshotProducer to validate that if a custom auth
			// server name is defined, that it maps to the name of a server configured in Settings
			settings := &gloov1.Settings{
				Extauth: &extauth.Settings{},
				NamedExtauth: map[string]*extauth.Settings{
					"custom-auth-server": {},
				},
			}
			reports := make(reporter.ResourceReports)
			reports.Accept(snapshot.AuthConfigs.AsInputResources()...)
			reports.Accept(snapshot.Proxies.AsInputResources()...)

			By("Execute ProduceXdsSnapshot")
			functionToBenchmark := func() {
				producer.ProduceXdsSnapshot(ctx, settings, snapshot, reports)
			}
			samples := 5
			results := make([]float64, samples)
			for i := 0; i < samples; i++ {
				results[i] = timeForFuncToComplete(functionToBenchmark)
			}
			sort.Float64s(results)

			By("Assert 80th percentile of ProduceXdsSnapshot completes within expected window")
			// 0:sample-1, 1:sample-2, 2:sample-3 [3:sample-4], 4:sample5
			eightyPercentile := results[samples-2]
			print(eightyPercentile) // Helpful to gather details about results in CI
			Expect(eightyPercentile).To(durationAssertion)
		},
		getBenchmarkingTestEntries(),
	)

})

// timeForFuncToComplete is a trivial replica of our benchmarking.TimeForFuncToComplete
// that benchmarking utility is going to be changed in the future and the purpose of this
// test is not to create strict benchmarking rules, but more demonstrate performance gains of
// some code changes. We expect to improve this benchmarking in the future
func timeForFuncToComplete(f func()) float64 {
	before := time.Now()
	f()
	return time.Since(before).Seconds()

}

func getBenchmarkingTestEntries() []TableEntry {
	if os.Getenv("GCLOUD_BUILD_ID") != "" {
		// We're running in CI
		return []TableEntry{
			Entry("proxySourcedXdsSnapshotProducer",
				extauthsyncer.NewProxySourcedXdsSnapshotProducer(),
				test_fixtures.ResourceFrequency{
					AuthConfigs:          1000,
					Proxies:              50,
					VirtualHostsPerProxy: 100,
					RoutesPerVirtualHost: 1000,
				},
				// Should take 12 +- 2 seconds
				// Recent CI results:
				// 12.230412642
				BeNumerically("~", 12, 2),
			),
		}
	}

	// We're running locally
	return []TableEntry{
		// (AC=1000, PXY=50, VH=100, R=1000)
		Entry("snapshotSourcedXdsSnapshotProducer (AC=1000, PXY=50, VH=100, R=1000)",
			extauthsyncer.NewSnapshotSourcedXdsSnapshotProducer(),
			test_fixtures.ResourceFrequency{
				AuthConfigs:          1000,
				Proxies:              50,
				VirtualHostsPerProxy: 100,
				RoutesPerVirtualHost: 1000,
			},
			// Should take 2 +- 1 seconds
			// Recent LOCAL results:
			// +1.630976e+000
			// +1.801740e+000
			// +1.680145e+000
			// +1.847926e+000
			// +2.005645e+000
			BeNumerically("~", 2, 1),
		),
		Entry("proxySourcedXdsSnapshotProducer (AC=1000, PXY=50, VH=100, R=1000)",
			extauthsyncer.NewProxySourcedXdsSnapshotProducer(),
			test_fixtures.ResourceFrequency{
				AuthConfigs:          1000,
				Proxies:              50,
				VirtualHostsPerProxy: 100,
				RoutesPerVirtualHost: 1000,
			},
			// Should take 7 +- 2 seconds
			// Recent LOCAL results:
			// +6.613514e+000
			// +6.603773e+000
			// +8.101754e+000
			// +7.821978e+000
			// +7.087994e+000
			// +8.308528e+000
			// +6.692507e+000
			BeNumerically("~", 7, 2),
		),
	}
}
