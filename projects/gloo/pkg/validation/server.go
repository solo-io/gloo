package validation

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/solo-io/gloo/pkg/utils/syncutil"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	sk_resources "github.com/solo-io/solo-kit/pkg/api/v1/resources"
	sk_kubernetes "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var SyncNotCalledError = eris.New("proxy validation called before the validation server received its first sync of resources")

type Validator interface {
	v1snap.ApiSyncer
	validation.GlooValidationServiceServer
	ValidateGloo(ctx context.Context, proxy *v1.Proxy, resource resources.Resource, shouldDelete bool) ([]*GlooValidationReport, error)
}

// ValidatorConfig is used to configure the validator
type ValidatorConfig struct {
	Ctx                 context.Context
	GlooValidatorConfig GlooValidatorConfig
}

type validator struct {
	// note to devs: this can be called in parallel by the validation webhook and main translation loops at the same time
	// any stateful fields should be protected by a mutex or themselves be synchronized (like the xds sanitizer / translator)
	lock           sync.RWMutex
	latestSnapshot *v1snap.ApiSnapshot
	validator      GlooValidator
	translator     translator.Translator
	notifyResync   map[*validation.NotifyOnResyncRequest]chan struct{}
	ctx            context.Context
}

func NewValidator(config ValidatorConfig) *validator {
	return &validator{
		translator:   config.GlooValidatorConfig.Translator,
		validator:    NewGlooValidator(config.GlooValidatorConfig),
		notifyResync: make(map[*validation.NotifyOnResyncRequest]chan struct{}, 1),
		ctx:          config.Ctx,
	}
}

// only call within a lock
// should we notify on this snap update
func (s *validator) shouldNotify(snap *v1snap.ApiSnapshot) bool {
	if s.latestSnapshot == nil {
		return true
	}
	// rather than compare the hash of the whole snapshot,
	// we compare the hash of resources that can affect
	// the validation result (which excludes Endpoints)
	hashFunc := func(snap *v1snap.ApiSnapshot) (uint64, error) {
		toHash := append([]interface{}{}, snap.Upstreams.AsInterfaces()...)
		toHash = append(toHash, snap.UpstreamGroups.AsInterfaces()...)
		toHash = append(toHash, snap.Secrets.AsInterfaces()...)
		toHash = append(toHash, snap.AuthConfigs.AsInterfaces()...)
		toHash = append(toHash, snap.Ratelimitconfigs.AsInterfaces()...)
		// we also include proxies as this will help
		// the gateway to resync in case the proxy was deleted
		toHash = append(toHash, snap.Proxies.AsInterfaces()...)

		hash, err := hashutils.HashAllSafe(nil, toHash...)
		if err != nil {
			contextutils.LoggerFrom(context.Background()).DPanic("this error should never happen, as this is safe hasher")
			return 0, errors.New("this error should never happen, as this is safe hasher")
		}
		return hash, nil
	}
	oldHash, oldHashErr := hashFunc(s.latestSnapshot)
	newHash, newHashErr := hashFunc(snap)
	// If we cannot hash then we choose to treat them as different hashes since this is just a performance optimization.
	// In worst case we'd prefer correctness
	hashChanged := oldHash != newHash || oldHashErr != nil || newHashErr != nil

	logger := contextutils.LoggerFrom(s.ctx)
	// stringifying the snapshot may be an expensive operation, so we'd like to avoid building the large
	// string if we're not even going to log it anyway
	if contextutils.GetLogLevel() == zapcore.DebugLevel {
		logger.Infow("last validation snapshot", zap.Any("latestSnapshot", syncutil.StringifySnapshot(s.latestSnapshot)))
		logger.Infow("current validation snapshot", zap.Any("currentSnapshot", syncutil.StringifySnapshot(snap)))
	}
	logger.Infow("validation hash comparison",
		zap.Bool("hashChanged", hashChanged),
		zap.String("issue", "8539"),
	)

	// notify if the hash of what we care about has changed
	return hashChanged
}

