package mocks

import (
	"context"
)

type MocksApiServer struct {
	cache Cache
}

func NewMocksApiServer(cache Cache) ApiServerServer {
	return &MocksApiServer{
		cache: cache,
	}
}

func (s *MocksApiServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return &RegisterResponse{}, s.cache.Register()
}

func (s *MocksApiServer) Read(ctx context.Context, req *ReadRequest) (*ReadResponse, error) {
	
}

func (s *MocksApiServer) Write(ctx context.Context, req *WriteRequest) (*WriteResponse, error) {}

func (s *MocksApiServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {}

func (s *MocksApiServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {}

func (s *MocksApiServer) Watch(req *WatchRequest, srv ApiServer_WatchServer) error {}
