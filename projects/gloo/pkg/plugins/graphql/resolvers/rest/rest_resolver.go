package rest

import (
	"strings"
	"time"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/types"

	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/dot_notation"
	resolver_utils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/utils"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	restResolverTypedExtensionConfigName = "io.solo.graphql.resolver.rest"

	// allowed setter keywords; for now just $body but may add $headers in the future
	BODY_SETTER = "body"
)

func TranslateRestResolver(upstreams types.UpstreamList, r *RESTResolver) (*v3.TypedExtensionConfig, error) {
	requestTransform, err := translateRequestTransform(r.Request)
	if err != nil {
		return nil, err
	}
	responseTransform, err := translateResponseTransform(r.Response)
	if err != nil {
		return nil, err
	}
	us, err := upstreams.Find(r.UpstreamRef.GetNamespace(), r.UpstreamRef.GetName())
	if err != nil {
		return nil, eris.Wrapf(err, "unable to find upstream `%s` in namespace `%s` to resolve schema", r.UpstreamRef.GetName(), r.UpstreamRef.GetNamespace())
	}
	restResolver := &v2.RESTResolver{
		ServerUri: &v3.HttpUri{
			Uri: "ignored", // ignored by graphql filter
			HttpUpstreamType: &v3.HttpUri_Cluster{
				Cluster: translator.UpstreamToClusterName(us.GetMetadata().Ref()),
			},
			Timeout: durationpb.New(1 * time.Second),
		},
		RequestTransform:      requestTransform,
		PreExecutionTransform: responseTransform,
		SpanName:              r.SpanName,
	}

	marshalledRestResolver, err := utils.MessageToAny(restResolver)
	if err != nil {
		return nil, eris.Wrapf(err, "unable to marshal restResolver")
	}

	return &v3.TypedExtensionConfig{
		Name:        restResolverTypedExtensionConfigName,
		TypedConfig: marshalledRestResolver,
	}, nil
}

func translateResponseTransform(transform *ResponseTemplate) (*v2.ResponseTemplate, error) {
	if transform == nil {
		return nil, nil
	}
	resultRoot, err := dot_notation.DotNotationToPathSegments(transform.ResultRoot)
	if err != nil {
		return nil, eris.Wrapf(err, "error translating result root path %s", transform.ResultRoot)
	}
	setters := map[string]*v2.TemplatedPath{}
	for k, path := range transform.Setters {
		templatedPath, err := TranslateSetter(path)
		if err != nil {
			return nil, eris.Wrapf(err, "error, unable to translate setter string to templated path")
		}
		setters[k] = templatedPath
	}
	result := &v2.ResponseTemplate{
		ResultRoot: resultRoot,
		Setters:    setters,
	}
	return result, nil
}

func translateRequestTransform(transform *RequestTemplate) (*v2.RequestTemplate, error) {
	if transform == nil {
		return nil, nil
	}
	headersMap, err := resolver_utils.TranslateStringValueProviderMap(transform.GetHeaders())
	if err != nil {
		return nil, err
	}
	queryParamsMap, err := resolver_utils.TranslateStringValueProviderMap(transform.GetQueryParams())
	if err != nil {
		return nil, err
	}
	jv, err := resolver_utils.TranslateJsonValue(transform.GetBody())
	if err != nil {
		return nil, err
	}
	rt := &v2.RequestTemplate{
		Headers:      headersMap,
		QueryParams:  queryParamsMap,
		OutgoingBody: jv, // filled in later
	}

	return rt, nil
}

func TranslateSetter(templatedPathString string) (*v2.TemplatedPath, error) {
	ret := &v2.TemplatedPath{
		PathTemplate: templatedPathString,
		NamedPaths:   map[string]*v2.Path{},
	}
	subMatches := resolver_utils.ProviderTemplateRegex.FindAllStringSubmatch(templatedPathString, -1)
	if len(subMatches) < 1 {
		return nil, eris.Errorf("malformed setter %s, needs to match template string regex %s", templatedPathString, resolver_utils.ProviderTemplateRegex.String())
	}

	for _, matches := range subMatches {
		match := matches[1]
		pathSegment, err := dot_notation.DotNotationToPathSegments(match)
		if err != nil {
			return nil, eris.Wrapf(err, "unable to parse graphql extraction string %s", templatedPathString)
		}
		if len(pathSegment) <= 1 {
			return nil, eris.New("invalid extraction string " + templatedPathString)
		}

		if pathSegment[0].GetKey() != BODY_SETTER {
			return nil, eris.New("currently only support for grabbing from the request body with {$body.}")
		}
		// remove $body from path segments, this was a special keyword
		// we didn't remove it from `match` (NamedPaths` key) so we could leave room in our api for new keywords, e.g. $headers
		pathSegment = pathSegment[1:]

		tp := &v2.Path{
			Segments: pathSegment,
		}
		matchSanitized := strings.ReplaceAll(match, ".", "")
		ret.NamedPaths[matchSanitized] = tp
		newMatchTemplate := "{" + matchSanitized + "}"
		ret.PathTemplate = strings.ReplaceAll(ret.PathTemplate, matches[0], newMatchTemplate)
		// if the match template is the only thing that exists in the path template,
		// then the user does not want an interpolated string, and wants the original value instead,
		// so if we do not pass a path template, envoy will know to keep the original value rather
		// than try and coerce it to a string for an interpolated template
		if len(subMatches) == 1 && len(ret.PathTemplate) == len(newMatchTemplate) {
			ret.PathTemplate = ""
		}
	}
	return ret, nil
}
