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
				// 1.338306s 1.553647s 1.091003s 1.260452s 1.512719s 1.191144s 1.37824s 1.436372s 1.174521s 1.191923s
				// 1.295572s 1.122904s 1.031022s 1.034975s 1.01578s 1.022888s 1.038558s 1.027296s 1.02098s 1.026748s
				// 1.587027s 1.508666s 1.706318s 1.132044s 1.586973s 1.533919s 1.263382s 1.067526s 1.482481s 1.146513s
				// 1.019298s 1.055201s 1.022608s 1.004021s 1.015732s 1.004964s 1.041359s 1.061431s 1.009227s 1.006839s
				// 1.894277s 1.694374s 1.28535s 1.830699s 1.748889s 1.504626s 1.43577s 1.518958s 1.251122s 1.236105s
				// 1.274074s 1.211417s 1.207307s 1.198589s 1.19741s 1.224671s 1.200129s 1.204121s 1.201809s 1.202677s
				// 1.465857s 1.542869s 1.075749s 1.177539s 1.288051s 1.575071s 1.680141s 1.055014s 1.235321s 1.503519s
				// 1.040999s 1.04226s 1.071544s 1.027743s 1.057385s 1.027545s 1.036268s 1.026675s 1.026807s 1.024249s
				matchers.HavePercentileWithin(80, 1500*time.Millisecond, 500*time.Millisecond),
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
