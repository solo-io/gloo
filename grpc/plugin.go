package grpc

import (
	"crypto/sha1"
	"fmt"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoytranscoder "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/transcoder/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/googleapis/google/api"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/common/transformation"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/plugin"
	"github.com/solo-io/gloo/pkg/protoutil"
)

func init() {
	plugin.Register(NewPlugin(), nil)
}

func NewPlugin() *Plugin {
	return &Plugin{
		serviceDescriptors: make(map[string]*descriptor.FileDescriptorSet),
		upstreamServices:   make(map[string]string),
		transformation:     transformation.NewTransformationPlugin(),
	}
}

type Plugin struct {
	// map service names to their descriptors
	serviceDescriptors map[string]*descriptor.FileDescriptorSet
	// keep track of which service belongs to which upstream
	upstreamServices map[string]string
	transformation   transformation.Plugin
}

const (
	filterName  = "envoy.grpc_json_transcoder"
	pluginStage = plugin.PreOutAuth

	ServiceTypeGRPC = "gRPC"
)

func (p *Plugin) GetDependencies(cfg *v1.Config) *plugin.Dependencies {
	deps := &plugin.Dependencies{}
	for _, us := range cfg.Upstreams {
		if !isOurs(us) {
			continue
		}
		serviceSpec, err := DecodeServiceProperties(us.ServiceInfo.Properties)
		if err != nil {
			log.Warnf("%v: error parsing service properties for upstream %v: %v",
				ServiceTypeGRPC, us.Name, err)
			continue
		}
		deps.FileRefs = append(deps.FileRefs, serviceSpec.DescriptorsFileRef)
	}
	return deps
}

func isOurs(in *v1.Upstream) bool {
	return in.ServiceInfo != nil && in.ServiceInfo.Type == ServiceTypeGRPC
}

