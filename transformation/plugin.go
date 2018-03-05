package transformation

import (
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"

	"fmt"
	"regexp"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/coreplugins/common"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/plugin"
	"github.com/solo-io/gloo/pkg/protoutil"
)

//go:generate protoc -I=. -I=${GOPATH}/src/github.com/gogo/protobuf/ --gogo_out=. transformation_filter.proto

const (
	filterName          = "io.solo.transformation"
	metadataRequestKey  = "request-transformation"
	metadataResponseKey = "response-transformation"
	pluginStage         = plugin.PreOutAuth
)

func init() {
	plugin.Register(&Plugin{CachedTransformations: make(map[string]*Transformation)}, nil)
}

type Plugin struct {
	CachedTransformations map[string]*Transformation
}

func (p *Plugin) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *Plugin) ProcessRoute(pluginParams *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	if err := p.processRequestTransformationsForRoute(pluginParams, in, out); err != nil {
		return errors.Wrap(err, "failed to process request transformation")
	}
	if err := p.processResponseTransformationsForRoute(pluginParams, in, out); err != nil {
		return errors.Wrap(err, "failed to process response transformation")
	}
	return nil
}

func (p *Plugin) processRequestTransformationsForRoute(pluginParams *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	if in.Extensions == nil {
		return nil
	}

	extension, err := DecodeRouteExtension(in.Extensions)
	if err != nil {
		return err
	}

	if extension.Parameters == nil {
		return nil
	}

	extractors := make(map[string]*Extraction)

	// special http2 headers, get the whole thing for free
	// as a convenience to the user
	// TODO: add more
	for _, header := range []string{
		"path",
		"method",
		"scheme",
		"authority",
	} {
		addHeaderExtractorFromParam(":"+header, "{"+header+"}", extractors)
	}
	// headers we support submatching on
	// custom as well as the path and authority/host header
	addHeaderExtractorFromParam(":path", extension.Parameters.Path, extractors)
	addHeaderExtractorFromParam(":authority", extension.Parameters.Authority, extractors)
	for headerName, headerValue := range extension.Parameters.Headers {
		addHeaderExtractorFromParam(headerName, headerValue, extractors)
	}

	// calculate the templates for all these transformations
	if err := p.setTransformationsForRoute(pluginParams.Upstreams, in, extractors, out); err != nil {
		return errors.Wrap(err, "resolving request transformations for route")
	}

	return nil
}

// TODO: clean up the response transformation
// params should live on the source (upstream/function)
func (p *Plugin) processResponseTransformationsForRoute(pluginParams *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	if in.Extensions == nil {
		return nil
	}

	extension, err := DecodeRouteExtension(in.Extensions)
	if err != nil {
		return err
	}

	if extension.ResponseTemplate == nil {
		return nil
	}

	extractors := make(map[string]*Extraction)

	if extension.ResponseParams != nil {
		for headerName, headerValue := range extension.ResponseParams.Headers {
			addHeaderExtractorFromParam(headerName, headerValue, extractors)
		}
	}

	// calculate the templates for all these transformations
	if err := p.setResponseTransformationForRoute(*extension.ResponseTemplate, extractors, out); err != nil {
		return errors.Wrap(err, "resolving request transformations for route")
	}

	return nil
}

func addHeaderExtractorFromParam(header, parameter string, extractors map[string]*Extraction) {
	if parameter == "" {
		return
	}
	// remember that the order of the param names correlates with their order in the regex
	paramNames, regexMatcher := getNamesAndRegexFromParamString(parameter)
	log.Debugf("transformation pluginN: extraction for header %v: parameters: %v regex matcher: %v", header, paramNames, regexMatcher)
	// if no regex, this is a "default variable" that the user gets for free
	if len(paramNames) == 0 {
		// extract everything
		// TODO(yuval): create a special extractor that doesn't use regex when we just want the whole thing
		extract := &Extraction{
			Header:   header,
			Regex:    "(.*)",
			Subgroup: uint32(1),
		}
		extractors[strings.TrimPrefix(header, ":")] = extract
	}

	// otherwise it's regex, and we need to create an extraction for each variable name they defined
	for i, name := range paramNames {
		extract := &Extraction{
			Header:   header,
			Regex:    regexMatcher,
			Subgroup: uint32(i + 1),
		}
		extractors[name] = extract
	}
}

func getNamesAndRegexFromParamString(paramString string) ([]string, string) {
	// escape regex
	// TODO: make sure all envoy regex is being escaped here
	rxp := regexp.MustCompile("\\{([[:word:]]+)\\}")
	parameterNames := rxp.FindAllString(paramString, -1)
	for i, name := range parameterNames {
		parameterNames[i] = strings.TrimSuffix(strings.TrimPrefix(name, "{"), "}")
	}

	return parameterNames, buildRegexString(rxp, paramString)
}

func buildRegexString(rxp *regexp.Regexp, paramString string) string {
	var regexString string
	var prevEnd int
	for _, startStop := range rxp.FindAllStringIndex(paramString, -1) {
		start := startStop[0]
		end := startStop[1]
		subStr := regexp.QuoteMeta(paramString[prevEnd:start]) + "([_[:alnum:]]+)"
		regexString += subStr
		prevEnd = end
	}

	return regexString + regexp.QuoteMeta(paramString[prevEnd:])
}

