package k8sadmission

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"slices"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	"github.com/solo-io/gloo/pkg/utils/statsutils"

	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"

	"github.com/solo-io/gloo/pkg/utils/namespaces"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	errors "github.com/rotisserie/eris"
	gwv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/validation"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/api/external/kubernetes/namespace"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	ValidationPath      = "/validation"
	SkipValidationKey   = "gateway.solo.io/skip_validation"
	SkipValidationValue = "true"
	// kubernetesCoreApiGroup is the GVK group name for Kubernetes core API resources
	kubernetesCoreApiGroup = ""
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	resourceTypeKey, _ = tag.NewKey("resource_type")
	resourceRefKey, _  = tag.NewKey("resource_ref")

	mGatewayResourcesAccepted = statsutils.MakeSumCounter("validation.gateway.solo.io/resources_accepted", "The number of resources accepted")
	mGatewayResourcesRejected = statsutils.MakeSumCounter("validation.gateway.solo.io/resources_rejected", "The number of resources rejected")

	unmarshalErrMsg     = "could not unmarshal raw object"
	UnmarshalErr        = errors.New(unmarshalErrMsg)
	WrappedUnmarshalErr = func(err error) error {
		return errors.Wrap(err, unmarshalErrMsg)
	}
	ListGVK = schema.GroupVersionKind{
		Version: "v1",
		Group:   "",
		Kind:    "List",
	}
)

const (
	ApplicationJson = "application/json"
	ApplicationYaml = "application/x-yaml"
)

func incrementMetric(ctx context.Context, resource string, ref *core.ResourceRef, m *stats.Int64Measure) {
	statsutils.MeasureOne(
		ctx,
		m,
		tag.Insert(resourceTypeKey, resource),
		tag.Insert(resourceRefKey, fmt.Sprintf("%v.%v", ref.GetNamespace(), ref.GetName())),
	)
}

func skipValidationCheck(annotations map[string]string) bool {
	if annotations == nil {
		return false
	}
	return annotations[SkipValidationKey] == SkipValidationValue
}

type WebhookConfig struct {
	ctx                           context.Context
	validator                     validation.Validator
	watchNamespaces               []string
	port                          int
	serverCertPath, serverKeyPath string
	alwaysAccept                  bool // accept all resources
	readGatewaysFromAllNamespaces bool
	webhookNamespace              string
	kubeGatewayEnabled            bool
}

func NewWebhookConfig(
	ctx context.Context,
	validator validation.Validator,
	watchNamespaces []string,
	port int,
	serverCertPath, serverKeyPath string,
	alwaysAccept, readGatewaysFromAllNamespaces bool,
	webhookNamespace string,
	kubeGatewayEnabled bool,
) WebhookConfig {
	return WebhookConfig{
		ctx:                           ctx,
		validator:                     validator,
		watchNamespaces:               watchNamespaces,
		port:                          port,
		serverCertPath:                serverCertPath,
		serverKeyPath:                 serverKeyPath,
		alwaysAccept:                  alwaysAccept,
		readGatewaysFromAllNamespaces: readGatewaysFromAllNamespaces,
		webhookNamespace:              webhookNamespace,
		kubeGatewayEnabled:            kubeGatewayEnabled,
	}
}

func NewGatewayValidatingWebhook(cfg WebhookConfig) (*http.Server, error) {
	ctx := cfg.ctx
	validator := cfg.validator
	watchNamespaces := cfg.watchNamespaces
	port := cfg.port
	serverCertPath := cfg.serverCertPath
	serverKeyPath := cfg.serverKeyPath
	alwaysAccept := cfg.alwaysAccept
	readGatewaysFromAllNamespaces := cfg.readGatewaysFromAllNamespaces
	webhookNamespace := cfg.webhookNamespace
	kubeGatewayEnabled := cfg.kubeGatewayEnabled

	certProvider, err := NewCertificateProvider(serverCertPath, serverKeyPath, log.New(&debugLogger{ctx: ctx}, "validation-webhook-certificate-watcher", log.LstdFlags), ctx, 10*time.Second)
	if err != nil {
		return nil, errors.Wrapf(err, "loading TLS certificate provider")
	}

	handler := NewGatewayValidationHandler(
		contextutils.WithLogger(ctx, "gateway-validation-webhook"),
		validator,
		watchNamespaces,
		alwaysAccept,
		readGatewaysFromAllNamespaces,
		webhookNamespace,
		kubeGatewayEnabled,
	)

	mux := http.NewServeMux()
	mux.Handle(ValidationPath, handler)

	return &http.Server{
		Addr:      fmt.Sprintf(":%v", port),
		TLSConfig: &tls.Config{GetCertificate: certProvider.GetCertificateFunc()},
		Handler:   mux,
		ErrorLog:  log.New(&debugLogger{ctx: ctx}, "validation-webhook-server", log.LstdFlags),
	}, nil
}