func (p *Plugin) ProcessUpstream(params *plugin.UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error {
	if !isOurs(in) {
		return nil
	}

	serviceProperties, err := DecodeServiceProperties(in.ServiceInfo.Properties)
	if err != nil {
		return errors.Wrap(err, "parsing service properties")
	}
	fileRef := serviceProperties.DescriptorsFileRef
	serviceName := serviceProperties.GRPCServiceName

	if fileRef == "" {
		return errors.New("service_info.properties.descriptors_file_ref cannot be empty")
	}
	if serviceName == "" {
		return errors.New("service_info.properties.service_name cannot be empty")
	}
	log.Printf("files: %v", params.Files)
	descriptorsFile, ok := params.Files[fileRef]
	if !ok {
		return errors.Errorf("descriptors file not found for file ref %v", fileRef)
	}
	descriptors, err := convertProto(descriptorsFile)
	if err != nil {
		return errors.Wrapf(err, "parsing file %v as a proto descriptor set", fileRef)
	}

	if err := addHttpRulesToProto(in.Name, serviceName, descriptors); err != nil {
		return errors.Wrap(err, "failed to generate http rules for proto descriptors")
	}

	p.transformation.ActivateFilterForCluster(out)

	// cache the descriptors; we'll need then when we create our grpc filters
	p.serviceDescriptors[serviceName] = descriptors
	// keep track of which service belongs to which upstream
	p.upstreamServices[in.Name] = serviceName

	return nil
}

func convertProto(b []byte) (*descriptor.FileDescriptorSet, error) {
	var fileDescriptor descriptor.FileDescriptorSet
	err := proto.Unmarshal(b, &fileDescriptor)
	return &fileDescriptor, err
}

func getPath(matcher *v1.RequestMatcher) string {
	switch path := matcher.Path.(type) {
	case *v1.RequestMatcher_PathPrefix:
		return path.PathPrefix
	case *v1.RequestMatcher_PathExact:
		return path.PathExact
	case *v1.RequestMatcher_PathRegex:
		return path.PathRegex
	}
	panic("invalid matcher")
}

func (p *Plugin) ProcessRoute(_ *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	if in.Extensions == nil {
		matcher, ok := in.Matcher.(*v1.Route_RequestMatcher)
		if ok {
			in.Extensions = transformation.EncodeRouteExtension(transformation.RouteExtension{
				Parameters: &transformation.Parameters{
					Path: getPath(matcher.RequestMatcher) + "?{query_string}",
				},
			})
		}
	}
	return p.transformation.AddRequestTransformationsToRoute(p.templateForFunction, in, out)
}

func (p *Plugin) templateForFunction(dest *v1.Destination_Function) (*transformation.TransformationTemplate, error) {
	upstreamName := dest.Function.UpstreamName
	serviceName, ok := p.upstreamServices[upstreamName]
	if !ok {
		// the upstream is not a grpc desintation
		return nil, nil
	}

	// method name should be function name in this case. TODO: document in the api
	methodName := dest.Function.FunctionName

	// create the transformation for the route

	outPath := httpPath(upstreamName, serviceName, methodName)

	// add query matcher to out path. kombina for now
	// TODO: support query for matching
	outPath += `?{{ default(query_string), "")}}`

	// we always choose post
	httpMethod := "POST"
	return &transformation.TransformationTemplate{
		Headers: map[string]*transformation.InjaTemplate{
			":method": {Text: httpMethod},
			":path":   {Text: outPath},
		},
		BodyTransformation: &transformation.TransformationTemplate_MergeExtractorsToBody{
			MergeExtractorsToBody: &transformation.MergeExtractorsToBody{},
		},
	}, nil
}

func addHttpRulesToProto(upstreamName, serviceName string, set *descriptor.FileDescriptorSet) error {
	var serviceFound, googleApiHttpFound, googleApiAnnotationsFound bool
	for _, file := range set.File {
		log.Printf("%v", *file.Name)
		if *file.Name == "google/api/http.proto" {
			googleApiHttpFound = true
			continue
		}
		if *file.Name == "google/api/annotations.proto" {
			googleApiAnnotationsFound = true
			continue
		}
	findService:
		for _, svc := range file.Service {
			if *svc.Name == serviceName {
				for _, method := range svc.Method {
					if err := proto.SetExtension(method.Options, api.E_Http, &api.HttpRule{
						Pattern: &api.HttpRule_Post{
							Post: httpPath(upstreamName, serviceName, *method.Name),
						},
						Body: "*",
					}); err != nil {
						return errors.Wrap(err, "setting http extensions for method.Options")
					}
					serviceFound = true
					break findService
				}
				return nil
			}
		}

	}

	if !googleApiHttpFound {
		addGoogleApisHttp(set)
	}

	if !googleApiAnnotationsFound {
		addGoogleApisAnnotations(set)
	}

	if !serviceFound {
		return errors.Errorf("could not find match: %v/%v", upstreamName, serviceName)
	}
	return nil
}

func httpPath(upstreamName, serviceName, methodName string) string {
	h := sha1.New()
	h.Write([]byte(upstreamName + serviceName))
	return "/" + fmt.Sprintf("%x", h.Sum(nil))[:8] + "/" + upstreamName + "/" + serviceName + "/" + methodName
}

func (p *Plugin) HttpFilters(params *plugin.FilterPluginParams) []plugin.StagedFilter {
	if len(p.serviceDescriptors) == 0 {
		return nil
	}

	transformationFilter := p.transformation.GetTransformationFilter()
	if transformationFilter == nil {
		log.Warnf("ERROR: nil transformation filter returned from transformation plugin")
		return nil
	}

	filters := []plugin.StagedFilter{*transformationFilter}

	for serviceName, protoDescriptor := range p.serviceDescriptors {
		descriptorBytes, err := proto.Marshal(protoDescriptor)
		if err != nil {
			log.Warnf("ERROR: marshaling proto descriptor: %v", err)
			continue
		}
		filterConfig, err := protoutil.MarshalStruct(&envoytranscoder.GrpcJsonTranscoder{
			DescriptorSet: &envoytranscoder.GrpcJsonTranscoder_ProtoDescriptorBin{
				ProtoDescriptorBin: descriptorBytes,
			},
			Services:               []string{serviceName},
			SkipRecalculatingRoute: true,
		})
		if err != nil {
			log.Warnf("ERROR: marshaling GrpcJsonTranscoder config: %v", err)
			return nil
		}
		filters = append(filters, plugin.StagedFilter{
			HttpFilter: &envoyhttp.HttpFilter{
				Name:   filterName,
				Config: filterConfig,
			},
			Stage: pluginStage,
		})
	}

	// clear cache
	p.serviceDescriptors = make(map[string]*descriptor.FileDescriptorSet)
	p.upstreamServices = make(map[string]string)

	return filters
}

//func FuncsForProto(serviceName string, set *descriptor.FileDescriptorSet) []*v1.Function {
//	var funcs []*v1.Function
//	for _, file := range set.File {
//		for _, svc := range file.Service {
//			if svc.Name == nil || *svc.Name != serviceName {
//				continue
//			}
//			for _, method := range svc.Method {
//				g, err := proto.GetExtension(method.Options, api.E_Http)
//				if err != nil {
//					log.Printf("missing http option on the extensions, skipping: %v", *method.Name)
//					continue
//				}
//				httpRule, ok := g.(*api.HttpRule)
//				if !ok {
//					panic(g)
//				}
//				log.Printf("rule: %v", httpRule)
//				verb, path := verbAndPathForRule(httpRule)
//				fn := &v1.Function{
//					Name: *method.Name,
//					Spec: transformation.EncodeFunctionSpec(transformation.Template{
//						Path:            toInjaTemplateFormat(path),
//						Header:          map[string]string{":method": verb},
//						PassthroughBody: true,
//					}),
//				}
//				funcs = append(funcs, fn)
//			}
//		}
//		log.Printf("%v", file.MessageType)
//	}
//	return funcs
//}
//
//func toInjaTemplateFormat(in string) string {
//	in = strings.Replace(in, "{", "{{", -1)
//	return strings.Replace(in, "}", "}}", -1)
//}
//
//func verbAndPathForRule(httpRule *api.HttpRule) (string, string) {
//	switch rule := httpRule.Pattern.(type) {
//	case *api.HttpRule_Get:
//		return "GET", rule.Get
//	case *api.HttpRule_Custom:
//		return rule.Custom.Kind, rule.Custom.Path
//	case *api.HttpRule_Delete:
//		return "DELETE", rule.Delete
//	case *api.HttpRule_Patch:
//		return "PATCH", rule.Patch
//	case *api.HttpRule_Post:
//		return "POST", rule.Post
//	case *api.HttpRule_Put:
//		return "PUT", rule.Put
//	}
//	panic("unknown rule type")
//}
//
//func lookupMessageType(inputType string, messageTypes []*descriptor.DescriptorProto) *descriptor.DescriptorProto {
//	for _, msg := range messageTypes {
//		if *msg.Name == inputType {
//			return msg
//		}
//	}
//	return nil
//}
//
