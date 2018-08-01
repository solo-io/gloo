package rbac

import (
	"context"
	"fmt"

	"github.com/ory/ladon"
	"github.com/solo-io/solo-kit/pkg/api/v1/apiserver"
	"github.com/solo-io/solo-kit/pkg/api/v1/rbac/policy"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

type Callbacks struct {
	resourceType resources.Resource
	warden       ladon.Warden
}

var _ apiserver.Callbacks = &Callbacks{}

func (cb *Callbacks) OnRegister(ctx context.Context, req *apiserver.RegisterRequest) (*apiserver.RegisterResponse, error) {
	accessRequest := &ladon.Request{
		Resource: cb.resourceName(),
	}
}

func (cb *Callbacks) OnRead(ctx context.Context, req *apiserver.ReadRequest) (*apiserver.ReadResponse, error) {
}

func (cb *Callbacks) OnWrite(ctx context.Context, req *apiserver.WriteRequest) (*apiserver.WriteResponse, error) {
}

func (cb *Callbacks) OnDelete(ctx context.Context, req *apiserver.DeleteRequest) (*apiserver.DeleteResponse, error) {
}

func (cb *Callbacks) OnList(ctx context.Context, req *apiserver.ListRequest) (*apiserver.ListResponse, error) {
}

func (cb *Callbacks) OnWatch(req *apiserver.WatchRequest, watch apiserver.ApiServer_WatchServer) error {
}

func (cb *Callbacks) accessRequestFromContext(ctx context.Context, capability policy.Capability) *ladon.Request {
	return &ladon.Request{}
}

func (cb *Callbacks) resourceName(namespace, name string) string {
	return fmt.Sprintf("%v.%v.%v", resources.Kind(cb.resourceType), namespace, name)
}
