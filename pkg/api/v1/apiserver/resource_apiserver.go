package apiserver

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	"google.golang.org/grpc"
)

type ApiServer struct {
	resourceClients map[string]clients.ResourceClient
}

func NewApiServer(s *grpc.Server, resourceClients ...clients.ResourceClient) ApiServerServer {
	mapped := make(map[string]clients.ResourceClient)
	for _, rc := range resourceClients {
		mapped[rc.Kind()] = rc
	}
	srv := &ApiServer{
		resourceClients: mapped,
	}
	RegisterApiServerServer(s, srv)
	return srv
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

func (s *ApiServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	rc, err := s.resourceClient(req.Kind)
	if err != nil {
		return nil, err
	}
	if err := rc.Delete(req.Name, clients.DeleteOpts{
		IgnoreNotExist: req.IgnoreNotExist,
		Namespace:      req.Namespace,
		Ctx:            contextutils.WithLogger(ctx, "apiserver.delete"),
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to delete resource %v", req.Kind)
	}
	return &DeleteResponse{}, nil
}

func (s *ApiServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	rc, err := s.resourceClient(req.Kind)
	if err != nil {
		return nil, err
	}
	resourceList, err := rc.List(clients.ListOpts{
		Namespace: req.Namespace,
		Ctx:       contextutils.WithLogger(ctx, "apiserver.read"),
	})
	var resourceListResponse []*Resource
	for _, resource := range resourceList {
		data, err := protoutils.MarshalStruct(resource)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal resource %v", req.Kind)
		}
		resourceListResponse = append(resourceListResponse, &Resource{
			Kind: rc.Kind(),
			Data: data,
		})
	}
	return &ListResponse{
		ResourceList: resourceListResponse,
	}, nil
}

func (s *ApiServer) Watch(req *WatchRequest, watch ApiServer_WatchServer) error {
	rc, err := s.resourceClient(req.Kind)
	if err != nil {
		return err
	}
	ctx := contextutils.WithLogger(watch.Context(), "apiserver.read")
	resourceWatch, errs, err := rc.Watch(clients.WatchOpts{
		RefreshRate: req.SyncFrequency,
		Namespace:   req.Namespace,
		Ctx:         ctx,
	})
	for {
		select {
		case resourceList := <-resourceWatch:
			var resourceListResponse []*Resource
			for _, resource := range resourceList {
				data, err := protoutils.MarshalStruct(resource)
				if err != nil {
					return errors.Wrapf(err, "failed to marshal resource %v", req.Kind)
				}
				resourceListResponse = append(resourceListResponse, &Resource{
					Kind: rc.Kind(),
					Data: data,
				})
			}
			if err := watch.Send(&ListResponse{
				ResourceList: resourceListResponse,
			}); err != nil {
				return errors.Wrapf(err, "failed to send list response on watch")
			}
		case err := <-errs:
			return errors.Wrapf(err, "error during %v watch", req.Kind)
		case <-ctx.Done():
			return nil
		}
	}
}