// only call within a lock
// notify all receivers
func (s *validator) pushNotifications() {
	logger := contextutils.LoggerFrom(s.ctx)
	receiverCount := len(s.notifyResync)
	logger.Infow("pushing notifications",
		zap.Int("receiverCount", receiverCount),
		zap.String("issue", "8539"),
	)
	for req, receiver := range s.notifyResync {
		logger.Infow("pushing notification for receiver",
			zap.Any("receiver", receiver),
			zap.Any("request", req),
			zap.String("issue", "8539"),
		)
		receiver := receiver
		go func() {
			select {
			// only write to channel if it's empty
			case receiver <- struct{}{}:
			default:
			}
		}()
	}
}

// the gloo snapshot has changed.
// update the local snapshot, notify subscribers
func (s *validator) Sync(ctx context.Context, snap *v1snap.ApiSnapshot) error {
	logger := contextutils.LoggerFrom(ctx)
	logger.Infow("Gloo Validator syncing snapshot",
		zap.String("issue", "8539"),
		zap.Int("upstreamCount", len(snap.Upstreams)),
		zap.Int("secretCount", len(snap.Secrets)),
		zap.Int("proxyCount", len(snap.Proxies)),
		zap.Int("upstreamGroupCount", len(snap.UpstreamGroups)),
		zap.Int("authConfigCount", len(snap.AuthConfigs)),
		zap.Int("rateLimitConfigCount", len(snap.Ratelimitconfigs)),
	)
	snapCopy := snap.Clone()
	s.lock.Lock()
	shouldNotify := s.shouldNotify(snap)
	if shouldNotify {
		logger.Infow("snapshot changed, pushing notifications", zap.String("issue", "8539"))
		s.pushNotifications()
	} else {
		logger.Infow("snapshot unchanged, skipping notifications", zap.String("issue", "8539"))
	}
	s.latestSnapshot = &snapCopy
	s.lock.Unlock()
	logger.Infow("Gloo Validator synced snapshot",
		zap.String("issue", "8539"),
		zap.Bool("notificationsTriggered", shouldNotify),
	)
	return nil
}

func (s *validator) NotifyOnResync(req *validation.NotifyOnResyncRequest, stream validation.GlooValidationService_NotifyOnResyncServer) error {
	// send initial response as ACK
	if err := stream.Send(&validation.NotifyOnResyncResponse{}); err != nil {
		return err
	}

	// initialize a receiver. this will receive all update notifications
	// size of one so we don't queue multiple notifications
	receiver := make(chan struct{}, 1)

	logger := contextutils.LoggerFrom(s.ctx)

	// add the receiver to our map
	s.lock.Lock()
	s.notifyResync[req] = receiver
	logger.Debug("added receiver to map", zap.Any("newReceiver", receiver), zap.Any("afterMap", s.notifyResync), zap.Any("validator", s))
	s.lock.Unlock()

	defer func() {
		// remove the receiver from the map
		s.lock.Lock()
		delete(s.notifyResync, req)
		logger.Debug("removed receiver from map", zap.Any("removedReceiver", req), zap.Any("afterMap", s.notifyResync), zap.Any("validator", s))
		s.lock.Unlock()
	}()

	// loop forever, sending a notification
	// whenever we read from the receiver channel
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-receiver:
			if err := stream.Send(&validation.NotifyOnResyncResponse{}); err != nil {
				contextutils.LoggerFrom(stream.Context()).Errorw("failed to send validation resync notification", zap.Error(err))
			}
		}
	}
}

