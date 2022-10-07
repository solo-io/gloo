package validation

import (
	"context"
	"sync"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/go-utils/hashutils"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	errors "github.com/rotisserie/eris"
	utils2 "github.com/solo-io/gloo/pkg/utils"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	"github.com/solo-io/gloo/projects/gateway/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloo_translator "github.com/solo-io/gloo/projects/gloo/pkg/translator"
	validationutils "github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	skprotoutils "github.com/solo-io/solo-kit/pkg/utils/protoutils"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Reports struct {
	Proxies         []*gloov1.Proxy
	ProxyReports    *ProxyReports
	UpstreamReports *UpstreamReports
}

func (r *Reports) GetProxies() []*gloov1.Proxy {
	if r == nil || r.Proxies == nil {
		return []*gloov1.Proxy{}
	}
	return r.Proxies
}

type ProxyReports []*validation.ProxyReport
type UpstreamReports []*validation.ResourceReport

var (
	NotReadyErr = errors.Errorf("validation is not yet available. Waiting for first snapshot")

	RouteTableDeleteErr = func(parentVirtualServices, parentRouteTables []*core.ResourceRef) error {
		return errors.Errorf("Deletion blocked because active Routes delegate to this Route Table. Remove delegate actions to this route table from the virtual services: %v and the route tables: %v, then try again",
			parentVirtualServices,
			parentRouteTables)
	}
	VirtualServiceDeleteErr = func(parentGateways []*core.ResourceRef) error {
		return errors.Errorf("Deletion blocked because active Gateways reference this Virtual Service. Remove refs to this virtual service from the gateways: %v, then try again",
			parentGateways)
	}
	unmarshalErrMsg     = "could not unmarshal raw object"
	WrappedUnmarshalErr = func(err error) error {
		return errors.Wrapf(err, unmarshalErrMsg)
	}

	GlooValidationResponseLengthError = func(resp *validation.GlooValidationServiceResponse) error {
		return errors.Errorf("Expected Gloo validation response to contain 1 report, but contained %d",
			len(resp.GetValidationReports()))
	}

	mValidConfig = utils2.MakeGauge("validation.gateway.solo.io/valid_config",
		"A boolean indicating whether gloo config is valid")
)

const (
	InvalidSnapshotErrMessage = "validation is disabled due to an invalid resource which has been written to storage. " +
		"Please correct any Rejected resources to re-enable validation."
)

var _ Validator = &validator{}

type Validator interface {
	gloov1snap.ApiSyncer
	ValidateList(ctx context.Context, ul *unstructured.UnstructuredList, dryRun bool) (*Reports, *multierror.Error)
	ValidateDeleteRef(ctx context.Context, gvk schema.GroupVersionKind, ref *core.ResourceRef, dryRun bool) error
	ValidateGvk(ctx context.Context, gvk schema.GroupVersionKind, resource resources.Resource, dryRun bool) (*Reports, error)
	ValidateGlooResource(ctx context.Context, resource resources.Resource, rv GlooResourceValidator) (*Reports, error)
	ValidateGatewayResource(ctx context.Context, resource resources.Resource, rv GatewayResourceValidator, dryRun bool) (*Reports, error)
	ValidateDeleteGlooResource(ctx context.Context, ref *core.ResourceRef, rv DeleteGlooResourceValidator) (*Reports, error)
	ValidateDeleteVirtualService(ctx context.Context, vs *core.ResourceRef, dryRun bool) error
	ValidateDeleteRouteTable(ctx context.Context, rt *core.ResourceRef, dryRun bool) error
	ValidationIsSupported(gvk schema.GroupVersionKind) bool
}

type ValidatorFunc = func(
	context.Context,
	*validation.GlooValidationServiceRequest,
) (*validation.GlooValidationServiceResponse, error)

type validator struct {
	lock              sync.RWMutex
	latestSnapshot    *gloov1snap.ApiSnapshot
	latestSnapshotErr error
	translator        translator.Translator
	//This function replaces a grpc client from when gloo and gateway pods were separate.
	validationFunc               ValidatorFunc
	ignoreProxyValidationFailure bool
	allowWarnings                bool
}

type ValidatorConfig struct {
	translator                   translator.Translator
	validatorFunc                ValidatorFunc
	writeNamespace               string
	ignoreProxyValidationFailure bool
	allowWarnings                bool
}

