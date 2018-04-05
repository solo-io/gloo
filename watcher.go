package storage

import "github.com/solo-io/gloo-api/pkg/api/types/v1"

type Watcher struct {
	runFunc func(stop <-chan struct{}, errs chan error)
}

func NewWatcher(runFunc func(stop <-chan struct{}, errs chan error)) *Watcher {
	return &Watcher{runFunc: runFunc}
}

func (w *Watcher) Run(stop <-chan struct{}, errs chan error) {
	w.runFunc(stop, errs)
}

type UpstreamEventHandler interface {
	OnAdd(updatedList []*v1.Upstream, obj *v1.Upstream)
	OnUpdate(updatedList []*v1.Upstream, newObj *v1.Upstream)
	OnDelete(updatedList []*v1.Upstream, obj *v1.Upstream)
}

type VirtualHostEventHandler interface {
	OnAdd(updatedList []*v1.VirtualHost, obj *v1.VirtualHost)
	OnUpdate(updatedList []*v1.VirtualHost, newObj *v1.VirtualHost)
	OnDelete(updatedList []*v1.VirtualHost, obj *v1.VirtualHost)
}

// UpstreamEventHandlerFuncs is an adaptor to let you easily specify as many or
// as few of the notification functions as you want while still implementing
// UpstreamEventHandler.
type UpstreamEventHandlerFuncs struct {
	AddFunc    func(updatedList []*v1.Upstream, obj *v1.Upstream)
	UpdateFunc func(updatedList []*v1.Upstream, newObj *v1.Upstream)
	DeleteFunc func(updatedList []*v1.Upstream, obj *v1.Upstream)
}

// OnAdd calls AddFunc if it's not nil.
func (r UpstreamEventHandlerFuncs) OnAdd(updatedList []*v1.Upstream, obj *v1.Upstream) {
	if r.AddFunc != nil {
		r.AddFunc(updatedList, obj)
	}
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r UpstreamEventHandlerFuncs) OnUpdate(updatedList []*v1.Upstream, newObj *v1.Upstream) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(updatedList, newObj)
	}
}

// OnDelete calls DeleteFunc if it's not nil.
func (r UpstreamEventHandlerFuncs) OnDelete(updatedList []*v1.Upstream, obj *v1.Upstream) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(updatedList, obj)
	}
}

// VirtualHostEventHandlerFuncs is an adaptor to let you easily specify as many or
// as few of the notification functions as you want while still implementing
// VirtualHostEventHandler.
type VirtualHostEventHandlerFuncs struct {
	AddFunc    func(updatedList []*v1.VirtualHost, obj *v1.VirtualHost)
	UpdateFunc func(updatedList []*v1.VirtualHost, newObj *v1.VirtualHost)
	DeleteFunc func(updatedList []*v1.VirtualHost, obj *v1.VirtualHost)
}

// OnAdd calls AddFunc if it's not nil.
func (r VirtualHostEventHandlerFuncs) OnAdd(updatedList []*v1.VirtualHost, obj *v1.VirtualHost) {
	if r.AddFunc != nil {
		r.AddFunc(updatedList, obj)
	}
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r VirtualHostEventHandlerFuncs) OnUpdate(updatedList []*v1.VirtualHost, newObj *v1.VirtualHost) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(updatedList, newObj)
	}
}

// OnDelete calls DeleteFunc if it's not nil.
func (r VirtualHostEventHandlerFuncs) OnDelete(updatedList []*v1.VirtualHost, obj *v1.VirtualHost) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(updatedList, obj)
	}
}