// Validate is a gRPC call that we use for validating resources against a request to add upstreams and secrets.
func (s *validator) Validate(ctx context.Context, req *validation.GlooValidationServiceRequest) (*validation.GlooValidationServiceResponse, error) {
	logger := contextutils.LoggerFrom(ctx)
	logger.Infow("received proxy validation request",
		zap.String("proxyName", req.GetProxy().GetMetadata().GetName()),
		zap.String("proxyNamespace", req.GetProxy().GetMetadata().GetNamespace()),
		zap.Bool("hasModifiedResources", req.GetModifiedResources() != nil),
		zap.Bool("hasDeletedResources", req.GetDeletedResources() != nil),
	)

	s.lock.Lock()
	// we may receive a Validate call before a Sync has occurred
	if s.latestSnapshot == nil {
		logger.Warnw("validation called before sync", zap.String("issue", "8539"))
		s.lock.Unlock()
		return nil, SyncNotCalledError
	}
	snapCopy := s.latestSnapshot.Clone() // cloning can mutate so we need a write lock
	s.lock.Unlock()

	// update the snapshot copy with the resources from the request
	applyRequestToSnapshot(&snapCopy, req)

	reports := s.validator.Validate(ctx, req.GetProxy(), &snapCopy, false)

	var validationReports []*validation.ValidationReport
	// convert the reports for the gRPC response
	for _, rep := range reports {
		validationReports = append(validationReports, convertToValidationReport(rep.ProxyReport, rep.ResourceReports, rep.Proxy))
	}

	logger.Infow("completed proxy validation request",
		zap.Int("reportCount", len(validationReports)),
	)

	return &validation.GlooValidationServiceResponse{
		ValidationReports: validationReports,
	}, nil
}

func HandleResourceDeletion(snapshot *v1snap.ApiSnapshot, resource resources.Resource) error {
	if _, ok := resource.(*sk_kubernetes.KubeNamespace); ok {
		// Special case to handle namespace deletion
		snapshot.RemoveMatches(func(metadata *core.Metadata) bool {
			return resource.GetMetadata().GetNamespace() == metadata.GetNamespace()
		})
		return nil
	} else {
		return snapshot.RemoveFromResourceList(resource)
	}
}

// ValidateGloo replaces the functionality of Validate.  Validate is still a method that needs to be
// exported because it is used as a gRPC service. A synced version of the snapshot is needed for
// gloo validation.
func (s *validator) ValidateGloo(ctx context.Context, proxy *v1.Proxy, resource resources.Resource, shouldDelete bool) ([]*GlooValidationReport, error) {
	logger := contextutils.LoggerFrom(ctx)

	var resourceName, resourceNamespace, resourceKind string
	if resource != nil {
		resourceName = resource.GetMetadata().GetName()
		resourceNamespace = resource.GetMetadata().GetNamespace()
		resourceKind = sk_resources.Kind(resource)
	}

	logger.Infow("received gloo validation request",
		zap.String("proxyName", proxy.GetMetadata().GetName()),
		zap.String("proxyNamespace", proxy.GetMetadata().GetNamespace()),
		zap.String("resourceName", resourceName),
		zap.String("resourceNamespace", resourceNamespace),
		zap.String("resourceKind", resourceKind),
		zap.Bool("shouldDelete", shouldDelete),
	)

	// the gateway validator will call this function to validate Gloo resources.
	s.lock.Lock()
	// we may receive a Validate call before a Sync has occurred
	if s.latestSnapshot == nil {
		logger.Warnw("gloo validation called before sync", zap.String("issue", "8539"))
		s.lock.Unlock()
		return nil, SyncNotCalledError
	}
	snapCopy := s.latestSnapshot.Clone() // cloning can mutate so we need a write lock
	s.lock.Unlock()

	if resource != nil {
		if shouldDelete {
			logger.Infow("handling resource deletion",
				zap.String("resourceName", resourceName),
				zap.String("resourceKind", resourceKind),
			)
			if err := HandleResourceDeletion(&snapCopy, resource); err != nil {
				logger.Errorw("failed to handle resource deletion",
					zap.Error(err),
					zap.String("resourceName", resourceName),
					zap.String("resourceKind", resourceKind),
				)
				return nil, err
			}

			// If we are deleting an Upstream with a Kube destination, we also want to remove the associated "fake" Upstream from the snapshot
			switch typedResource := resource.(type) {
			case *v1.Upstream:
				if typedResource.GetKube() != nil {
					fakeUpstreamName := fmt.Sprintf("%s%s", kubernetes.FakeUpstreamNamePrefix, resource.GetMetadata().GetName())
					logger.Infow("removing associated fake upstream",
						zap.String("fakeUpstreamName", fakeUpstreamName),
						zap.String("originalUpstreamName", resourceName),
					)
					kubeSvcUs := &v1.Upstream{
						Metadata: &core.Metadata{
							Namespace: resource.GetMetadata().GetNamespace(),
							Name:      fakeUpstreamName,
						},
					}
					if err := snapCopy.RemoveFromResourceList(kubeSvcUs); err != nil {
						logger.Errorw("failed to remove fake upstream",
							zap.Error(err),
							zap.String("fakeUpstreamName", fakeUpstreamName),
						)
						return nil, err
					}
				}
			}
		} else {
			logger.Infow("upserting resource to snapshot",
				zap.String("resourceName", resourceName),
				zap.String("resourceKind", resourceKind),
			)
			if err := snapCopy.UpsertToResourceList(resource); err != nil {
				logger.Errorw("failed to upsert resource",
					zap.Error(err),
					zap.String("resourceName", resourceName),
					zap.String("resourceKind", resourceKind),
				)
				return nil, err
			}
		}
	}

	reports := s.validator.Validate(ctx, proxy, &snapCopy, shouldDelete)

	logger.Infow("completed gloo validation",
		zap.Int("reportCount", len(reports)),
		zap.String("proxyName", proxy.GetMetadata().GetName()),
	)

	return reports, nil
}