func NewValidatorConfig(
	translator translator.Translator,
	validatorFunc ValidatorFunc,
	ignoreProxyValidationFailure, allowWarnings bool,
) ValidatorConfig {
	return ValidatorConfig{
		translator:                   translator,
		validatorFunc:                validatorFunc,
		ignoreProxyValidationFailure: ignoreProxyValidationFailure,
		allowWarnings:                allowWarnings,
	}
}

func NewValidator(cfg ValidatorConfig) *validator {
	return &validator{
		translator:                   cfg.translator,
		validationFunc:               cfg.validatorFunc,
		ignoreProxyValidationFailure: cfg.ignoreProxyValidationFailure,
		allowWarnings:                cfg.allowWarnings,
	}
}

func (v *validator) ready() bool {
	return v.latestSnapshot != nil
}

func (v *validator) Sync(ctx context.Context, snap *gloov1snap.ApiSnapshot) error {
	v.lock.Lock() // hashing and cloning resources may mutate the object, so we need to lock
	defer v.lock.Unlock()
	if !v.gatewayUpdate(snap) {
		return nil
	}
	snapCopy := snap.Clone()
	gatewaysByProxy := utils.GatewaysByProxyName(snap.Gateways)
	var errs error
	for proxyName, gatewayList := range gatewaysByProxy {
		_, reports := v.translator.Translate(ctx, proxyName, snap, gatewayList)
		validate := reports.ValidateStrict
		if v.allowWarnings {
			validate = reports.Validate
		}
		if err := validate(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	// When the pod is first starting (aka the first snapshot is received),
	// set the value of mValidConfig with respect to the translation loop above.
	// Without this, mValidConfig will not be exported on /metrics until a new
	// resource is applied (https://github.com/solo-io/gloo/issues/5949).
	if v.latestSnapshot == nil {
		if errs == nil {
			utils2.MeasureOne(ctx, mValidConfig)
		} else {
			utils2.MeasureZero(ctx, mValidConfig)
		}
	}

	v.latestSnapshotErr = errs
	v.latestSnapshot = &snapCopy

	if errs != nil {
		return errors.Wrapf(errs, InvalidSnapshotErrMessage)
	}

	return nil
}

func (v *validator) ValidationIsSupported(gvk schema.GroupVersionKind) bool {
	_, hit := GvkToGatewayResourceValidator[gvk]
	if !hit {
		_, hit := GvkToGlooValidator[gvk]
		return hit
	}
	return true
}

type applyResource func(snap *gloov1snap.ApiSnapshot) (proxyNames []string, resource resources.Resource, ref *core.ResourceRef)

func (v *validator) gatewayUpdate(snap *gloov1snap.ApiSnapshot) bool {

	if v.latestSnapshot == nil {
		return true
	}
	//look at the hash of resources that affect the gateway snapshot
	hashFunc := func(snap *gloov1snap.ApiSnapshot) uint64 {
		toHash := append([]interface{}{}, snap.VirtualHostOptions.AsInterfaces()...)
		toHash = append(toHash, snap.VirtualServices.AsInterfaces()...)
		toHash = append(toHash, snap.Gateways.AsInterfaces()...)
		toHash = append(toHash, snap.RouteOptions.AsInterfaces()...)
		toHash = append(toHash, snap.RouteTables.AsInterfaces()...)
		toHash = append(toHash, snap.HttpGateways.AsInterfaces()...)
		toHash = append(toHash, snap.VirtualHostOptions.AsInterfaces()...)
		hash, err := hashutils.HashAllSafe(nil, toHash...)
		if err != nil {
			panic("this error should never happen, as this is safe hasher")
		}
		return hash
	}
	hashChanged := hashFunc(v.latestSnapshot) != hashFunc(snap)
	return hashChanged
}

// update internal snapshot to handle race where a lot of resources may be deleted at once, before syncer updates
// should be called within a lock
func (v *validator) deleteFromLocalSnapshot(resource resources.Resource) {
	ref := resource.GetMetadata().Ref()
	switch resource.(type) {
	case *v1.VirtualService:
		for i, rt := range v.latestSnapshot.VirtualServices {
			if rt.GetMetadata().Ref().Equal(ref) {
				v.latestSnapshot.VirtualServices = append(v.latestSnapshot.VirtualServices[:i],
					v.latestSnapshot.VirtualServices[i+1:]...)
				break
			}
		}
	case *v1.RouteTable:
		for i, rt := range v.latestSnapshot.RouteTables {
			if rt.GetMetadata().Ref().Equal(ref) {
				v.latestSnapshot.RouteTables = append(v.latestSnapshot.RouteTables[:i], v.latestSnapshot.RouteTables[i+1:]...)
				break
			}
		}
	}
}

func (v *validator) validateSnapshotThreadSafe(ctx context.Context, apply applyResource, dryRun bool) (
	*Reports,
	error,
) {
	// thread-safe implementation of validateSnapshot

	v.lock.Lock()
	defer v.lock.Unlock()

	return v.validateSnapshot(ctx, apply, dryRun)
}

func (v *validator) validateSnapshot(ctx context.Context, apply applyResource, dryRun bool) (*Reports, error) {
	// validate that a snapshot can be modified
	// should be called within a lock
	//
	// validation occurs by the following steps:
	//	1. Clone the most recent snapshot
	//	2. Apply the changes to that snapshot clone
	//	3. Validate the generated proxy of that snapshot clone by making a gRPC call to Gloo.
	//	4. If the proxy is valid, we know that the requested mutation is valid. If this request happens
	//		during a dry run, we don't want to actually apply the change, since this will modify the internal
	//		state of the validator, which is shared across requests. Therefore, only if we are not in a dry run,
	//		we apply the mutation.
	if !v.ready() {
		return nil, NotReadyErr
	}

	ctx = contextutils.WithLogger(ctx, "gateway-validator")

	snapshotClone := v.latestSnapshot.Clone()
	if v.latestSnapshotErr != nil {
		if !dryRun {
			utils2.MeasureZero(ctx, mValidConfig)
		}
		contextutils.LoggerFrom(ctx).Errorw(InvalidSnapshotErrMessage, zap.Error(v.latestSnapshotErr))
		// allow writes if storage is already broken
		return nil, nil
	}

	// verify the mutation against a snapshot clone first, only apply the change to the actual snapshot if this passes
	proxyNames, resource, ref := apply(&snapshotClone)

	gatewaysByProxy := utils.GatewaysByProxyName(snapshotClone.Gateways)

	var (
		errs         error
		proxyReports ProxyReports
		proxies      []*gloov1.Proxy
	)
	for _, proxyName := range proxyNames {
		gatewayList := gatewaysByProxy[proxyName]
		proxy, reports := v.translator.Translate(ctx, proxyName, &snapshotClone, gatewayList)
		validate := reports.ValidateStrict
		if v.allowWarnings {
			validate = reports.Validate
		}
		if err := validate(); err != nil {
			errs = multierr.Append(errs, errors.Wrapf(err, "could not render proxy"))
			continue
		}

		// a nil proxy may have been returned if 0 listeners were created
		if proxy == nil {
			continue
		}

		proxies = append(proxies, proxy)
		// validate the proxy with gloo
		glooValidationResponse, err := v.sendGlooValidationServiceRequest(ctx, &validation.GlooValidationServiceRequest{
			Proxy: proxy,
		})
		if err != nil {
			err = errors.Wrapf(err, "failed to communicate with Gloo validation server")
			if v.ignoreProxyValidationFailure {
				contextutils.LoggerFrom(ctx).Error(err)
			} else {
				errs = multierr.Append(errs, err)
			}
			continue
		}
		if len(glooValidationResponse.GetValidationReports()) != 1 {
			// This was likely caused by a development error
			err := GlooValidationResponseLengthError(glooValidationResponse)
			errs = multierr.Append(errs, err)
			continue
		}

		proxyReport := glooValidationResponse.GetValidationReports()[0].GetProxyReport()
		proxyReports = append(proxyReports, proxyReport)
		if err := validationutils.GetProxyError(proxyReport); err != nil {
			errs = multierr.Append(errs, errors.Wrapf(err, "failed to validate Proxy with Gloo validation server"))
			continue
		}
		if warnings := validationutils.GetProxyWarning(proxyReport); !v.allowWarnings && len(warnings) > 0 {
			for _, warning := range warnings {
				errs = multierr.Append(errs, errors.New(warning))
			}
			continue
		}
	}

	if errs != nil {
		contextutils.LoggerFrom(ctx).Debugf("Rejected %T %v: %v", resource, ref, errs)
		if !dryRun {
			utils2.MeasureZero(ctx, mValidConfig)
		}
		return &Reports{ProxyReports: &proxyReports, Proxies: proxies}, errors.Wrapf(errs,
			"validating %T %v",
			resource,
			ref)
	}

	contextutils.LoggerFrom(ctx).Debugf("Accepted %T %v", resource, ref)
	if !dryRun {
		utils2.MeasureOne(ctx, mValidConfig)
	}

	if !dryRun {
		// update internal snapshot to handle race where a lot of resources may be applied at once, before syncer updates
		apply(v.latestSnapshot)
	}

	return &Reports{ProxyReports: &proxyReports, Proxies: proxies}, nil
}

func (v *validator) ValidateDeleteRef(ctx context.Context, gvk schema.GroupVersionKind, ref *core.ResourceRef, dryRun bool) error {
	if gvk.Group == GatewayGroup {
		if rv, hit := GvkToDeleteGatewayResourceValidator[gvk]; hit {
			return rv.DeleteResource(ctx, ref, v, dryRun)
		}
	} else if gvk.Group == GlooGroup {
		if rv, hit := GvkToDeleteGlooValidator[gvk]; hit {
			_, err := v.ValidateDeleteGlooResource(ctx, ref, rv)
			return err
		}
	}
	contextutils.LoggerFrom(ctx).Debugf("unsupported validation for resource delete ref namespace [%s] name [%s] group [%s] kind [%s]", ref.GetNamespace(), ref.GetName(), gvk.Group, gvk.Kind)
	return nil
}

func (v *validator) ValidateGvk(ctx context.Context, gvk schema.GroupVersionKind, resource resources.Resource, dryRun bool) (*Reports, error) {
	var reports *Reports
	// Gloo has two types of Groups: Gateway, Gloo. This statement is splitting the Validation based off
	// the resource group type.
	if gvk.Group == GatewayGroup {
		if rv, hit := GvkToGatewayResourceValidator[gvk]; hit {
			reports, err := v.ValidateGatewayResource(ctx, resource, rv, dryRun)
			if err != nil {
				return reports, &multierror.Error{Errors: []error{errors.Wrapf(err, "Validating %T failed", resource)}}
			}
			return reports, nil
		}
	} else if gvk.Group == GlooGroup {
		if rv, hit := GvkToGlooValidator[gvk]; hit {
			reports, err := v.ValidateGlooResource(ctx, resource, rv)
			if err != nil {
				return reports, &multierror.Error{Errors: []error{errors.Wrapf(err, "Validating %T failed", resource)}}
			}
			return reports, nil
		}
	}
	return reports, &multierror.Error{Errors: []error{errors.Errorf("failed validating the resoruce [%T] because the group [%s] does not get validated", resource, gvk.Group)}}
}

func (v *validator) ValidateList(ctx context.Context, ul *unstructured.UnstructuredList, dryRun bool) (
	*Reports,
	*multierror.Error,
) {
	var (
		proxies      []*gloov1.Proxy
		proxyReports = ProxyReports{}
		errs         = &multierror.Error{}
	)

	v.lock.Lock()
	originalSnapshot := v.latestSnapshot.Clone()
	v.lock.Unlock()

	for _, item := range ul.Items {

		// this will lock
		var itemProxyReports, err = v.processItem(ctx, item)

		errs = multierror.Append(errs, err)
		if itemProxyReports != nil && itemProxyReports.ProxyReports != nil {
			// ok to return final proxy reports as the latest result includes latest proxy calculated
			// for each resource, as we process incrementally, storing new state in memory as we go
			proxyReports = append(proxyReports, *itemProxyReports.ProxyReports...)
			proxies = append(proxies, itemProxyReports.Proxies...)
		}
	}

	if dryRun {
		// to validate the entire list of changes against one another, each item was applied to the latestSnapshot
		// if this is a dry run, latestSnapshot needs to be reset back to its original value without any of the changes
		v.lock.Lock()
		v.latestSnapshot = &originalSnapshot
		v.lock.Unlock()
	}

	return &Reports{ProxyReports: &proxyReports, Proxies: proxies}, errs
}

func (v *validator) processItem(ctx context.Context, item unstructured.Unstructured) (*Reports, error) {
	// process a single change in a list of changes
	//
	// when calling the specific internal validate method, dryRun and acquireLock are always false:
	// 	dryRun=false: this enables items to be validated against other items in the list
	// 	acquireLock=false: the entire list of changes are called within a single lock
	gv, err := schema.ParseGroupVersion(item.GetAPIVersion())
	if err != nil {
		return &Reports{ProxyReports: &ProxyReports{}}, err
	}

	itemGvk := schema.GroupVersionKind{
		Version: gv.Version,
		Group:   gv.Group,
		Kind:    item.GetKind(),
	}

	jsonBytes, err := item.MarshalJSON()
	if err != nil {
		return &Reports{ProxyReports: &ProxyReports{}}, err
	}

	if newResourceFunc, hit := gloosnapshot.ApiGvkToHashableResource[itemGvk]; hit {
		resource := newResourceFunc()
		if unmarshalErr := skprotoutils.UnmarshalResource(jsonBytes, resource); unmarshalErr != nil {
			return &Reports{ProxyReports: &ProxyReports{}}, WrappedUnmarshalErr(unmarshalErr)
		}
		return v.ValidateGvk(ctx, itemGvk, resource, false)
	}
	// TODO(mitchaman): Handle upstreams
	// should not happen
	return &Reports{ProxyReports: &ProxyReports{}}, errors.Errorf("Unknown group/version/kind, %v", itemGvk)
}

func (v *validator) ValidateDeleteVirtualService(ctx context.Context, vsRef *core.ResourceRef, dryRun bool) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	if !v.ready() {
		return errors.Errorf("Gateway validation is yet not available. Waiting for first snapshot")
	}
	snap := v.latestSnapshot.Clone()

	vs, err := snap.VirtualServices.Find(vsRef.Strings())
	if err != nil {
		// if it's not present in the snapshot, allow deletion
		return nil
	}

	var parentGateways []*core.ResourceRef
	snap.Gateways.Each(func(element *v1.Gateway) {
		virtualServices := virtualServicesForGateway(ctx, snap, element)
		if len(virtualServices) == 0 {
			contextutils.LoggerFrom(ctx).Debugw("Accepted deletion of Virtual Service %v", vsRef)
			return
		}

		for _, ref := range virtualServices {
			if ref.Equal(vsRef) {
				// this gateway points at this virtual service
				parentGateways = append(parentGateways, element.GetMetadata().Ref())

				break
			}
		}
	})

	if len(parentGateways) > 0 {
		err := VirtualServiceDeleteErr(parentGateways)
		if !v.allowWarnings {
			contextutils.LoggerFrom(ctx).Infof("Rejected deletion of Virtual Service %v: %v", vsRef, err)
			return err
		}
		contextutils.LoggerFrom(ctx).Warn("Allowed deletion of Virtual Service %v with warning: %v", vsRef, err)
	} else {
		contextutils.LoggerFrom(ctx).Debugw("Accepted deletion of Virtual Service %v", vsRef)
	}

	if !dryRun {
		v.deleteFromLocalSnapshot(vs)
	}
	return nil
}

func (v *validator) ValidateDeleteRouteTable(ctx context.Context, rtRef *core.ResourceRef, dryRun bool) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	if !v.ready() {
		return errors.Errorf("Gateway validation is yet not available. Waiting for first snapshot")
	}
	snap := v.latestSnapshot.Clone()

	rt, err := snap.RouteTables.Find(rtRef.Strings())
	if err != nil {
		// if it's not present in the snapshot, allow deletion
		return nil
	}

	refsToDelete := sets.NewString(gloo_translator.UpstreamToClusterName(rtRef))

	var parentVirtualServices []*core.ResourceRef
	snap.VirtualServices.Each(func(element *v1.VirtualService) {
		// for each VS, check if its routes contain a ref to deleted rt
		if routesContainRefs(element.GetVirtualHost().GetRoutes(), refsToDelete) {
			parentVirtualServices = append(parentVirtualServices, element.GetMetadata().Ref())
		}
	})

	var parentRouteTables []*core.ResourceRef
	snap.RouteTables.Each(func(element *v1.RouteTable) {
		// for each RT, check if its routes contain a ref to deleted rt
		if routesContainRefs(element.GetRoutes(), refsToDelete) {
			parentRouteTables = append(parentRouteTables, element.GetMetadata().Ref())
		}
	})

	if len(parentVirtualServices) > 0 || len(parentRouteTables) > 0 {
		err := RouteTableDeleteErr(parentVirtualServices, parentRouteTables)
		if !v.allowWarnings {
			contextutils.LoggerFrom(ctx).Debugw("Rejected deletion of Route Table %v: %v", rtRef, err)
			return err
		}
		contextutils.LoggerFrom(ctx).Warn("Allowed deletion of Route Table %v with warning: %v", rtRef, err)
	} else {
		contextutils.LoggerFrom(ctx).Debugw("Accepted Route Table deletion %v", rtRef)
	}

	if !dryRun {
		v.deleteFromLocalSnapshot(rt)
	}
	return nil
}

