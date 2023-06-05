package translator_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
	mock_consul "github.com/solo-io/gloo/projects/gloo/pkg/upstreams/consul/mocks"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	gloohelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	. "github.com/solo-io/gloo/projects/gloo/pkg/translator"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	validationutils "github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
	"github.com/solo-io/gloo/test/ginkgo/labels"
	"sort"
	"time"
)

type benchmarkEntry struct {
	// Name for this test
	desc string

	// The snapshot to translate
	snapshot *v1snap.ApiSnapshot

	// Configuration for the benchmarking
	tries          int
	maxDur         time.Duration
	benchmarkFuncs []benchmarkFunc
}

type benchmarkFunc func(durations []time.Duration)

var _ = FDescribe("Translation - Benchmarking Tests", Serial, Label(labels.Performance), func() {
	var (
		ctrl       *gomock.Controller
		settings   *v1.Settings
		translator Translator

		registeredPlugins []plugins.Plugin
	)

	BeforeEach(func() {

		ctrl = gomock.NewController(T)

		settings = &v1.Settings{}
		memoryClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		opts := bootstrap.Opts{
			Settings:  settings,
			Secrets:   memoryClientFactory,
			Upstreams: memoryClientFactory,
			Consul: bootstrap.Consul{
				ConsulWatcher: mock_consul.NewMockConsulWatcher(ctrl), // just needed to activate the consul plugin
			},
		}
		registeredPlugins = registry.Plugins(opts)

		pluginRegistry := registry.NewPluginRegistry(registeredPlugins)

		translator = NewTranslatorWithHasher(glooutils.NewSslConfigTranslator(), settings, pluginRegistry, EnvoyCacheResourcesListToFnvHash)
	})

	DescribeTable("Translate",
		func(ent benchmarkEntry) {

			params := plugins.Params{
				Ctx:      context.Background(),
				Snapshot: ent.snapshot,
			}

			var (
				snap   cache.Snapshot
				errs   reporter.ResourceReports
				report *validation.ProxyReport
			)

			experiment := gmeasure.NewExperiment("translate")

			AddReportEntry(experiment.Name, experiment)

			statName := fmt.Sprintf("translating %s", ent.desc)
			experiment.Sample(func(idx int) {

				// Time translation
				experiment.MeasureDuration(statName, func() {
					snap, errs, report = translator.Translate(params, gloohelpers.Proxy())
				})

				// Assert expected results
				Expect(errs.Validate()).NotTo(HaveOccurred())
				Expect(snap).NotTo(BeNil())
				Expect(report).To(Equal(validationutils.MakeReport(gloohelpers.Proxy())))
			}, gmeasure.SamplingConfig{N: ent.tries, Duration: ent.maxDur})

			durations := experiment.Get(statName).Durations

			for _, bench := range ent.benchmarkFuncs {
				bench(durations)
			}
		},
		Entry("basic", basicCase),
		Entry("10 upstreams", upstreamScale(10)),
		Entry("100 upstreams", upstreamScale(100)),
		Entry("1000 upstreams", upstreamScale(1000)),
		Entry("10 endpoints", endpointScale(10)),
		Entry("100 endpoints", endpointScale(100)),
		Entry("1000 endpoints", endpointScale(1000)),
	)
})

var basicCase = benchmarkEntry{
	desc: "basic",
	snapshot: &v1snap.ApiSnapshot{
		Endpoints: []*v1.Endpoint{gloohelpers.Endpoint},
		Upstreams: []*v1.Upstream{gloohelpers.Upstream},
	},
	tries:  1000,
	maxDur: time.Second,
	benchmarkFuncs: []benchmarkFunc{
		median(5 * time.Millisecond),
		percentile(90, 10*time.Millisecond),
	},
}

var upstreamScale = func(numUpstreams int) benchmarkEntry {
	return benchmarkEntry{
		desc: fmt.Sprintf("%d upstreams", numUpstreams),
		snapshot: gloohelpers.ScaledSnapshot(gloohelpers.ScaleConfig{
			Upstreams: numUpstreams,
			Endpoints: 1,
		}),
		tries:  1000,
		maxDur: time.Second,
		benchmarkFuncs: []benchmarkFunc{
			median(5 * time.Millisecond),
			percentile(90, 10*time.Millisecond),
		},
	}
}

var endpointScale = func(numEndpoints int) benchmarkEntry {
	return benchmarkEntry{
		desc: fmt.Sprintf("%d endpoints", numEndpoints),
		snapshot: gloohelpers.ScaledSnapshot(gloohelpers.ScaleConfig{
			Endpoints: numEndpoints,
			Upstreams: 1,
		}),
		tries:  1000,
		maxDur: time.Second,
		benchmarkFuncs: []benchmarkFunc{
			median(5 * time.Millisecond),
			percentile(90, 10*time.Millisecond),
		},
	}
}

var median = func(target time.Duration) benchmarkFunc {
	return func(durations []time.Duration) {
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		var median time.Duration
		if l := len(durations); l%2 == 1 {
			median = durations[l/2]
		} else {
			median = (durations[l/2] + durations[l/2-1]) / 2
		}
		Expect(median).To(BeNumerically("<", target))
	}
}

var percentile = func(percentile int, target time.Duration) benchmarkFunc {
	return func(durations []time.Duration) {
		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
		pct := durations[int(float64(len(durations))*(float64(percentile)/float64(100)))]
		Expect(pct).To(BeNumerically("<", target))
	}
}