// gets all transformations a route may need
// if single destination, just one transformation
// if multi destination, one transformation for each functional
// that specifies a transformation spec
func (p *Plugin) setTransformationsForRoute(upstreams []*v1.Upstream, in *v1.Route, extractors map[string]*Extraction, out *envoyroute.Route) error {
	switch {
	case in.MultipleDestinations != nil:
		for _, dest := range in.MultipleDestinations {
			err := p.setTransformationForFunction(upstreams, dest.Destination, extractors, out)
			if err != nil {
				return errors.Wrap(err, "getting transformation for function")
			}
		}
	case in.SingleDestination != nil:
		err := p.setTransformationForFunction(upstreams, in.SingleDestination, extractors, out)
		if err != nil {
			return errors.Wrap(err, "getting transformation for function")
		}
	}
	return nil
}

func (p *Plugin) setTransformationForFunction(upstreams []*v1.Upstream, dest *v1.Destination, extractors map[string]*Extraction, out *envoyroute.Route) error {
	hash, transformation, err := getTransformationForFunction(upstreams, dest, extractors)
	if err != nil {
		return errors.Wrap(err, "getting transformation for function")
	}
	// no transformations for this destination
	if transformation == nil {
		return nil
	}

	// cache the transformation, the filter config needs to contain all of them
	p.CachedTransformations[hash] = transformation

	// set the filter metadata on the route
	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	filterMetadata := common.InitFilterMetadataField(filterName, metadataRequestKey, out.Metadata)
	filterMetadata.Kind = &types.Value_StringValue{StringValue: hash}

	return nil
}

func getTransformationForFunction(upstreams []*v1.Upstream, dest *v1.Destination, extractors map[string]*Extraction) (string, *Transformation, error) {
	fnDestination, ok := dest.DestinationType.(*v1.Destination_Function)
	if !ok {
		// not a functional route, nothing to do
		return "", nil, nil
	}
	fn, err := findFunction(upstreams, fnDestination.Function.UpstreamName, fnDestination.Function.FunctionName)
	if err != nil {
		return "", nil, errors.Wrap(err, "finding function")
	}
	outputTemplates, err := DecodeFunctionSpec(fn.Spec)
	if err != nil {
		return "", nil, errors.Wrap(err, "decoding function spec")
	}

	// create templates
	// right now it's just a no-op, user writes inja directly
	headerTemplates := make(map[string]*InjaTemplate)
	for k, v := range outputTemplates.Header {
		headerTemplates[k] = &InjaTemplate{Text: v}
	}

	if outputTemplates.Path != "" {
		headerTemplates[":path"] = &InjaTemplate{Text: outputTemplates.Path}
	}

	t := Transformation{
		Extractors: extractors,
		RequestTemplate: &RequestTemplate{
			Body: &InjaTemplate{
				Text: outputTemplates.Body,
			},
			Headers: headerTemplates,
		},
	}

	intHash, err := hashstructure.Hash(t, nil)
	if err != nil {
		return "", nil, err
	}

	hash := fmt.Sprintf("%v", intHash)

	return hash, &t, nil
}

func findFunction(upstreams []*v1.Upstream, upstreamName, functionName string) (*v1.Function, error) {
	for _, us := range upstreams {
		if us.Name == upstreamName {
			for _, fn := range us.Functions {
				if fn.Name == functionName {
					return fn, nil
				}
			}
		}
	}
	return nil, errors.Errorf("function %v/%v not found", upstreamName, functionName)
}

func (p *Plugin) setResponseTransformationForRoute(template Template, extractors map[string]*Extraction, out *envoyroute.Route) error {
	// create templates
	// right now it's just a no-op, user writes inja directly
	headerTemplates := make(map[string]*InjaTemplate)
	for k, v := range template.Header {
		headerTemplates[k] = &InjaTemplate{Text: v}
	}

	transformation := Transformation{
		Extractors: extractors,
		RequestTemplate: &RequestTemplate{
			Body: &InjaTemplate{
				Text: template.Body,
			},
			Headers: headerTemplates,
		},
	}

	intHash, err := hashstructure.Hash(transformation, nil)
	if err != nil {
		return errors.Wrap(err, "generating hash")
	}

	hash := fmt.Sprintf("%v", intHash)

	// cache the transformation, the filter config needs to contain all of them
	p.CachedTransformations[hash] = &transformation

	// set the filter metadata on the route
	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	filterMetadata := common.InitFilterMetadataField(filterName, metadataResponseKey, out.Metadata)
	filterMetadata.Kind = &types.Value_StringValue{StringValue: hash}

	return nil
}
func (p *Plugin) HttpFilters(params *plugin.FilterPluginParams) []plugin.StagedFilter {
	filterConfig, err := protoutil.MarshalStruct(&Transformations{
		Transformations: p.CachedTransformations,
	})
	if err != nil {
		runtime.HandleError(err)
		return nil
	}

	// clear cache
	p.CachedTransformations = make(map[string]*Transformation)

	return []plugin.StagedFilter{{HttpFilter: &envoyhttp.HttpFilter{
		Name:   filterName,
		Config: filterConfig,
	}, Stage: pluginStage}}
}
