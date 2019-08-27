package status

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

//go:generate mockgen -destination mocks/input_resource_status_getter_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/status InputResourceStatusGetter

type InputResourceStatusGetter interface {
	GetApiStatusFromResource(inputResource resources.InputResource) *v1.Status
}

type inputResourceStatusGetter struct{}

var _ InputResourceStatusGetter = inputResourceStatusGetter{}

func NewInputResourceStatusGetter() InputResourceStatusGetter {
	return &inputResourceStatusGetter{}
}

func (inputResourceStatusGetter) GetApiStatusFromResource(inputResource resources.InputResource) *v1.Status {
	var (
		metadata = inputResource.GetMetadata()
		status   = inputResource.GetStatus()
	)

	switch status.State {
	case core.Status_Pending:
		return &v1.Status{Code: v1.Status_WARNING, Message: ResourcePending(metadata.Namespace, metadata.Name)}
	case core.Status_Accepted:
		return &v1.Status{Code: v1.Status_OK}
	case core.Status_Rejected:
		return &v1.Status{Code: v1.Status_ERROR, Message: ResourceRejected(metadata.Namespace, metadata.Name, status.Reason)}
	default: // unhandled status value
		return &v1.Status{Code: v1.Status_ERROR, Message: UnknownFailure(metadata.Namespace, metadata.Name, status.State)}
	}
}