type GetProxies func(ctx context.Context, resource resources.HashableInputResource, snap *gloov1snap.ApiSnapshot) ([]string, error)

func (v *validator) ValidateDeleteGlooResource(ctx context.Context, ref *core.ResourceRef, rv DeleteGlooResourceValidator) (*Reports, error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	response, err := v.sendGlooValidationServiceRequest(ctx, rv.CreateDeleteRequest(ref))
	logger := contextutils.LoggerFrom(ctx)
	if err != nil {
		if v.ignoreProxyValidationFailure {
			logger.Error(err)
		} else {
			return &Reports{}, err
		}
	}
	// dont log the responsse as proxies and status reports may be too large in large envs.
	logger.Debugf("Got response from GlooValidationService with %d reports", len(response.GetValidationReports()))

	return v.getReportsFromGlooValidationResponse(response)
}

// ValidateGlooResource will validate gloo group presources
func (v *validator) ValidateGlooResource(ctx context.Context, resource resources.Resource, rv GlooResourceValidator) (*Reports, error) {
	return v.validateGlooResource(ctx, resource, rv)
}

func (v *validator) validateGlooResource(ctx context.Context, resource resources.Resource, rv GlooResourceValidator) (*Reports, error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	response, err := v.sendGlooValidationServiceRequest(ctx, rv.CreateModifiedRequest(resource))
	logger := contextutils.LoggerFrom(ctx)
	if err != nil {
		if v.ignoreProxyValidationFailure {
			logger.Error(err)
		} else {
			return &Reports{}, err
		}
	}
	// dont log the responsse as proxies and status reports may be too large in large envs.
	logger.Debugf("Got response from GlooValidationService with %d reports", len(response.GetValidationReports()))

	return v.getReportsFromGlooValidationResponse(response)
}

