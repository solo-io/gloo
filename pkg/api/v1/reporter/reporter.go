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

func (e ResourceErrors) AddError(res resources.InputResource, err error) {
	if err == nil {
		return
	}
	e[res] = multierror.Append(e[res], err)
}

type Reporter interface {
	WriteReports(ctx context.Context, errs ResourceErrors) error
}

type reporter struct {
	clients map[string]clients.ResourceClient
}

func NewReporter(resourceClients ...clients.ResourceClient) Reporter {
	clientsByKind := make(map[string]clients.ResourceClient)
	for _, client := range resourceClients {
		clientsByKind[client.Kind()] = client
	}
	return &reporter{
		clients: clientsByKind,
	}
}

func (r *reporter) WriteReports(ctx context.Context, resourceErrs ResourceErrors) error {
	ctx = contextutils.WithLogger(ctx, "reporter")
	for resource, validationError := range resourceErrs {
		kind := resources.Kind(resource)
		client, ok := r.clients[kind]
		if !ok {
			return errors.Errorf("reporter: was passed resource of kind %v but no client to support it", kind)
		}
		status := statusFromError(validationError)
		resourceToWrite := resources.Clone(resource).(resources.InputResource)
		resourceToWrite.SetStatus(status)
		if _, err := client.Write(resourceToWrite, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: true,
		}); err != nil {
			return errors.Wrapf(err, "failed to write status %v for resource %v", status, resource.GetMetadata().Name)
		}
	}
	return nil
}

func statusFromError(err error) core.Status {
	if err != nil {
		return core.Status{
			State:  core.Status_Rejected,
			Reason: err.Error(),
		}
	}
	return core.Status{
		State: core.Status_Accepted,
	}
}
