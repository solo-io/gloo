package extauth_test_server

import (
	"fmt"
	"sync/atomic"

	"github.com/golang/protobuf/ptypes/wrappers"
	oauth_utils "github.com/solo-io/ext-auth-service/pkg/config/oauth/test_utils"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/user_info"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	IntrospectionAccessToken = "valid-access-token"
)

type TokenIntrospectionServer struct {
	AuthHandlers *oauth_utils.AuthHandlers
	s            *oauth_utils.AuthServer
	port         uint32
}

func NewTokenIntrospectionServer() *TokenIntrospectionServer {
	port := atomic.AddUint32(&baseOauth2Port, 1) + uint32(parallel.GetPortOffset())
	return &TokenIntrospectionServer{
		port:         port,
		AuthHandlers: &oauth_utils.AuthHandlers{},
	}
}

func (t *TokenIntrospectionServer) Start() {
	if t.port == 0 {
		// If port has not been set yet, generate a new one.
		// It's possible that the port has been set with `NewTokenIntrospectionServer`.
		t.port = atomic.AddUint32(&baseOauth2Port, 1) + uint32(parallel.GetPortOffset())
	}

	t.s = oauth_utils.NewAuthServer(
		fmt.Sprintf(":%d", t.port),
		&oauth_utils.AuthEndpoints{
			TokenIntrospectionEndpoint: "/introspection",
			UserInfoEndpoint:           "/userinfo",
		},
		t.AuthHandlers,
		sets.NewString("valid-access-token"),
		map[string]user_info.UserInfo{},
	)

	t.s.Start()
}

func (t *TokenIntrospectionServer) Stop() {
	t.s.Stop()
}

func (t *TokenIntrospectionServer) GetOauthTokenIntrospectionConfig(clientId string, clientSecretRef *core.ResourceRef, disableClientSecret bool) *extauth.OAuth2_AccessTokenValidation {
	return &extauth.OAuth2_AccessTokenValidation{
		AccessTokenValidation: &extauth.AccessTokenValidation{
			ValidationType: &extauth.AccessTokenValidation_Introspection{
				Introspection: &extauth.IntrospectionValidation{
					IntrospectionUrl:    fmt.Sprintf("http://localhost:%d/introspection", t.port),
					ClientId:            clientId,
					ClientSecretRef:     clientSecretRef,
					DisableClientSecret: &wrappers.BoolValue{Value: disableClientSecret},
				},
			},
			UserinfoUrl:  fmt.Sprintf("http://localhost:%d/userinfo", t.port),
			CacheTimeout: nil,
		},
	}
}

func (t *TokenIntrospectionServer) GetOauthTokenIntrospectionUrlConfig() *extauth.OAuth2_AccessTokenValidation {
	return &extauth.OAuth2_AccessTokenValidation{
		AccessTokenValidation: &extauth.AccessTokenValidation{
			ValidationType: &extauth.AccessTokenValidation_IntrospectionUrl{
				IntrospectionUrl: fmt.Sprintf("http://localhost:%d/introspection", t.port),
			},
			UserinfoUrl:  fmt.Sprintf("http://localhost:%d/userinfo", t.port),
			CacheTimeout: nil,
		},
	}
}