type debugLogger struct{ ctx context.Context }

func (l *debugLogger) Write(p []byte) (n int, err error) {
	contextutils.LoggerFrom(l.ctx).Debug(string(p))
	return len(p), nil
}

type gatewayValidationWebhook struct {
	// everything living here MUST be synchronized to avoid data races
	ctx                           context.Context
	validator                     validation.Validator // the function calls are synchronized
	watchNamespaces               []string             // we make a deep copy when we read this; original is read only so no races
	alwaysAccept                  bool                 // read only so no races
	readGatewaysFromAllNamespaces bool                 // read only so no races
	webhookNamespace              string               // read only so no races
	kubeGatewayEnabled            bool                 // read only so no races
}

type AdmissionReviewWithProxies struct {
	AdmissionRequestWithProxies
	AdmissionResponseWithProxies
}

type AdmissionRequestWithProxies struct {
	v1beta1.AdmissionReview
	ReturnProxies bool `json:"returnProxies,omitempty"`
}

// Validation webhook works properly even if extra fields are provided in the response
type AdmissionResponseWithProxies struct {
	*v1beta1.AdmissionResponse
	Proxies []*gloov1.Proxy `json:"proxies,omitempty"`
}

func NewGatewayValidationHandler(
	ctx context.Context,
	validator validation.Validator,
	watchNamespaces []string,
	alwaysAccept bool,
	readGatewaysFromAllNamespaces bool,
	webhookNamespace string,
	kubeGatewayEnabled bool,
) *gatewayValidationWebhook {
	return &gatewayValidationWebhook{
		ctx:                           ctx,
		validator:                     validator,
		watchNamespaces:               watchNamespaces,
		alwaysAccept:                  alwaysAccept,
		readGatewaysFromAllNamespaces: readGatewaysFromAllNamespaces,
		webhookNamespace:              webhookNamespace,
		kubeGatewayEnabled:            kubeGatewayEnabled,
	}
}