// updates the given snapshot with the resources from the request
func applyRequestToSnapshot(snap *v1snap.ApiSnapshot, req *validation.GlooValidationServiceRequest) {
	// if we want to change the type, we could use API snapshots as containers.  Like this struct
	// projects/gloo/pkg/validation/api_snapshot_request.go
	if req.GetModifiedResources() != nil {
		existingUpstreams := snap.Upstreams.AsResources()
		modifiedUpstreams := utils.UpstreamsToResourceList(req.GetModifiedResources().GetUpstreams())
		mergedUpstreams := utils.MergeResourceLists(existingUpstreams, modifiedUpstreams)
		snap.Upstreams = utils.ResourceListToUpstreamList(mergedUpstreams)
	} else if req.GetDeletedResources() != nil {
		// Upstreams
		existingUpstreams := snap.Upstreams.AsResources()
		deletedUpstreamRefs := req.GetDeletedResources().GetUpstreamRefs()
		// If we are deleting an Upstream with a Kube destination, we also want to remove the associated "fake" Upstream from the snapshot
		// Since we only have refs here, attempt to delete the "fake" Upstream corresponding with all refs
		// If none exists this will be a no-op
		for _, ref := range req.GetDeletedResources().GetUpstreamRefs() {
			deletedUpstreamRefs = append(deletedUpstreamRefs, &core.ResourceRef{
				Namespace: ref.GetNamespace(),
				Name:      fmt.Sprintf("%s%s", kubernetes.FakeUpstreamNamePrefix, ref.GetName()),
			})
		}
		finalUpstreams := utils.DeleteResources(existingUpstreams, deletedUpstreamRefs)
		snap.Upstreams = utils.ResourceListToUpstreamList(finalUpstreams)
		// Secrets
		existingSecrets := snap.Secrets.AsResources()
		deletedSecretRefs := req.GetDeletedResources().GetSecretRefs()
		finalSecrets := utils.DeleteResources(existingSecrets, deletedSecretRefs)
		snap.Secrets = utils.ResourceListToSecretList(finalSecrets)
	}
}

