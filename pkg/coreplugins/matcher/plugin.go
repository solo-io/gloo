package matcher

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugin"
)

const (
	eventPath        = "/events"
	headerXEventType = "X-Event-Type"
)

type Plugin struct{}

func (p *Plugin) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *Plugin) ProcessRoute(_ *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	switch matcher := in.Matcher.(type) {
	case *v1.Route_EventMatcher:
		return createEventMatcher(matcher.EventMatcher, out)
	case *v1.Route_RequestMatcher:
		return createRequestMatcher(matcher.RequestMatcher, out)
	}
	return errors.New("invalid or unspecified matcher")
}

func createEventMatcher(eventMatcher *v1.EventMatcher, out *envoyroute.Route) error {
	eventType := eventMatcher.EventType
	if eventType == "" {
		return errors.New("must specify event_type")
	}

	out.Match.PathSpecifier = &envoyroute.RouteMatch_Path{
		Path: eventPath,
	}
	out.Match.Headers = append(out.Match.Headers, &envoyroute.HeaderMatcher{
		Name:  headerXEventType,
		Value: eventType,
	})
	return nil
}

func createRequestMatcher(requestMatcher *v1.RequestMatcher, out *envoyroute.Route) error {
	switch path := requestMatcher.Path.(type) {
	case *v1.RequestMatcher_PathRegex:
		out.Match.PathSpecifier = &envoyroute.RouteMatch_Regex{
			Regex: path.PathRegex,
		}
	case *v1.RequestMatcher_PathPrefix:
		out.Match.PathSpecifier = &envoyroute.RouteMatch_Prefix{
			Prefix: path.PathPrefix,
		}
	case *v1.RequestMatcher_PathExact:
		out.Match.PathSpecifier = &envoyroute.RouteMatch_Path{
			Path: path.PathExact,
		}
	}
	for headerName, headerValue := range requestMatcher.Headers {
		var regex bool
		if headerValue == "" {
			headerValue = ".*"
		}
		if strings.Contains(headerValue, ".*") {
			regex = true
		}
		out.Match.Headers = append(out.Match.Headers, &envoyroute.HeaderMatcher{
			Name:  headerName,
			Value: headerValue,
			Regex: &types.BoolValue{Value: regex},
		})
	}
	for paramName, paramValue := range requestMatcher.QueryParams {
		var regex bool
		if paramValue == "" {
			paramValue = ".*"
			regex = true
		}
		out.Match.QueryParameters = append(out.Match.QueryParameters, &envoyroute.QueryParameterMatcher{
			Name:  paramName,
			Value: paramValue,
			Regex: &types.BoolValue{Value: regex},
		})
	}
	if len(requestMatcher.Verbs) > 0 {
		out.Match.Headers = append(out.Match.Headers, &envoyroute.HeaderMatcher{
			Name:  ":method",
			Value: strings.Join(requestMatcher.Verbs, "|"),
			Regex: &types.BoolValue{Value: true},
		})
	}
	return nil
}