// ValidateGatewayResource will validate gateway group resources
func (v *validator) ValidateGatewayResource(ctx context.Context, resource resources.Resource, rv GatewayResourceValidator, dryRun bool) (*Reports, error) {
	return v.validateGatewayResource(ctx, resource, rv, dryRun, true)
}

func (v *validator) validateGatewayResource(ctx context.Context, resource resources.Resource, rv GatewayResourceValidator, dryRun, acquireLock bool) (*Reports, error) {
	apply := func(snap *gloov1snap.ApiSnapshot) ([]string, resources.Resource, *core.ResourceRef) {
		resourceRef := resource.GetMetadata().Ref()

		var isUpdate bool
		listOfInputResources, err := snap.GetResourcesList(resource)
		if err != nil {
			contextutils.LoggerFrom(ctx).Error(err)
		}
		for i, existingResource := range listOfInputResources {
			if existingResource.GetMetadata().Ref().Equal(resourceRef) {
				if err := snap.ReplaceResource(i, resource); err != nil {
					contextutils.LoggerFrom(ctx).Error(err)
				}
				isUpdate = true
				break
			}
		}
		if !isUpdate {
			if err := snap.AddToResourceList(resource); err != nil {
				contextutils.LoggerFrom(ctx).Error(err)
			}
		}

		// TODO GetProxies(ctx, resource, snap)
		proxiesToConsider, err := rv.GetProxies(ctx, resource, snap)
		if err != nil {
			contextutils.LoggerFrom(ctx).Error(eris.Wrapf(err, "the resource is %+v", resource.GetMetadata()))
		}
		return proxiesToConsider, resource, resourceRef
	}

	if acquireLock {
		return v.validateSnapshotThreadSafe(ctx, apply, dryRun)
	} else {
		return v.validateSnapshot(ctx, apply, dryRun)
	}
}

