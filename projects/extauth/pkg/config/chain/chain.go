package chain

import (
	"context"
	"fmt"

	"github.com/solo-io/go-utils/errors"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
)

var DuplicateAuthServiceNameError = func(name string) error {
	return errors.Errorf("duplicate auth service name [%s]", name)
}

func NewAuthServiceChain() *authServiceChain {
	return &authServiceChain{}
}

type AuthServiceChain interface {
	api.AuthService
	AddAuthService(name string, authService api.AuthService) error
	ListAuthServices() []api.AuthService
}

var _ AuthServiceChain = &authServiceChain{}

type authServiceWithName struct {
	name        string
	authService api.AuthService
}

// Used to wrap a collection of auth services and expose them as a single AuthService Implementation
type authServiceChain struct {
	authServices []authServiceWithName
	started      bool
	names        []string
}

// Returns true if the chain contains a auth service with the given name
func (s *authServiceChain) contains(name string) bool {
	for _, existingName := range s.names {
		if existingName == name {
			return true
		}
	}
	return false
}

func (s *authServiceChain) AddAuthService(name string, authService api.AuthService) error {
	if s.started {
		panic("cannot add authService to started authServiceChain!")
	}
	if s.contains(name) {
		return DuplicateAuthServiceNameError(name)
	}
	s.authServices = append(s.authServices, authServiceWithName{
		name:        name,
		authService: authService,
	})
	// Pre-compute the list of names so we don't have to loop during actual requests
	s.names = append(s.names, name)
	return nil
}

func (s *authServiceChain) ListAuthServices() (out []api.AuthService) {
	for _, svc := range s.authServices {
		out = append(out, svc.authService)
	}
	return
}

func (s *authServiceChain) Start(ctx context.Context) error {
	for _, p := range s.authServices {
		if err := p.authService.Start(ctx); err != nil {
			return err
		}
	}
	s.started = true
	return nil
}

func (s *authServiceChain) Authorize(ctx context.Context, request *api.AuthorizationRequest) (*api.AuthorizationResponse, error) {

	// Base case: allow request
	lastResponse := api.AuthorizedResponse()

	for i, p := range s.authServices {

		response, err := p.authService.Authorize(ctx, request)
		if err != nil {
			contextutils.LoggerFrom(ctx).Errorw("Error during authService execution",
				zap.Any("authService", p.name),
				zap.Any("error", err),
			)
			return nil, err
		}

		// If response is not OK return without executing any further authService. Nil status means OK
		if status := response.CheckResponse.Status; status != nil && status.Code != int32(rpc.OK) {
			contextutils.LoggerFrom(ctx).Debugw("Access denied by auth authService", zap.Any("authService", p.name))
			if i < len(s.authServices)-1 {
				contextutils.LoggerFrom(ctx).Debugw("Skipping execution of following authServices",
					zap.Any("skippedauthServices", s.names[i+1:]))
			}
			return response, nil
		}

		// Response is OK, merge headers into previous request
		responseHeaders := mergeHeaders(lastResponse.CheckResponse.GetOkResponse(), response.CheckResponse.GetOkResponse())

		// If no new user id given, merge the one from the last response
		if response.UserInfo.UserID == "" {
			response.UserInfo = lastResponse.UserInfo
		}

		if responseHeaders != nil {
			response.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_OkResponse{
				OkResponse: responseHeaders,
			}
		}

		lastResponse = response
	}

	s.started = true

	return lastResponse, nil
}

// This gets called only if both the last and new responses are OK
func mergeHeaders(last, new *envoyauthv2.OkHttpResponse) *envoyauthv2.OkHttpResponse {

	// If the new response does not have any additional header information, use the ones from the previous one
	if new == nil {
		return last
	}

	// Default if last response did not have any additional header information
	if last == nil {
		return new
	}

	// Index last response headers
	lastHeadersMap := map[string]*core.HeaderValueOption{}
	for _, h := range last.Headers {
		// Clone so we don't modify the input
		lastHeadersMap[h.Header.Key] = h
	}

	// Add new headers to last ones, overwriting if necessary
	for _, newHeader := range new.Headers {

		lastHeader, ok := lastHeadersMap[newHeader.Header.Key]

		// Header was not present in last response OR new header should overwrite old one
		if !ok || newHeader.Append.Value == false {
			lastHeadersMap[newHeader.Header.Key] = newHeader
			continue
		}

		// Append header value to the previous one
		lastHeader.Header.Value = fmt.Sprintf("%s, %s", lastHeader.Header.Value, newHeader.Header.Value)
		lastHeader.Append = newHeader.Append
	}

	var result []*core.HeaderValueOption
	for _, headerValueOption := range lastHeadersMap {
		result = append(result, headerValueOption)
	}

	return &envoyauthv2.OkHttpResponse{
		Headers: result,
	}
}
