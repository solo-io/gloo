package v1helpers

import (
	"context"
	"fmt"
	"net"

	structpb "github.com/golang/protobuf/ptypes/struct"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
)

// this code was copied from https://github.com/solo-io/ext-auth-service/blob/v0.16.4/pkg/config/passthrough/test_utils/fake_grpc_auth_server.go
const (
	// This is the key used to store PassThrough state in the CheckRequest and CheckResponse metadata
	MetadataStateKey  = "solo.auth.passthrough"
	MetadataConfigKey = "solo.auth.passthrough.config"
)

type AuthChecker func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error)

// Represents an external auth service that always returns the same response/error
func NewGrpcAuthServerWithResponse(response *envoy_service_auth_v3.CheckResponse, err error) *GrpcAuthServer {
	return &GrpcAuthServer{
		AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
			return response, err
		},
	}
}

// Represents an external auth service that only returns ok if all the requiredKeys are passed in the filter metadata
func NewGrpcAuthServerWithRequiredMetadata(requiredKeys []string) *GrpcAuthServer {
	return &GrpcAuthServer{
		AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
			if passThroughFilterMetadata, ok := req.GetAttributes().GetMetadataContext().GetFilterMetadata()[MetadataStateKey]; ok {
				passThroughFields := passThroughFilterMetadata.GetFields()
				for _, requiredKey := range requiredKeys {
					if _, ok := passThroughFields[requiredKey]; !ok {
						// Required key was not passed in FilterMetadata, fail request
						return DeniedResponse(), nil
					}
				}
				// All keys were passed in FilterMetadata, succeed request
				return OkResponse(), nil
			}
			// No passthrough properties were sent in FilterMetadata, fail request
			return DeniedResponse(), nil
		},
	}
}

// Represents a runnable instance of an external auth service that performs a custom auth check
type GrpcAuthServer struct {
	AuthChecker AuthChecker

	port    int
	address string
	server  *grpc.Server
}

func (s *GrpcAuthServer) Check(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	return s.AuthChecker(ctx, req)
}

func (s *GrpcAuthServer) Start(port int) error {
	s.Stop()

	srv := grpc.NewServer()
	envoy_service_auth_v3.RegisterAuthorizationServer(srv, s)

	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.server = srv
	s.port = port
	s.address = address

	go func() {
		defer GinkgoRecover()
		err := srv.Serve(listener)
		Expect(err).ToNot(HaveOccurred())
	}()

	return nil
}

func (s *GrpcAuthServer) Stop() {
	if s.server != nil {
		s.server.Stop()
	}
}

func (s *GrpcAuthServer) GetAddress() string {
	return s.address
}

func OkResponse() *envoy_service_auth_v3.CheckResponse {
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.OK),
		},
		HttpResponse: &envoy_service_auth_v3.CheckResponse_OkResponse{
			OkResponse: &envoy_service_auth_v3.OkHttpResponse{},
		},
		DynamicMetadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		},
	}
}

func OkResponseWithDynamicMetadata(passThroughDynamicMetadata *structpb.Struct) *envoy_service_auth_v3.CheckResponse {
	response := OkResponse()
	response.DynamicMetadata.Fields[MetadataStateKey] = &structpb.Value{
		Kind: &structpb.Value_StructValue{
			StructValue: passThroughDynamicMetadata,
		},
	}
	return response
}

func DeniedResponse() *envoy_service_auth_v3.CheckResponse {
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &envoy_service_auth_v3.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
				Status: &envoy_type_v3.HttpStatus{
					Code: envoy_type_v3.StatusCode_Unauthorized,
				},
				Body: "PERMISSION_DENIED",
			},
		},
		DynamicMetadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		},
	}
}
