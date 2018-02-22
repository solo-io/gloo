package server

import (
	"log"
	"time"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
)

// TODO(ashish) - rename to Subscriber
type Handler interface {
	Update(*v1.Upstream)
	Remove(string)
}

type controller struct {
	upstreams storage.Upstreams
	watcher   *storage.Watcher
	handlers  []Handler
}

func newController(resyncDelay time.Duration, upstreams storage.Upstreams) *controller {
	c := &controller{upstreams: upstreams}
	watcher, err := upstreams.Watch(&upstreamHandler{c: c})
	if err != nil {
		// do something
	}
	c.watcher = watcher
	return c
}

func (c *controller) AddHandler(h Handler) {
	c.handlers = append(c.handlers, h)
}

func (c *controller) Run(stop <-chan struct{}) {
	errs := make(chan error)
	go func() {
		for {
			select {
			case <-errs:
			case <-stop:
				return
			}
		}
	}()
	go func() { c.watcher.Run(stop, errs) }()
	log.Println("Controller started")
}

type upstreamHandler struct {
	c *controller
}

func (uh *upstreamHandler) OnAdd(updatedList []*v1.Upstream, obj *v1.Upstream) {
	uh.OnUpdate(updatedList, obj)
}

func (uh *upstreamHandler) OnUpdate(updatedList []*v1.Upstream, obj *v1.Upstream) {
	for _, h := range uh.c.handlers {
		h.Update(obj)
	}
}

func (uh *upstreamHandler) OnDelete(updatedList []*v1.Upstream, obj *v1.Upstream) {
	for _, h := range uh.c.handlers {
		// FIXME (ashish) - file based storage have nil obj
		h.Remove(obj.Name)
	}
}
