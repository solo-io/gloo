package failover

import (
	"context"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/solo-io/solo-kit/pkg/utils/prototime"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
)

// Status represents the batch of properties on a FailoverScheme
type Status interface {
	// GetState The current state of the resource.
	GetState() fed_types.FailoverSchemeStatus_State
	// GetMessage A human readable message about the current state of the object.
	GetMessage() string
	// GetObservedGeneration The most recently observed generation of the resource. This value corresponds to the `metadata.generation` of
	// a kubernetes resource.
	GetObservedGeneration() int64
	// GetProcessingTime The time at which this status was recorded.
	GetProcessingTime() *timestamp.Timestamp
}

// StatusBuilder allows a FailoverSchemeStatus to be built incrementally
type StatusBuilder interface {
	Pending(message string) StatusBuilder
	Accept() StatusBuilder
	Fail(err error) StatusBuilder
	Invalidate(err error) StatusBuilder
	Set(Status) StatusBuilder // useful for testing
	Build() *fedv1.FailoverScheme
}

type StatusManager struct {
	client    fedv1.FailoverSchemeClient
	namespace string
}

func NewStatusManager(client fedv1.FailoverSchemeClient, namespace string) *StatusManager {
	return &StatusManager{
		client:    client,
		namespace: namespace,
	}
}

func (m StatusManager) GetStatus(obj *fedv1.FailoverScheme) Status {
	statuses := obj.Status.GetNamespacedStatuses()
	if statuses == nil {
		return nil
	}
	resourceStatus := statuses[m.namespace]
	return &StatusImpl{
		State:              resourceStatus.GetState(),
		Message:            resourceStatus.GetMessage(),
		ObservedGeneration: resourceStatus.GetObservedGeneration(),
		ProcessingTime:     resourceStatus.GetProcessingTime(),
	}
}

func (m StatusManager) NewStatusBuilder(obj *fedv1.FailoverScheme) StatusBuilder {
	return &statusBuilderImpl{
		statusNamespace: m.namespace,
		obj:             obj,
	}
}

func (m StatusManager) UpdateStatus(ctx context.Context, builder StatusBuilder) error {
	return m.client.UpdateFailoverSchemeStatus(ctx, builder.Build())
}

type statusBuilderImpl struct {
	obj             *fedv1.FailoverScheme
	statusNamespace string
	status          Status
}

func (f *statusBuilderImpl) Pending(message string) StatusBuilder {
	f.status = &StatusImpl{
		State:              fed_types.FailoverSchemeStatus_PENDING,
		Message:            message,
		ObservedGeneration: f.obj.GetGeneration(),
		ProcessingTime:     prototime.Now(),
	}
	return f
}

func (f *statusBuilderImpl) Accept() StatusBuilder {
	f.status = &StatusImpl{
		State:              fed_types.FailoverSchemeStatus_ACCEPTED,
		ObservedGeneration: f.obj.GetGeneration(),
		ProcessingTime:     prototime.Now(),
	}
	return f
}

func (f *statusBuilderImpl) Fail(err error) StatusBuilder {
	f.status = &StatusImpl{
		State:              fed_types.FailoverSchemeStatus_FAILED,
		Message:            err.Error(),
		ObservedGeneration: f.obj.GetGeneration(),
		ProcessingTime:     prototime.Now(),
	}
	return f
}

func (f *statusBuilderImpl) Invalidate(err error) StatusBuilder {
	f.status = &StatusImpl{
		State:              fed_types.FailoverSchemeStatus_INVALID,
		Message:            err.Error(),
		ObservedGeneration: f.obj.GetGeneration(),
		ProcessingTime:     prototime.Now(),
	}
	return f
}

func (f *statusBuilderImpl) Set(status Status) StatusBuilder {
	f.status = status
	return f
}

func (f *statusBuilderImpl) Build() *fedv1.FailoverScheme {
	statuses := f.obj.Status.GetNamespacedStatuses()
	if statuses == nil {
		f.obj.Status = fed_types.FailoverSchemeStatus{
			NamespacedStatuses: map[string]*fed_types.FailoverSchemeStatus_Status{
				f.statusNamespace: {
					State:              f.status.GetState(),
					Message:            f.status.GetMessage(),
					ObservedGeneration: f.status.GetObservedGeneration(),
					ProcessingTime:     f.status.GetProcessingTime(),
				},
			},
		}
		return f.obj
	}

	statuses[f.statusNamespace] = &fed_types.FailoverSchemeStatus_Status{
		State:              f.status.GetState(),
		Message:            f.status.GetMessage(),
		ObservedGeneration: f.status.GetObservedGeneration(),
		ProcessingTime:     f.status.GetProcessingTime(),
	}
	return f.obj
}

var _ Status = new(StatusImpl)

// StatusImpl represents the batch of properties on a FailoverScheme
type StatusImpl struct {
	// The current state of the resource.
	State fed_types.FailoverSchemeStatus_State
	// A human readable message about the current state of the object.
	Message string
	// The most recently observed generation of the resource. This value corresponds to the `metadata.generation` of
	// a kubernetes resource.
	ObservedGeneration int64
	// The time at which this status was recorded.
	ProcessingTime *timestamp.Timestamp
}

func (s StatusImpl) GetState() fed_types.FailoverSchemeStatus_State {
	return s.State
}

func (s StatusImpl) GetMessage() string {
	return s.Message
}

func (s StatusImpl) GetObservedGeneration() int64 {
	return s.ObservedGeneration
}

func (s StatusImpl) GetProcessingTime() *timestamp.Timestamp {
	return s.ProcessingTime
}
