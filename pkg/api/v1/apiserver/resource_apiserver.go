package apiserver

import (
	"context"
)

type MocksApiServer struct{}

func NewMocksApiServer() ApiServerServer {
	return &MocksApiServer{}
}

func (s *MocksApiServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {}
func (s *MocksApiServer) Read(context.Context, *ReadRequest) (*ReadResponse, error)             {}
func (s *MocksApiServer) Write(context.Context, *WriteRequest) (*WriteResponse, error)          {}
func (s *MocksApiServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)       {}
func (s *MocksApiServer) List(context.Context, *ListRequest) (*ListResponse, error)             {}
func (s *MocksApiServer) Watch(*WatchRequest, ApiServer_WatchServer) error                      {}
