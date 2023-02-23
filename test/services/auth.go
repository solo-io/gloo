package services

// This is a replica of services defined in github.com/solo-io/ext-auth-service
// We have duplicated them in solo-projects while upgrading Ginkgo to use v2, since the
// ext-auth-service repository still uses Ginkgo v1 and that leads to panics in the tests
// https://onsi.github.io/ginkgo/MIGRATING_TO_V2#can-i-mix-ginkgo-v1-and-ginkgo-v2
//
// We can remove this when we resolve: https://github.com/solo-io/solo-projects/issues/4590

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"google.golang.org/grpc/metadata"

	"github.com/solo-io/ext-auth-service/pkg/config/oauth/token_validation"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/user_info"
	"k8s.io/apimachinery/pkg/util/sets"

	passthrough "github.com/solo-io/ext-auth-service/pkg/config/passthrough/grpc"

	structpb "github.com/golang/protobuf/ptypes/struct"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
)

var (
	tokenHeaderRegex = regexp.MustCompile("^Bearer (.+)$")
)

type AuthEndpoints struct {
	// both of these are something like "/user-info", or "/introspection"
	UserInfoEndpoint           string
	TokenIntrospectionEndpoint string
}

type AuthHandlers struct {
	TokenIntrospectionHandler func(writer http.ResponseWriter, request *http.Request)
}

type AuthServer struct {
	srv             *http.Server
	authEndpoints   *AuthEndpoints
	authHandlers    *AuthHandlers
	baseAddress     string
	validTokens     sets.String
	tokenToUserInfo map[string]user_info.UserInfo
}

// NewAuthServer provides user_info and token introspection endpoints
// every arg must be non-nil
// base address is something like ":8000" to run on localhost port 8000
func NewAuthServer(baseAddress string, endpoints *AuthEndpoints, handlers *AuthHandlers, validTokens sets.String, tokenToUserInfo map[string]user_info.UserInfo) *AuthServer {
	return &AuthServer{
		authEndpoints:   endpoints,
		authHandlers:    handlers,
		baseAddress:     baseAddress,
		validTokens:     validTokens,
		tokenToUserInfo: tokenToUserInfo,
	}
}

func (a *AuthServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = a.srv.Shutdown(ctx)
}

func (a *AuthServer) Start() {
	router := mux.NewRouter().StrictSlash(true)

	if a.authEndpoints.UserInfoEndpoint != "" {
		router.HandleFunc(a.authEndpoints.UserInfoEndpoint, func(writer http.ResponseWriter, request *http.Request) {
			authHeader := request.Header.Get("Authorization")
			matches := tokenHeaderRegex.FindStringSubmatch(authHeader)
			if len(matches) != 2 {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			token := matches[1]

			if userInfo, ok := a.tokenToUserInfo[token]; ok {
				bytes, err := json.Marshal(userInfo)
				if err != nil {
					panic(err)
				}
				writer.Header().Set("Content-Type", "application/json; charset=utf-8")
				writer.Write(bytes)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
			}
		})
	}

	if a.authEndpoints.TokenIntrospectionEndpoint != "" {
		if a.authHandlers.TokenIntrospectionHandler != nil {
			router.HandleFunc(a.authEndpoints.TokenIntrospectionEndpoint, a.authHandlers.TokenIntrospectionHandler)
		} else {
			// Default handler doesn't require authentication
			router.HandleFunc(a.authEndpoints.TokenIntrospectionEndpoint, func(writer http.ResponseWriter, request *http.Request) {
				err := request.ParseForm()
				if err != nil {
					panic(err)
				}

				token := request.Form.Get("token")

				response := &token_validation.IntrospectionResponse{}
				if a.validTokens.Has(token) {
					response.Active = true
				}

				bytes, err := json.Marshal(response)
				if err != nil {
					panic(err)
				}
				writer.Write(bytes)
			})
		}
	}

	a.srv = &http.Server{
		Addr:    a.baseAddress,
		Handler: router,
	}

	go func() {
		defer GinkgoRecover()
		err := a.srv.ListenAndServe()
		if err != http.ErrServerClosed {
			Expect(err).NotTo(HaveOccurred())
		}
	}()
}

type AuthChecker func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error)

// NewGrpcAuthServerWithResponse Represents an external auth service that always returns the same response/error
func NewGrpcAuthServerWithResponse(response *envoy_service_auth_v3.CheckResponse, err error) *GrpcAuthServer {
	return &GrpcAuthServer{
		AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
			return response, err
		},
	}
}

// NewGrpcAuthServerWithDelayedResponse Represents an external auth service that always returns the same response/error after the a 1s delay
func NewGrpcAuthServerWithDelayedResponse(response *envoy_service_auth_v3.CheckResponse, err error) *GrpcAuthServer {
	return &GrpcAuthServer{
		AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
			time.Sleep(1 * time.Second)
			return response, err
		},
	}
}

// NewGrpcAuthServerWithRequiredMetadata Represents an external auth service that only returns ok if all the requiredKeys are passed in the filter metadata
func NewGrpcAuthServerWithRequiredMetadata(requiredKeys []string) *GrpcAuthServer {
	return &GrpcAuthServer{
		AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
			if passThroughFilterMetadata, ok := req.GetAttributes().GetMetadataContext().GetFilterMetadata()[passthrough.MetadataStateKey]; ok {
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

func NewGrpcAuthServerWithTracingRequired() *GrpcAuthServer {
	return &GrpcAuthServer{
		AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
			inMD, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				// No incoming metadata at all, fail request
				return DeniedResponse(), nil
			}
			if len(inMD.Get("x-b3-traceid")) > 0 {
				// Incoming metadata contained expected tracing data, succeed request
				return OkResponse(), nil
			}
			// Expected tracing header not found, fail request
			return DeniedResponse(), nil
		},
	}
}

// GrpcAuthServer Represents a runnable instance of an external auth service that performs a custom auth check
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

// OkResponse returns an OK HTTP Response with a status of 200 OK
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

// OkResponseWithDynamicMetadata returns an OK HTTP Response with PassThrough dynamic metadata set
func OkResponseWithDynamicMetadata(passThroughDynamicMetadata *structpb.Struct) *envoy_service_auth_v3.CheckResponse {
	response := OkResponse()
	response.DynamicMetadata.Fields[passthrough.MetadataStateKey] = &structpb.Value{
		Kind: &structpb.Value_StructValue{
			StructValue: passThroughDynamicMetadata,
		},
	}
	return response
}

// DeniedResponse returns a denied HTTP Response with a status of 401 Unauthorized
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

// ServerErrorResponse returns a denied HTTP Response with a status of 503 Service Unavailable
func ServerErrorResponse() *envoy_service_auth_v3.CheckResponse {
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.UNAVAILABLE),
		},
		HttpResponse: &envoy_service_auth_v3.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
				Status: &envoy_type_v3.HttpStatus{
					Code: envoy_type_v3.StatusCode_ServiceUnavailable,
				},
				Body: "SERVICE_UNAVAILABLE",
			},
		},
		DynamicMetadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		},
	}
}