// Converts the GlooValidationServiceResponse into Reports.
func (v *validator) getReportsFromGlooValidationResponse(validationResponse *validation.GlooValidationServiceResponse) (
	*Reports,
	error,
) {
	var (
		errs            error
		upstreamReports UpstreamReports
		proxyReports    ProxyReports
		proxies         []*gloov1.Proxy
	)
	for _, report := range validationResponse.GetValidationReports() {
		// Append upstream errors
		for _, usRpt := range report.GetUpstreamReports() {
			upstreamReports = append(upstreamReports, usRpt)
			if err := resourceReportToMultiErr(usRpt); err != nil {
				errs = multierr.Append(errs, errors.Wrapf(err, "failed to validate Upstream with Gloo validation server"))
			}
			if warnings := usRpt.GetWarnings(); !v.allowWarnings && len(warnings) > 0 {
				for _, warning := range warnings {
					errs = multierr.Append(errs, errors.New(warning))
				}
			}
		}

		// Append proxies and proxy reports
		if report.GetProxy() != nil {
			proxies = append(proxies, report.GetProxy())
		}
		if proxyReport := report.GetProxyReport(); proxyReport != nil {
			proxyReports = append(proxyReports, report.GetProxyReport())
			if err := validationutils.GetProxyError(proxyReport); err != nil {
				errs = multierr.Append(errs, errors.Wrapf(err, "failed to validate Proxy with Gloo validation server"))
			}
			if warnings := validationutils.GetProxyWarning(proxyReport); !v.allowWarnings && len(warnings) > 0 {
				for _, warning := range warnings {
					errs = multierr.Append(errs, errors.New(warning))
				}
			}
		}
	}
	return &Reports{
		ProxyReports:    &proxyReports,
		UpstreamReports: &upstreamReports,
		Proxies:         proxies,
	}, errs
}

