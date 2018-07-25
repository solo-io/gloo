package apiserver

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
)

type ApiServer struct {
	resourceClients map[string]clients.ResourceClient
}

func NewApiServer(resourceClients ...clients.ResourceClient) ApiServerServer {
	mapped := make(map[string]clients.ResourceClient)
	for _, rc := range resourceClients {
		mapped[rc.Kind()] = rc
	}
	return &ApiServer{
		resourceClients: mapped,
	}
}

func (s *ApiServer) resourceClient(kind string) (clients.ResourceClient, error) {
	rc, ok := s.resourceClients[kind]
	if !ok {
		return nil, errors.Errorf("no resource client registered for kind %s", kind)
	}
	return rc, nil
}

func (s *ApiServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	for _, rc := range s.resourceClients {
		if err := rc.Register(); err != nil {
			return nil, errors.Wrapf(err, "failed to register client %v", rc.Kind())
		}
	}
	return &RegisterResponse{}, nil
}

func (s *ApiServer) Read(ctx context.Context, req *ReadRequest) (*ReadResponse, error) {
	rc, err := s.resourceClient(req.Kind)
	if err != nil {
		return nil, err
	}
	resource, err := rc.Read(req.Name, clients.ReadOpts{
		Namespace: req.Namespace,
		Ctx:       contextutils.WithLogger(ctx, "apiserver.read"),
	})
	data, err := protoutils.MarshalStruct(resource)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal resource")
	}
	return &ReadResponse{
		Resource: &Resource{
			Kind: rc.Kind(),
			Data: data,
		},
	}, nil
}

func (s *ApiServer) Write(ctx context.Context, req *WriteRequest) (*WriteResponse, error) {
	rc, err := s.resourceClient(req.Resource.Kind)
	if err != nil {
		return nil, err
	}
	resource := rc.NewResource()
	if err := protoutils.UnmarshalStruct(req.Resource.Data, resource); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal resource %v", rc.Kind())
	}
	resource, err = rc.Write(resource, clients.WriteOpts{
		OverwriteExisting: req.OverwriteExisting,
		Ctx:               contextutils.WithLogger(ctx, "apiserver.write"),
	})
	data, err := protoutils.MarshalStruct(resource)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal resource")
	}
	return &WriteResponse{
		Resource: &Resource{
			Kind: rc.Kind(),
			Data: data,
		},
	}, nil
}

func (s *ApiServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {}

func (s *ApiServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {}

func (s *ApiServer) Watch(*WatchRequest, ApiServer_WatchServer) error {}