func (wh *gatewayValidationWebhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := contextutils.LoggerFrom(wh.ctx)

	logger.Debug("received validation request on webhook")

	b, _ := httputil.DumpRequest(r, true)
	logger.Debugf("validation request dump:\n %s", string(b))

	// Verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != ApplicationJson && contentType != ApplicationYaml {
		logger.Errorf("contentType=%s, expecting application/json or application/x-yaml", contentType)
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
		defer r.Body.Close()
	}
	if len(body) == 0 {
		logger.Errorf("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	var (
		admissionResponse = &AdmissionResponseWithProxies{}
		review            AdmissionReviewWithProxies
		err               error
	)

	if contentType == ApplicationYaml {
		if err = yaml.Unmarshal(body, &review); err == nil {
			admissionResponse = wh.makeAdmissionResponse(wh.ctx, &review)
		}
	} else {
		if _, _, err := deserializer.Decode(body, nil, &review); err == nil {
			admissionResponse = wh.makeAdmissionResponse(wh.ctx, &review)
		}
	}

	if err != nil {
		logger.Errorf("Can't decode body: %v", err)
		admissionResponse.AdmissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	admissionReview := AdmissionReviewWithProxies{}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse.AdmissionResponse
		if review.ReturnProxies {
			admissionReview.Proxies = admissionResponse.Proxies
		}
		if review.Request != nil {
			admissionReview.Response.UID = review.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		logger.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}
	logger.Debug("ready to write response ...")
	if _, err := w.Write(resp); err != nil {
		logger.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
	logger.Debugf("responded with review: %s", resp)
}

func (wh *gatewayValidationWebhook) makeAdmissionResponse(ctx context.Context, review *AdmissionReviewWithProxies) *AdmissionResponseWithProxies {
	oldLogger := contextutils.LoggerFrom(ctx)
	req := review.Request

	logger := oldLogger.With("Kind", req.Kind,
		"Namespace", req.Namespace,
		"Name", req.Name,
		"UID", req.UID,
		"PatchOperation", req.Operation,
		"UserInfo", req.UserInfo,
	)
	ctx = contextutils.WithExistingLogger(ctx, logger)

	logger.Debugf("Start AdmissionReview")

	gvk := schema.GroupVersionKind{
		Group:   req.Kind.Group,
		Version: req.Kind.Version,
		Kind:    req.Kind.Kind,
	}

	watchNamespaces := make([]string, len(wh.watchNamespaces))
	copy(watchNamespaces, wh.watchNamespaces) // important we make a deep copy

	if req.Kind.Kind == "Namespace" {
		// if it's not a namespace we watch, do not validate
		if !slices.Contains(watchNamespaces, req.Name) {
			return &AdmissionResponseWithProxies{
				AdmissionResponse: &v1beta1.AdmissionResponse{
					Allowed: true,
				},
			}
		}
	} else {
		// If we've specified to NOT read gateway requests from all namespaces, then only
		// check gateway requests for the same namespace as this webhook, regardless of the
		// contents of watchNamespaces. It's assumed that if it's non-empty, watchNamespaces
		// contains the webhook's own namespace, since this was checked during setup in setup_syncer.go
		if gvk == gwv1.GatewayGVK && !wh.readGatewaysFromAllNamespaces && !namespaces.AllNamespaces(wh.watchNamespaces) {
			watchNamespaces = []string{wh.webhookNamespace}
		}

		// ensure the request applies to a watched namespace, if watchNamespaces is set
		var validatingForNamespace bool
		if len(watchNamespaces) > 0 {
			for _, ns := range watchNamespaces {
				if ns == metav1.NamespaceAll || ns == req.Namespace {
					validatingForNamespace = true
					break
				}
			}
		} else {
			validatingForNamespace = true
		}

		// if it's not our namespace, do not validate
		if !validatingForNamespace {
			return &AdmissionResponseWithProxies{
				AdmissionResponse: &v1beta1.AdmissionResponse{
					Allowed: true,
				},
			}
		}
	}

	ref := &core.ResourceRef{
		Namespace: req.Namespace,
		Name:      req.Name,
	}

	reports, validationErrs := wh.validateAdmissionRequest(ctx, gvk, ref, req)

	hasUnmarshalErr := false
	if validationErrs != nil {
		for _, e := range validationErrs.Errors {
			if errors.Is(e, UnmarshalErr) {
				hasUnmarshalErr = true
			}
		}
	}

	// even if validation is set to always accept, we want to fail on unmarshal errors
	if validationErrs.ErrorOrNil() == nil || (wh.alwaysAccept && !hasUnmarshalErr) {
		logger.Debugf("Succeeded, alwaysAccept: %v validationErrs: %v", wh.alwaysAccept, validationErrs)
		incrementMetric(ctx, gvk.String(), ref, mGatewayResourcesAccepted)
		return &AdmissionResponseWithProxies{
			AdmissionResponse: &v1beta1.AdmissionResponse{
				Allowed: true,
			},
			Proxies: reports.GetProxies(),
		}
	}

	incrementMetric(ctx, gvk.String(), ref, mGatewayResourcesRejected)
	logger.Errorf("Validation failed: %v", validationErrs)

	finalErr := errors.Errorf("resource incompatible with current Gloo snapshot: %v", validationErrs.Errors)
	details := &metav1.StatusDetails{
		Name:   req.Name,
		Group:  gvk.Group,
		Kind:   gvk.Kind,
		Causes: getFailureCauses(validationErrs),
	}

	return &AdmissionResponseWithProxies{
		AdmissionResponse: &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: finalErr.Error(),
				Details: details,
			},
		},
		Proxies: reports.GetProxies(),
	}
}

func getFailureCauses(validationErr *multierror.Error) []metav1.StatusCause {
	var causes []metav1.StatusCause
	for _, e := range validationErr.Errors {
		causes = append(causes, metav1.StatusCause{
			Message: fmt.Sprintf("Error %v", e.Error()),
		})
	}
	return causes
}

