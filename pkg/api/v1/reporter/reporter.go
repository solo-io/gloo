package reporter

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

type ResourceErrors map[resources.InputResource]error

func (e ResourceErrors) Accept(res ...resources.InputResource) ResourceErrors {
	for _, r := range res {
		e[r] = nil
	}
	return e
}

func (e ResourceErrors) Merge(resErrs ResourceErrors) {
	for k, v := range resErrs {
		e[k] = v
	}
}

func (e ResourceErrors) AddError(res resources.InputResource, err error) {
	if err == nil {
		return
	}
	e[res] = multierror.Append(e[res], err)
}

func (e ResourceErrors) Validate() error {
	var errs error
	for res, err := range e {
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "invalid resource %v.%v", res.GetMetadata().Namespace, res.GetMetadata().Name))
		}
	}
	return errs
}

type Reporter interface {
	WriteReports(ctx context.Context, errs ResourceErrors, subresourceStatuses map[string]*core.Status) error
}

type reporter struct {
	clients clients.ResourceClients
	ref     string
}

func NewReporter(reporterRef string, resourceClients ...clients.ResourceClient) Reporter {
	clientsByKind := make(clients.ResourceClients)
	for _, client := range resourceClients {
		clientsByKind[client.Kind()] = client
	}
	return &reporter{
		ref:     reporterRef,
		clients: clientsByKind,
	}
}

func (r *reporter) WriteReports(ctx context.Context, resourceErrs ResourceErrors, subresourceStatuses map[string]*core.Status) error {
	ctx = contextutils.WithLogger(ctx, "reporter")
	logger := contextutils.LoggerFrom(ctx)

	var merr *multierror.Error

	for resource, validationError := range resourceErrs {
		kind := resources.Kind(resource)
		client, ok := r.clients[kind]
		if !ok {
			return errors.Errorf("reporter: was passed resource of kind %v but no client to support it", kind)
		}
		status := statusFromError(r.ref, validationError, subresourceStatuses)
		resourceToWrite := resources.Clone(resource).(resources.InputResource)
		if status.Equal(resource.GetStatus()) {
			logger.Debugf("skipping report for %v as it has not changed", resourceToWrite.GetMetadata().Ref())
			continue
		}
		resourceToWrite.SetStatus(status)
		res, err := client.Write(resourceToWrite, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			err := errors.Wrapf(err, "failed to write status %v for resource %v", status, resource.GetMetadata().Name)
			logger.Warn(err)
			merr = multierror.Append(merr, err)
			continue
		}
		resources.UpdateMetadata(resource, func(meta *core.Metadata) {
			meta.ResourceVersion = res.GetMetadata().ResourceVersion
		})

		logger.Infof("wrote report %v : %v", resourceToWrite.GetMetadata().Ref(), status)
	}
	return merr.ErrorOrNil()
}

func statusFromError(ref string, err error, subresourceStatuses map[string]*core.Status) core.Status {
	if err != nil {
		return core.Status{
			State:               core.Status_Rejected,
			Reason:              err.Error(),
			ReportedBy:          ref,
			SubresourceStatuses: subresourceStatuses,
		}
	}
	return core.Status{
		State:               core.Status_Accepted,
		ReportedBy:          ref,
		SubresourceStatuses: subresourceStatuses,
	}
}
