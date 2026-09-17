package reconciler_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/statusutils"
	. "github.com/solo-io/gloo/projects/gateway/pkg/reconciler"
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	validationutils "github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
	"github.com/solo-io/gloo/test/samples"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-kit/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Regression tests for interrupted validation
//
// ReconcileProxies checks its own context rather than relying on the gloo validator to report the
// interruption, so these tests stub the worst case: a validation server that answers successfully
// with the interruption flattened into a listener-level ProcessingError, indistinguishable from a
// rejection.

// interruptedReason is the validator's wrap of runner.ErrValidationInterrupted as a string, the text
// a validation server would leave in a listener report if it flattened the error instead of failing
// the call.
const interruptedReason = "envoy validation of envoy.filters.http.waf config was interrupted: " +
	"envoy validation process was interrupted before it completed (ctxErr=context canceled): " +
	"command \"envoy --mode validate\" failed with error: signal: killed"

var _ = Describe("ReconcileProxies with a context cancelled during validation", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		ns           = "namespace"
		snap         *gloov1snap.ApiSnapshot
		proxy        *gloov1.Proxy
		reports      reporter.ResourceReports
		proxyToWrite GeneratedProxies

		proxyClient  gloov1.ProxyClient
		statusClient resources.StatusClient
	)

	BeforeEach(func() {
		var err error
		ctx, cancel = context.WithCancel(context.Background())

		proxyClient, err = gloov1.NewProxyClient(ctx, &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		})
		Expect(err).NotTo(HaveOccurred())
		statusClient = statusutils.GetStatusClientFromEnvOrDefault(ns)

		snap = samples.SimpleGlooSnapshot(ns)
		tx := translator.NewDefaultTranslator(translator.Opts{WriteNamespace: ns})
		proxy, reports = tx.Translate(ctx, "proxy-name", snap, snap.Gateways)
		proxyToWrite = GeneratedProxies{proxy: reports}
	})

	AfterEach(func() {
		cancel()
	})

	// interruptedValidation reports the proxy's first listener as errored with the interruption text,
	// optionally cancelling the reconcile context first.
	interruptedValidation := func(cancelDuringValidation bool) func(
		context.Context, *validation.GlooValidationServiceRequest,
	) (*validation.GlooValidationServiceResponse, error) {
		return func(_ context.Context, req *validation.GlooValidationServiceRequest) (*validation.GlooValidationServiceResponse, error) {
			if cancelDuringValidation {
				cancel()
			}
			report := validationutils.MakeReport(req.GetProxy())
			validationutils.AppendListenerError(
				report.GetListenerReports()[0],
				validation.ListenerReport_Error_ProcessingError,
				interruptedReason,
			)
			return &validation.GlooValidationServiceResponse{
				ValidationReports: []*validation.ValidationReport{{ProxyReport: report}},
			}, nil
		}
	}

	It("aborts without writing a proxy derived from the interrupted validation", func() {
		reconciler := NewProxyReconciler(interruptedValidation(true), proxyClient, statusClient)

		err := reconciler.ReconcileProxies(ctx, proxyToWrite, ns, clients.ListOpts{Selector: map[string]string{}})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("context ended before or during proxy validation"))

		_, err = proxyClient.Read(ns, proxy.GetMetadata().GetName(), clients.ReadOpts{})
		Expect(errors.IsNotExist(err)).To(BeTrue(), "no proxy should be written from an aborted reconcile")
	})

	// The same validation result under a live context is a genuine rejection and must still strip.
	// This proves the context check, not the report content, prevents the write above.
	It("still strips a listener reported as errored when the context is live, and logs why", func() {
		observerCore, observedLogs := observer.New(zapcore.WarnLevel)
		ctx := contextutils.WithExistingLogger(ctx, zap.New(observerCore).Sugar())

		reconciler := NewProxyReconciler(interruptedValidation(false), proxyClient, statusClient)

		err := reconciler.ReconcileProxies(ctx, proxyToWrite, ns, clients.ListOpts{Selector: map[string]string{}})
		Expect(err).NotTo(HaveOccurred())

		written, err := proxyClient.Read(ns, proxy.GetMetadata().GetName(), clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(written.GetListeners()).To(HaveLen(len(proxy.GetListeners()) - 1))

		warnings := observedLogs.FilterMessage("stripping invalid listener from proxy").All()
		Expect(warnings).To(HaveLen(1))

		// Assert against the encoded line, since what matters is how the fields render in a log.
		encoded, encodeErr := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()).
			EncodeEntry(warnings[0].Entry, warnings[0].Context)
		Expect(encodeErr).NotTo(HaveOccurred())
		line := encoded.String()

		Expect(line).To(ContainSubstring(`"rejectedResources":[{"kind":"Gateway","ref":"namespace.gateway-proxy"`),
			"the warning should name the rejected resource by the kind and ref the user wrote")
		Expect(line).To(ContainSubstring("was interrupted"),
			"the warning should carry the errors from the rejected resource")
		Expect(line).To(ContainSubstring(`"proxy":"namespace.proxy-name"`),
			"refs should render as plain namespace.name, not as prototext")
		Expect(line).NotTo(ContainSubstring("error occurred"),
			"the multierror preamble should be flattened into one message per error")
		Expect(line).NotTo(ContainSubstring(`\n`),
			"no field should carry escaped newlines into the log line")
	})
})