func (wh *gatewayValidationWebhook) validateAdmissionRequest(
	ctx context.Context,
	gvk schema.GroupVersionKind,
	ref *core.ResourceRef,
	admissionRequest *v1beta1.AdmissionRequest,
) (*validation.Reports, *multierror.Error) {

	isDelete := admissionRequest.Operation == v1beta1.Delete
	isUpdate := admissionRequest.Operation == v1beta1.Update
	dryRun := isDryRun(admissionRequest)

	if gvk == ListGVK {
		return wh.validateList(ctx, admissionRequest.Object.Raw, dryRun)
	}

	// Kubernetes Secrets and Namespace are the only non-Solo API resources we support validation requests for.
	// For a these resources, we want to skip validation on operations we don't support. For Namespaces, we
	// support DELETE and UPDATE. For Secrets, we support DELETE.
	// Else, we expect to find the resource in our ApiGvkToHashableResource map - if the resource is supported.

	if gvk.Group == kubernetesCoreApiGroup && gvk.Kind == "Secret" {
		if !isDelete {
			contextutils.LoggerFrom(ctx).Infof("unsupported operation validation [%s] for resource namespace [%s] name [%s] group [%s] kind [%s]", admissionRequest.Operation, ref.GetNamespace(), ref.GetName(), gvk.Group, gvk.Kind)
			return &validation.Reports{}, nil
		}
	} else if gvk.Group == kubernetesCoreApiGroup && gvk.Kind == "Namespace" {
		if !isDelete && !isUpdate {
			contextutils.LoggerFrom(ctx).Infof("unsupported operation validation [%s] for resource namespace [%s] name [%s] group [%s] kind [%s]", admissionRequest.Operation, ref.GetNamespace(), ref.GetName(), gvk.Group, gvk.Kind)
			return &validation.Reports{}, nil
		}

		if isUpdate {
			var namespace namespace.KubeNamespace
			err := json.Unmarshal(admissionRequest.Object.Raw, &namespace)
			if err != nil {
				return nil, &multierror.Error{Errors: []error{err}}
			}
			kns := kubernetes.KubeNamespace{
				KubeNamespace: namespace,
			}
			// At this point, any namespace update we receive is a result of modifying a namespace we watch
			// Check if we still watch the namespace after the modification
			namespaceStillWatched, err := settingsutil.NamespaceWatched(settingsutil.FromContext(ctx), kns)
			if err != nil {
				return nil, &multierror.Error{Errors: []error{err}}
			}

			if namespaceStillWatched {
				// The namespace has changed in a way that does not affect us so we let it through
				return &validation.Reports{}, nil
			} else {
				// The namespace used to be watched but will now be removed from the list of namespaces we watch
				// Run validation on this updated list of namespaces and report if anything breaks
				return wh.deleteRef(ctx, gvk, ref, admissionRequest)
			}

		}
	} else if _, hit := gloosnapshot.ApiGvkToHashableResource[gvk]; !hit {
		contextutils.LoggerFrom(ctx).Infof("unsupported validation for resource namespace [%s] name [%s] group [%s] kind [%s]", ref.GetNamespace(), ref.GetName(), gvk.Group, gvk.Kind)
		return &validation.Reports{}, nil
	}

	if isDelete {
		return wh.deleteRef(ctx, gvk, ref, admissionRequest)
	}

	return wh.validateGvk(ctx, gvk, ref, admissionRequest)
}

func (wh *gatewayValidationWebhook) deleteRef(ctx context.Context, gvk schema.GroupVersionKind, ref *core.ResourceRef, admissionRequest *v1beta1.AdmissionRequest) (*validation.Reports, *multierror.Error) {
	newResourceFunc := gloosnapshot.ApiGvkToHashableResource[gvk]
	// Special case for Kubernetes secrets and namespaces, since they are not handled by our hashable resource.
	// We can reuse the NewSecretHashableResource resource.Resource, since all that matters is the metadata for deletion.
	if gvk.Group == kubernetesCoreApiGroup && gvk.Kind == "Secret" {
		newResourceFunc = gloov1.NewSecretHashableResource
	}
	if gvk.Group == kubernetesCoreApiGroup && gvk.Kind == "Namespace" {
		newResourceFunc = kubernetes.NewKubeNamespaceHashableResource
	}
	newResource := newResourceFunc()
	newResource.SetMetadata(&core.Metadata{
		Namespace: ref.GetNamespace(),
		Name:      ref.GetName(),
	})
	err := wh.validator.ValidateDeletedGvk(ctx, gvk, newResource, isDryRun(admissionRequest))
	if err != nil {
		return nil, &multierror.Error{Errors: []error{errors.Wrapf(err, "failed validating the deletion of resource namespace: %s name: %s", ref.GetNamespace(), ref.GetName())}}
	}
	return &validation.Reports{}, nil
}

