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
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/common"
	"github.com/solo-io/gloo/pkg/plugin"
	"github.com/solo-io/gloo/pkg/protoutil"
	"k8s.io/apimachinery/pkg/util/runtime"
)

//go:generate protoc -I=. -I=${GOPATH}/src/github.com/gogo/protobuf/ --gogo_out=. transformation_filter.proto

const (
	filterName  = "io.solo.transformation"
	metadataKey = "transformation"
	pluginStage = plugin.PreOutAuth
)

func init() {
	plugin.Register(&Plugin{cachedTransformations: make(map[string]*Transformation)}, nil)
}

type Plugin struct {
	cachedTransformations map[string]*Transformation
}

func (p *Plugin) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func parseParameterToRegex(paramString string) ([]string, string) {
	// escape regex
	// TODO: make sure all envoy regex is being escaped here
	paramString = regexp.QuoteMeta(paramString)
	rxp := regexp.MustCompile("\\{([[:word:]]+)\\}")
	parameterNames := rxp.FindAllString(paramString, -1)
	for i, name := range parameterNames {
		parameterNames[i] = strings.TrimSuffix(strings.TrimPrefix(name, "{"), "}")
	}
	// remember that the order of the param names correlates with their order in the regex
	return parameterNames, rxp.ReplaceAllString(paramString, "([_[:alnum:]]+)")
}

func (p *Plugin) setHeaderFromParam(header, parameter string, extractors map[string]*Extraction, out *envoyroute.Route) {
	if parameter != "" {
		paramNames, pathRegexMatcher := parseParameterToRegex(parameter)
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
				Regex:    pathRegexMatcher,
				Subgroup: uint32(i + 1),
			}
			extractors[name] = extract
		}
	}
}

func (p *Plugin) ProcessRoute(_ *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	if in.Extensions == nil {
		return nil
	}
	spec, err := DecodeTransformationSpec(in.Extensions)
	if err != nil {
		return err
	}

	extractors := make(map[string]*Extraction)

	// special http2 headers
	p.setHeaderFromParam(":path", spec.Input.PathParameter, extractors, out)
	p.setHeaderFromParam(":method", spec.Input.MethodParameter, extractors, out)
	p.setHeaderFromParam(":scheme", spec.Input.MethodParameter, extractors, out)
	p.setHeaderFromParam(":authority", spec.Input.AuthorityParameter, extractors, out)
	// extra headers
	for headerName, headerValue := range spec.Input.HeaderParameters {
		p.setHeaderFromParam(headerName, headerValue, extractors, out)
	}

	// create templates
	// right now it's just a no-op, user writes inja directly
	headerTemplates := make(map[string]*InjaTemplate)
	for k, v := range spec.Output.HeaderTemplates {
		headerTemplates[k] = &InjaTemplate{Text: v}
	}

	t := &Transformation{
		Extractors: extractors,
		RequestTemplate: &RequestTemplate{
			Body: &InjaTemplate{
				Text: spec.Output.BodyTemplate,
			},
			Headers: headerTemplates,
		},
	}

	intHash, err := hashstructure.Hash(t, nil)
	if err != nil {
		return err
	}

	hash := fmt.Sprintf("%v", intHash)

	p.cachedTransformations[hash] = t

	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	filterMetadata := common.InitFilterMetadataField(filterName, metadataKey, out.Metadata)
	filterMetadata.Kind = &types.Value_StringValue{StringValue: hash}

	return nil
}

func (p *Plugin) HttpFilter(_ *plugin.FilterPluginParams) (*envoyhttp.HttpFilter, plugin.Stage) {
	filterConfig, err := protoutil.MarshalStruct(&Transformations{
		Transformations: p.cachedTransformations,
	})
	if err != nil {
		runtime.HandleError(err)
		return nil, 0
	}

	// clear cache
	p.cachedTransformations = make(map[string]*Transformation)

	return &envoyhttp.HttpFilter{
		Name:   filterName,
		Config: filterConfig,
	}, pluginStage
}