func convertToValidationReport(proxyReport *validation.ProxyReport, resourceReports reporter.ResourceReports, proxy *v1.Proxy) *validation.ValidationReport {
	var upstreamReports []*validation.ResourceReport

	for resource, report := range resourceReports {
		switch sk_resources.Kind(resource) {
		case "*v1.Upstream":
			upstreamReports = append(upstreamReports, &validation.ResourceReport{
				ResourceRef: resource.GetMetadata().Ref(),
				Warnings:    report.Warnings,
				Errors:      getErrors(report.Errors),
			})
		}
		// TODO add other resources types here
	}

	return &validation.ValidationReport{
		ProxyReport:     proxyReport,
		UpstreamReports: upstreamReports,
		Proxy:           proxy,
	}
}

func getErrors(err error) []string {
	if err == nil {
		return []string{}
	}
	switch err.(type) {
	case *multierror.Error:
		var errorStrings []string
		for _, e := range err.(*multierror.Error).Errors {
			errorStrings = append(errorStrings, e.Error())
		}
		return errorStrings
	}
	return []string{err.Error()}
}

// Update the validation report so that route errors that were changed into warnings during sanitization
// are also switched in the report results
func routeErrorToWarnings(resourceReport reporter.ResourceReports, validationReport *validation.ProxyReport) {
	// Only proxy reports are needed
	resourceReport = resourceReport.FilterByKind("*v1.Proxy")
	resourceReportErrors := make(map[string]struct{})
	resourceReportWarnings := make(map[string]struct{})
	for _, report := range resourceReport {
		if report.Errors != nil {
			for _, rError := range report.Errors.(*multierror.Error).Errors {
				resourceReportErrors[rError.Error()] = struct{}{}
			}
		}

		for _, rWarning := range report.Warnings {
			resourceReportWarnings[rWarning] = struct{}{}
		}
	}

	for _, listenerReport := range validationReport.GetListenerReports() {
		virtualHostReports := utils.GetVhostReportsFromListenerReport(listenerReport)

		for _, virtualHostReport := range virtualHostReports {
			for _, routeReport := range virtualHostReport.GetRouteReports() {
				modifiedErrors := make([]*validation.RouteReport_Error, 0)
				for _, rError := range routeReport.GetErrors() {
					if _, inErrors := resourceReportErrors[rError.String()]; !inErrors {
						if _, inWarnings := resourceReportWarnings[rError.String()]; inWarnings {
							warning := &validation.RouteReport_Warning{
								Type:   validation.RouteReport_Warning_Type(rError.GetType()),
								Reason: rError.GetReason(),
							}
							routeReport.Warnings = append(routeReport.GetWarnings(), warning)
						}
					} else {
						modifiedErrors = append(modifiedErrors, rError)
					}
				}
				routeReport.Errors = modifiedErrors
			}
		}
	}
}

type ValidationServer interface {
	validation.GlooValidationServiceServer
	SetValidator(v Validator)
	Register(grpcServer *grpc.Server)
}

type validationServer struct {
	lock      sync.RWMutex
	validator Validator
}

func NewValidationServer() *validationServer {
	return &validationServer{}
}

func (s *validationServer) SetValidator(v Validator) {
	s.lock.Lock()
	s.validator = v
	s.lock.Unlock()
}

func (s *validationServer) Register(grpcServer *grpc.Server) {
	validation.RegisterGlooValidationServiceServer(grpcServer, s)
}

func (s *validationServer) NotifyOnResync(req *validation.NotifyOnResyncRequest, stream validation.GlooValidationService_NotifyOnResyncServer) error {
	s.lock.RLock()
	validator := s.validator
	s.lock.RUnlock()

	return validator.NotifyOnResync(req, stream)
}

func (s *validationServer) Validate(ctx context.Context, req *validation.GlooValidationServiceRequest) (*validation.GlooValidationServiceResponse, error) {
	s.lock.RLock()
	validator := s.validator
	s.lock.RUnlock()

	return validator.Validate(ctx, req)
}