// sendGlooValidationServiceRequest this will call validator validationFunc
func (v *validator) sendGlooValidationServiceRequest(
	ctx context.Context,
	req *validation.GlooValidationServiceRequest,
) (*validation.GlooValidationServiceResponse, error) {
	logger := contextutils.LoggerFrom(ctx)
	logger.Debugf("Sending request validation request modified:%s, deleted:%s",
		req.GetModifiedResources().String(), req.GetDeletedResources().String())

	// Validate() in  https://github.com/solo-io/gloo/blob/master/projects/gloo/pkg/validation/server.go#L165
	return v.validationFunc(ctx, req)
}

func virtualServicesForGateway(ctx context.Context, snap gloov1snap.ApiSnapshot, gateway *v1.Gateway) []*core.ResourceRef {
	var virtualServices []*core.ResourceRef

	switch gatewayType := gateway.GetGatewayType().(type) {
	case *v1.Gateway_TcpGateway:
		// TcpGateway does not configure VirtualServices
		break
	case *v1.Gateway_HttpGateway:
		virtualServices = gatewayType.HttpGateway.GetVirtualServices()
	case *v1.Gateway_HybridGateway:
		matchedGateways := gatewayType.HybridGateway.GetMatchedGateways()
		if matchedGateways != nil {
			for _, matchedGateway := range matchedGateways {
				if httpGateway := matchedGateway.GetHttpGateway(); httpGateway != nil {
					virtualServices = append(virtualServices, httpGateway.GetVirtualServices()...)
				}
			}
		} else {
			delegatedGateways := gatewayType.HybridGateway.GetDelegatedHttpGateways()
			httpGatewayList := translator.NewHttpGatewaySelector(snap.HttpGateways).SelectMatchableHttpGateways(delegatedGateways, func(err error) {
				logger := contextutils.LoggerFrom(ctx)
				logger.Warnf("failed to select matchable http gateways on gateway: %v", err.Error())
			})

			httpGatewayList.Each(func(element *v1.MatchableHttpGateway) {
				virtualServices = append(virtualServices, element.GetHttpGateway().GetVirtualServices()...)
			})
		}
	}

	return virtualServices
}

// Returns true if any of the given routes delegate to any of the given route tables via a direct reference.
// This is used to determine which route tables are affected when a route table is deleted. Since selectors do not
// represent hard referential constraints, we only need to check direct references here (we can safely remove a route
// table that matches via a selector).
func routesContainRefs(routes []*v1.Route, refs sets.String) bool {
	for _, r := range routes {
		delegate := r.GetDelegateAction()
		if delegate == nil {
			continue
		}

		rtRef := GetDelegateRef(delegate)
		if rtRef == nil {
			continue
		}

		if _, ok := refs[gloo_translator.UpstreamToClusterName(rtRef)]; ok {
			return true
		}
	}
	return false
}

func resourceReportToMultiErr(resourceRpt *validation.ResourceReport) error {
	var multiErr error
	for _, errStr := range resourceRpt.GetErrors() {
		multiErr = multierr.Append(multiErr, errors.New(errStr))
	}
	return multiErr
}
