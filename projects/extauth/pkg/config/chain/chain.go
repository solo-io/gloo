package chain

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	errors "github.com/rotisserie/eris"
	"go.uber.org/atomic"

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
	SetAuthorizer(expr ast.Expr) error
}

var _ AuthServiceChain = &authServiceChain{}

type authServiceWithName struct {
	name        string
	authService api.AuthService
}

// Used to wrap a collection of auth services and expose them as a single AuthService Implementation
type authServiceChain struct {
	authServices []authServiceWithName
	started      atomic.Bool
	names        []string
	authorizer   authorizer
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
	if s.started.Load() {
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
	s.started.Store(true)
	return nil
}

func (s *authServiceChain) Authorize(ctx context.Context, request *api.AuthorizationRequest) (*api.AuthorizationResponse, error) {
	var err error
	// Base case: allow request
	finalResponse := api.AuthorizedResponse()
	finalResponse, err = s.authorizer.Authorize(ctx, request, finalResponse)
	if err != nil {
		return nil, err
	}

	s.started.Store(true)

	return finalResponse, nil
}

// call this function after all auth services have been added to the chain to set the authorizer
func (s *authServiceChain) SetAuthorizer(expr ast.Expr) error {
	if expr == nil {
		// if no expr provided, default behavior is to AND the entire auth chain together
		var authServiceNames []string
		for _, as := range s.authServices {
			authServiceNames = append(authServiceNames, as.name)
		}

		andStr := strings.Join(authServiceNames, " && ")
		var parseErr error
		expr, parseErr = parser.ParseExpr(andStr)
		if parseErr != nil {
			return parseErr // should never happen
		}
	}

	var err error
	s.authorizer, err = newAuthorizer(expr, s)
	if err != nil {
		return err
	}

	return nil
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

var _ authorizer = &andAuthorizerImpl{}
var _ authorizer = &orAuthorizerImpl{}
var _ authorizer = &notAuthorizerImpl{}
var _ authorizer = &identityAuthorizerImpl{}

type authorizer interface {
	// finalResp is passed to recursively allow each request to merge headers into the final successful response.
	// the return response should be used instead of finalResp upon completion -- failed responses are returned but don't modify finalResp
	Authorize(ctx context.Context, request *api.AuthorizationRequest, finalResp *api.AuthorizationResponse) (*api.AuthorizationResponse, error)
}

type andAuthorizerImpl struct {
	x authorizer
	y authorizer
	c *authServiceChain
}

type orAuthorizerImpl struct {
	x authorizer
	y authorizer
	c *authServiceChain
}

type notAuthorizerImpl struct {
	x authorizer
	c *authServiceChain
}

func (a *andAuthorizerImpl) Authorize(ctx context.Context, request *api.AuthorizationRequest, finalResp *api.AuthorizationResponse) (*api.AuthorizationResponse, error) {
	x, err := a.x.Authorize(ctx, request, finalResp)
	if err != nil {
		return nil, err
	}
	if x.CheckResponse.Status.Code != int32(rpc.OK) {
		return x, nil
	}
	y, err := a.y.Authorize(ctx, request, finalResp)
	if err != nil {
		return nil, err
	}
	if y.CheckResponse.Status.Code != int32(rpc.OK) {
		return y, nil
	}
	return finalResp, nil
}

func newAndAuthorizer(binary *ast.BinaryExpr, c *authServiceChain) (andAuthorizer *andAuthorizerImpl, err error) {
	andAuth := &andAuthorizerImpl{}
	if _, ok := binary.Y.(*ast.ParenExpr); ok {
		andAuth.x, err = newAuthorizer(binary.Y, c)
		if err != nil {
			return nil, err
		}
		andAuth.y, err = newAuthorizer(binary.X, c)
		if err != nil {
			return nil, err
		}
	} else {
		andAuth.x, err = newAuthorizer(binary.X, c)
		if err != nil {
			return nil, err
		}
		andAuth.y, err = newAuthorizer(binary.Y, c)
		if err != nil {
			return nil, err
		}
	}
	return andAuth, nil
}

func (o *orAuthorizerImpl) Authorize(ctx context.Context, request *api.AuthorizationRequest, finalResp *api.AuthorizationResponse) (*api.AuthorizationResponse, error) {
	x, err := o.x.Authorize(ctx, request, finalResp)
	if err != nil {
		return nil, err
	}
	if x.CheckResponse.Status.Code == int32(rpc.OK) {
		return x, nil
	}
	y, err := o.y.Authorize(ctx, request, finalResp)
	if err != nil {
		return nil, err
	}
	if y.CheckResponse.Status.Code == int32(rpc.OK) {
		return y, nil
	}
	// unclear whether to return denied response from x or y, so do neither
	// in the future we may want to implement some header merging
	return api.UnauthorizedResponse(), nil
}

func newOrAuthorizer(binary *ast.BinaryExpr, c *authServiceChain) (orAuthorizer *orAuthorizerImpl, err error) {
	orAuth := &orAuthorizerImpl{}
	if _, ok := binary.Y.(*ast.ParenExpr); ok {
		orAuth.x, err = newAuthorizer(binary.Y, c)
		if err != nil {
			return nil, err
		}
		orAuth.y, err = newAuthorizer(binary.X, c)
		if err != nil {
			return nil, err
		}
	} else {
		orAuth.x, err = newAuthorizer(binary.X, c)
		if err != nil {
			return nil, err
		}
		orAuth.y, err = newAuthorizer(binary.Y, c)
		if err != nil {
			return nil, err
		}
	}
	return orAuth, nil
}

func (n *notAuthorizerImpl) Authorize(ctx context.Context, request *api.AuthorizationRequest, finalResp *api.AuthorizationResponse) (*api.AuthorizationResponse, error) {
	x, err := n.x.Authorize(ctx, request, finalResp)
	if err != nil {
		return nil, err
	}
	if x.CheckResponse.Status.Code != int32(rpc.OK) {
		return finalResp, nil // finalResp is assumed to be an allowed response
	}
	return api.UnauthorizedResponse(), nil
}

func newNotAuthorizer(unary *ast.UnaryExpr, c *authServiceChain) (notAuthorizer *notAuthorizerImpl, err error) {
	notAuth := &notAuthorizerImpl{}
	notAuth.x, err = newAuthorizer(unary.X, c)
	if err != nil {
		return nil, err
	}
	return notAuth, nil
}

func newBinaryAuthorizer(binary *ast.BinaryExpr, c *authServiceChain) (authorizer authorizer, err error) {
	switch binary.Op {
	case token.LAND:
		return newAndAuthorizer(binary, c)
	case token.LOR:
		return newOrAuthorizer(binary, c)
	}

	return nil, errors.New(fmt.Sprintf("invalid binary expression, %T. Expected AND (&&) or OR (||) binary", binary))
}

func newUnaryAuthorizer(unary *ast.UnaryExpr, c *authServiceChain) (authorizer authorizer, err error) {
	switch unary.Op {
	case token.NOT:
		return newNotAuthorizer(unary, c)
	}
	return nil, errors.New(fmt.Sprintf("invalid unary expression, %T. Expected NOT (!) unary", unary))
}

func newAuthorizer(expr ast.Expr, c *authServiceChain) (authorizer, error) {
	if binary, ok := expr.(*ast.BinaryExpr); ok {
		return newBinaryAuthorizer(binary, c)
	}
	if unary, ok := expr.(*ast.UnaryExpr); ok {
		return newUnaryAuthorizer(unary, c)
	}
	if paren, ok := expr.(*ast.ParenExpr); ok {
		return newAuthorizer(paren.X, c)
	}
	if ident, ok := expr.(*ast.Ident); ok {
		return newIdentityAuthorizer(ident, c), nil
	}
	return nil, errors.New(fmt.Sprintf("Unexpected expression, %T. Expected a binary (i.e., && or ||), unary (i.e., !), parenthesis, or identity (i.e., a-zA-Z string) expression", expr))
}

func newIdentityAuthorizer(identifier *ast.Ident, c *authServiceChain) (identityAuthorizer *identityAuthorizerImpl) {
	return &identityAuthorizerImpl{
		chain:      c,
		identifier: identifier,
	}
}

type identityAuthorizerImpl struct {
	identifier *ast.Ident
	chain      *authServiceChain
}

func (i *identityAuthorizerImpl) Authorize(ctx context.Context, request *api.AuthorizationRequest, finalResp *api.AuthorizationResponse) (*api.AuthorizationResponse, error) {

	for _, as := range i.chain.authServices {
		if as.name == i.identifier.Name {
			response, err := as.authService.Authorize(ctx, request)
			if err != nil {
				contextutils.LoggerFrom(ctx).Errorw("Error during authService execution",
					zap.Any("authService", as.name),
					zap.Any("error", err),
				)
				return nil, err
			}

			// If response is not OK return without executing any further authService. Nil status means OK
			if status := response.CheckResponse.Status; status != nil && status.Code != int32(rpc.OK) {
				contextutils.LoggerFrom(ctx).Debugw("Access denied by auth authService", zap.Any("authService", as.name))
				return response, nil
			}

			// Response is OK, merge headers into previous request
			responseHeaders := mergeHeaders(finalResp.CheckResponse.GetOkResponse(), response.CheckResponse.GetOkResponse())

			// If no new user id given, merge the one from the last response
			if response.UserInfo.UserID == "" {
				response.UserInfo = finalResp.UserInfo
			}

			if responseHeaders != nil {
				response.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_OkResponse{
					OkResponse: responseHeaders,
				}
			}

			*finalResp = *response

			return response, nil
		}
	}

	var asNames []string
	for _, as := range i.chain.authServices {
		asNames = append(asNames, as.name)
	}
	return nil, errors.New(fmt.Sprintf("No matching auth service for %s in %v", i.identifier.Name, asNames))
}
