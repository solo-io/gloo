package plugins

import (
	"context"
	"fmt"

	"github.com/solo-io/go-utils/errors"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
)

var DuplicatePluginNameError = func(name string) error {
	return errors.Errorf("duplicate plugin name [%s]", name)
}

func NewExtAuthPluginChain() *pluginChain {
	return &pluginChain{}
}

type ExtAuthPluginChain interface {
	api.AuthService
	AddPlugin(name string, plugin api.AuthService) error
}

type pluginWithName struct {
	name   string
	plugin api.AuthService
}

// Used to wrap a collection of auth plugins and expose them as a single AuthService Implementation
type pluginChain struct {
	plugins []pluginWithName
	started bool
	names   []string
}

// Returns true if the chain contains a plugin with the given name
func (s *pluginChain) contains(name string) bool {
	for _, existingName := range s.names {
		if existingName == name {
			return true
		}
	}
	return false
}

func (s *pluginChain) AddPlugin(name string, plugin api.AuthService) error {
	if s.started {
		panic("cannot add plugin to started plugin chain!")
	}
	if s.contains(name) {
		return DuplicatePluginNameError(name)
	}
	s.plugins = append(s.plugins, pluginWithName{
		name:   name,
		plugin: plugin,
	})
	// Pre-compute the list of names so we don't have to loop during actual requests
	s.names = append(s.names, name)
	return nil
}

func (s *pluginChain) Start(ctx context.Context) error {
	for _, p := range s.plugins {
		if err := p.plugin.Start(ctx); err != nil {
			return err
		}
	}
	s.started = true
	return nil
}

func (s *pluginChain) Authorize(ctx context.Context, request *envoyauthv2.CheckRequest) (*api.AuthorizationResponse, error) {

	// Base case: allow request
	lastResponse := api.AuthorizedResponse()

	for i, p := range s.plugins {

		response, err := p.plugin.Authorize(ctx, request)
		if err != nil {
			contextutils.LoggerFrom(ctx).Errorw("Error during plugin execution",
				zap.Any("plugin", p.name),
				zap.Any("error", err),
			)
			return nil, err
		}

		// If response is not OK return without executing any further plugin. Nil status means OK
		if status := response.CheckResponse.Status; status != nil && status.Code != int32(rpc.OK) {
			contextutils.LoggerFrom(ctx).Infow("Access denied by auth plugin", zap.Any("plugin", p.name))
			if i < len(s.plugins)-1 {
				contextutils.LoggerFrom(ctx).Debugw("Skipping execution of following plugins",
					zap.Any("skippedPlugins", s.names[i+1:]))
			}
			return response, nil
		}

		// Response is OK, merge headers into previous request
		responseHeaders := mergeHeaders(lastResponse.CheckResponse.GetOkResponse(), response.CheckResponse.GetOkResponse())

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
