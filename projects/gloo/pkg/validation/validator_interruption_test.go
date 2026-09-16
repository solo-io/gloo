package validation_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/envoyinit/pkg/runner"
	validationgrpc "github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	glooplugins "github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	pluginregistry "github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
	glootranslator "github.com/solo-io/gloo/projects/gloo/pkg/translator"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	gloovalidation "github.com/solo-io/gloo/projects/gloo/pkg/validation"
	testsamples "github.com/solo-io/gloo/test/samples"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	skreporter "github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

// Regression tests for interrupted validation
//
// These specs run a real translator so they exercise the actual path a plugin error takes into a
// report.
var _ = Describe("GlooValidator with an interrupted envoy validation", func() {

	var (
		ctx       context.Context
		snap      *gloov1snap.ApiSnapshot
		sanitizer *recordingSanitizer
	)

	BeforeEach(func() {
		ctx = context.Background()
		snap = testsamples.SimpleGlooSnapshot("gloo-system")
		sanitizer = &recordingSanitizer{}
	})

	// newValidator builds a gloo validator whose only http filter plugin fails with the given error.
	newValidator := func(pluginErr error) gloovalidation.GlooValidator {
		settings := &gloov1.Settings{}
		registry := pluginregistry.NewPluginRegistry([]glooplugins.Plugin{
			&failingFilterPlugin{err: pluginErr},
		})

		return gloovalidation.NewGlooValidator(gloovalidation.GlooValidatorConfig{
			Translator: glootranslator.NewTranslatorWithHasher(
				glooutils.NewSslConfigTranslator(), settings, registry,
				glootranslator.EnvoyCacheResourcesListToFnvHash),
			XdsSanitizer: sanitizer,
			Settings:     settings,
		})
	}

	It("fails the call instead of returning reports from the interrupted run", func() {
		interrupted := eris.Wrapf(runner.ErrValidationInterrupted,
			"envoy validation of io.solo.filters.http.modsecurity config was interrupted")

		reports, err := newValidator(interrupted).Validate(ctx, snap.Proxies[0], snap, false)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("envoy config validation was interrupted"))
		Expect(reports).To(BeEmpty(),
			"reports from an interrupted run carry no validation information and must not reach a caller that acts on them as a validation result")
	})

	It("does not sanitize an interrupted translation", func() {
		interrupted := eris.Wrapf(runner.ErrValidationInterrupted, "interrupted")

		_, err := newValidator(interrupted).Validate(ctx, snap.Proxies[0], snap, false)

		Expect(err).To(HaveOccurred())
		// The route-replacing sanitizer acts on errored route reports, so it must not see reports from
		// an interrupted run.
		Expect(sanitizer.called).To(BeFalse(), "the guard must precede sanitization")
	})

	It("fails the call at the gRPC boundary rather than returning reports", func() {
		interrupted := eris.Wrapf(runner.ErrValidationInterrupted,
			"envoy validation of io.solo.filters.http.modsecurity config was interrupted")
		settings := &gloov1.Settings{}
		registry := pluginregistry.NewPluginRegistry([]glooplugins.Plugin{
			&failingFilterPlugin{err: interrupted},
		})
		server := gloovalidation.NewValidator(gloovalidation.ValidatorConfig{
			Ctx: ctx,
			GlooValidatorConfig: gloovalidation.GlooValidatorConfig{
				Translator: glootranslator.NewTranslatorWithHasher(
					glooutils.NewSslConfigTranslator(), settings, registry,
					glootranslator.EnvoyCacheResourcesListToFnvHash),
				XdsSanitizer: sanitizer,
				Settings:     settings,
			},
		})
		Expect(server.Sync(ctx, snap)).To(Succeed())

		resp, err := server.Validate(ctx, &validationgrpc.GlooValidationServiceRequest{Proxy: snap.Proxies[0]})
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil(),
			"reports from an interrupted run carry no validation information and must not reach a caller that acts on them as a validation result")
		Expect(err.Error()).To(ContainSubstring("was interrupted"),
			"the error should carry the interruption cause")
	})

	// Only an interruption fails the call. A genuine rejection must still come back as a report.
	It("still reports a genuine plugin rejection through the reports, with no call error", func() {
		reports, err := newValidator(eris.New("plugin rejected the config")).
			Validate(ctx, snap.Proxies[0], snap, false)

		Expect(err).NotTo(HaveOccurred())
		Expect(reports).To(HaveLen(1))
		Expect(reports[0].ResourceReports.ValidateStrict()).To(HaveOccurred(),
			"a real rejection must still be reported so the resource is rejected")
		Expect(sanitizer.called).To(BeTrue())
	})
})

// failingFilterPlugin is an http filter plugin that always fails with err.
type failingFilterPlugin struct {
	err error
}

func (p *failingFilterPlugin) Name() string { return "failing-filter-plugin" }

func (p *failingFilterPlugin) Init(glooplugins.InitParams) {}

func (p *failingFilterPlugin) HttpFilters(glooplugins.Params, *gloov1.HttpListener) ([]glooplugins.StagedHttpFilter, error) {
	return nil, p.err
}

// recordingSanitizer records whether it ran.
type recordingSanitizer struct {
	called bool
}

func (s *recordingSanitizer) SanitizeSnapshot(
	_ context.Context,
	_ *gloov1snap.ApiSnapshot,
	xdsSnapshot envoycache.Snapshot,
	_ skreporter.ResourceReports,
) envoycache.Snapshot {
	s.called = true
	return xdsSnapshot
}
