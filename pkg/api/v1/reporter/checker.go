package reporter

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// check reports on objects
type StatusHandler interface {
	OnAccepted(resource resources.Resource)
	OnRejected(resource resources.Resource)
	OnPending(resource resources.Resource)
}

func Check(handler StatusHandler, resources ...resources.InputResource) {
	for _, res := range resources {
		switch res.GetStatus().State {
		case core.Status_Pending:
			handler.OnPending(res)
		}
	}
}

type StatusHandlerFuncs struct {
	OnAcceptedF func(resource resources.Resource)
	OnRejectedF func(resource resources.Resource)
	OnPendingF  func(resource resources.Resource)
}

func (f *StatusHandlerFuncs) OnAccepted(resource resources.Resource) {
	if f.OnAcceptedF != nil {
		f.OnAcceptedF(resource)
	}
}

func (f *StatusHandlerFuncs) OnRejected(resource resources.Resource) {
	if f.OnRejectedF != nil {
		f.OnRejectedF(resource)
	}
}

func (f *StatusHandlerFuncs) OnPending(resource resources.Resource) {
	if f.OnPendingF != nil {
		f.OnPendingF(resource)
	}
}