func (wh *gatewayValidationWebhook) validateGvk(ctx context.Context, gvk schema.GroupVersionKind, _ *core.ResourceRef, admissionRequest *v1beta1.AdmissionRequest) (*validation.Reports, *multierror.Error) {
	var reports *validation.Reports
	newResourceFunc := gloosnapshot.ApiGvkToHashableResource[gvk]

	// check to see if we should perform validation for this GVK
	if !wh.shouldValidateGvk(gvk) {
		return nil, nil
	}

	newResource := newResourceFunc()
	oldResource := newResourceFunc()

	shouldValidate, shouldValidateErr := wh.shouldValidateResource(ctx, admissionRequest, newResource, oldResource)
	if shouldValidateErr != nil {
		return nil, &multierror.Error{Errors: []error{shouldValidateErr}}
	}
	if !shouldValidate {
		return nil, nil
	}

	reports, err := wh.validator.ValidateModifiedGvk(ctx, gvk, newResource, isDryRun(admissionRequest))
	if err != nil {
		return reports, &multierror.Error{Errors: []error{errors.Wrapf(err, "Validating %T failed", newResource)}}
	}
	return reports, nil
}

func (wh *gatewayValidationWebhook) validateList(ctx context.Context, rawJson []byte, dryRun bool) (*validation.Reports, *multierror.Error) {
	var (
		ul      unstructured.UnstructuredList
		reports *validation.Reports
		errs    *multierror.Error
	)

	if err := ul.UnmarshalJSON(rawJson); err != nil {
		return nil, &multierror.Error{Errors: []error{WrappedUnmarshalErr(err)}}
	}
	if reports, errs = wh.validator.ValidateList(ctx, &ul, dryRun); errs != nil {
		return reports, errs
	}
	return reports, nil
}

func (wh *gatewayValidationWebhook) shouldValidateGvk(gvk schema.GroupVersionKind) bool {
	if gvk == gwv1.RouteOptionGVK || gvk == gwv1.VirtualHostOptionGVK {
		// only validate RouteOption and VirtualHostOption resources if K8s Gateway is enabled
		return wh.kubeGatewayEnabled
	}

	// no other special considerations at this point, so continue with validation
	return true
}

// shouldValidateResource determines if a resource should be validated AND populates `resource` (and `oldResource` if applicable) from the objects[s] in the `admissionRequest`
func (wh *gatewayValidationWebhook) shouldValidateResource(ctx context.Context, admissionRequest *v1beta1.AdmissionRequest, resource, oldResource resources.HashableResource) (bool, error) {
	logger := contextutils.LoggerFrom(ctx)

	if err := validation.UnmarshalResource(admissionRequest.Object.Raw, resource); err != nil {
		return false, &multierror.Error{Errors: []error{WrappedUnmarshalErr(err)}}
	}
	if skipValidationCheck(resource.GetMetadata().GetAnnotations()) {
		logger.Debugf("Skipping validation. Reason: detected skip validation annotation")
		return false, nil
	}

	if admissionRequest.Operation != v1beta1.Update {
		return true, nil
	}

	// For update requests, we check to see if this is a status update
	// If it is, we do not need to validate the resource
	if err := validation.UnmarshalResource(admissionRequest.OldObject.Raw, oldResource); err != nil {
		return false, &multierror.Error{Errors: []error{WrappedUnmarshalErr(err)}}
	}

	equalResources := resource.MustHash() == oldResource.MustHash()
	if equalResources {
		logger.Debugf("Skipping validation. Reason: status only update")
		return false, nil
	}
	return true, nil
}

func isDryRun(admissionRequest *v1beta1.AdmissionRequest) bool {
	if admissionRequest.DryRun != nil {
		return *admissionRequest.DryRun
	}
	return false
}
