package extauth_test

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/solo-io/gloo/test/testutils"

	"github.com/onsi/gomega/gmeasure"
	"github.com/solo-io/gloo/test/gomega/matchers"

	"github.com/solo-io/gloo/test/ginkgo/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	extauthsyncer "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth/test_fixtures"

	gloohelpers "github.com/solo-io/gloo/test/helpers"
)

type xdsSnapshotProducerType int

const (
	proxySourced xdsSnapshotProducerType = iota
	snapshotSourced
)

// syncer.TranslatorSyncerExtension Sync methods are a frequently executed code path.
// We are progressively adding micro-benchmarking to this area of the code to ensure
// we don't introduce unintended latency to this space.
var _ = Describe("ExtAuth Translation - Benchmarking Tests", Label(labels.Performance), func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
	})

	// https://onsi.github.io/ginkgo/#the-spectimeout-and-nodetimeout-decorators
	// It would be preferable to apply a SpecTimeout to these tests to allow them to run for longer than the default
	// I opted to adjust them to fit within the 1-minute window for now, though am leaving this as a note
	// in case we want to expand the test in the future
	DescribeTable("xDS Snapshot Producer",
		func(producerType xdsSnapshotProducerType, resourceFrequency test_fixtures.ResourceFrequency, durationAssertion types.GomegaMatcher) {
			producer := getProducer(producerType)
			desc := generateDesc(producerType, resourceFrequency, durationAssertion)

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
			experiment := gmeasure.NewExperiment(fmt.Sprintf("Experiment - %s", desc))
			AddReportEntry(experiment.Name, experiment)

			experiment.Sample(func(idx int) {
				var before, after runtime.MemStats
				runtime.ReadMemStats(&before)

				// Time ProduceXdsSnapshot
				res, err := gloohelpers.Measure(func() {
					producer.ProduceXdsSnapshot(ctx, settings, snapshot, reports)
				})

				runtime.ReadMemStats(&after)
				var memDiff uint64
				memDiff = (after.TotalAlloc - before.TotalAlloc) / 1000 // KB

				Expect(err).NotTo(HaveOccurred())

				experiment.RecordDuration(desc, res.Total)
				experiment.RecordValue("memory allocated (KB)", float64(memDiff))
			}, gmeasure.SamplingConfig{N: 20})

			durations := experiment.Get(desc).Durations

			Expect(durations).Should(durationAssertion, "Assert ProduceXdsSnapshot meets expected performance targets")
		},
		generateDesc,
		getBenchmarkingTestEntries(),
	)

})

func generateDesc(producerType xdsSnapshotProducerType, rf test_fixtures.ResourceFrequency, _ types.GomegaMatcher) string {
	var typeString string
	switch producerType {
	case proxySourced:
		typeString = "proxySourcedXdsSnapshotProducer"
	case snapshotSourced:
		typeString = "snapshotSourcedXdsSnapshotProducer"
	}

	return fmt.Sprintf("%s (AC=%d, PXY=%d, VH=%d, R=%d)", typeString, rf.AuthConfigs, rf.Proxies, rf.VirtualHostsPerProxy, rf.RoutesPerVirtualHost)
}

func getProducer(producerType xdsSnapshotProducerType) extauthsyncer.XdsSnapshotProducer {
	switch producerType {
	case proxySourced:
		return extauthsyncer.NewProxySourcedXdsSnapshotProducer()
	case snapshotSourced:
		return extauthsyncer.NewSnapshotSourcedXdsSnapshotProducer()
	}
	return nil
}

func getBenchmarkingTestEntries() []TableEntry {
	// As of #5142 we do not expect to run performance tests in cloudbuild
	// These targets are left so that there are appropriate targets if for whatever reason the tests are run in
	// cloudbuild in the future
	if os.Getenv(testutils.GcloudBuildId) != "" {
		// We're running a cloudbuild
		return []TableEntry{
			Entry(nil,
				proxySourced,
				test_fixtures.ResourceFrequency{
					AuthConfigs:          1000,
					Proxies:              50,
					VirtualHostsPerProxy: 100,
					RoutesPerVirtualHost: 1000,
				},
				// Should take 12 +- 2 seconds
				// Recent CI results:
				// 12.230412642
				matchers.HavePercentileWithin(80, 12*time.Second, 2*time.Second),
			),
		}
	}

	if os.Getenv(testutils.GithubAction) != "" {
		// We're running a GHA
		// These targets assume we're running on the 8-core, 32GB RAM runner
		return []TableEntry{
			// (AC=1000, PXY=50, VH=100, R=1000)
			Entry(nil,
				snapshotSourced,
				test_fixtures.ResourceFrequency{
					AuthConfigs:          1000,
					Proxies:              50,
					VirtualHostsPerProxy: 100,
					RoutesPerVirtualHost: 1000,
				},
				// Recent CI results:
				// 1.038826s 1.042333s 1.038221s 1.043138s 1.047517s 1.042832s 1.025864s 1.02505s 1.034961s 1.035683s
				// 1.036788s 1.036823s 1.029266s 1.026531s 1.032527s 1.034193s 1.029518s 1.034541s 1.028277s 1.024386s
				// 1.138405s 1.15853s 1.150084s 1.134609s 1.167894s 1.15503s 1.150983s 1.151669s 1.152028s 1.144216s
				// 1.14433s 1.153812s 1.200387s 1.185184s 1.165962s 1.185261s 1.189101s 1.159826s 1.149311s 1.146693s
				// 1.062801s 1.057882s 1.063839s 1.060852s 1.061825s 1.042442s 1.049299s 1.046706s 1.044362s 1.047164s
				// 1.04883s 1.042133s 1.04047s 1.041431s 1.04804s 1.035113s 1.043892s 1.04273s 1.049597s 1.049994s
				// 1.042341s 1.042542s 1.03826s 1.044233s 1.028609s 1.028543s 1.035026s 1.031438s 1.031079s 1.03304s
				// 1.029206s 1.045725s 1.040016s 1.073661s 1.054951s 1.042765s 1.039053s 1.035061s 1.036125s 1.038115s
				matchers.HavePercentileWithin(80, 1100*time.Millisecond, 300*time.Millisecond),
			),
		}
	}

	// We're running locally
	return []TableEntry{
		// (AC=1000, PXY=50, VH=100, R=1000)
		Entry(nil,
			snapshotSourced,
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
			matchers.HavePercentileWithin(80, 2*time.Second, time.Second),
		),
		Entry(nil,
			proxySourced,
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
			matchers.HavePercentileWithin(80, 7*time.Second, 2*time.Second),
		),
	}
}
