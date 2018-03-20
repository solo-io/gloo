package grpc

import (
	"strings"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"github.com/gogo/googleapis/google/api"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/types"

	"crypto/sha1"

	"fmt"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/common/annotations"
	"github.com/solo-io/gloo-plugins/transformation"
	"github.com/solo-io/gloo/pkg/coreplugins/common"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/plugin"
)

type Plugin struct {
	transformation *transformation.Plugin
}

const (
	filterName  = "envoy.grpc_json_transcoder"
	pluginStage = plugin.PostInAuth

	ServiceTypeGRPC = "gRPC-Functions"
)

/*

NOTES; allow auto-generation of routes that match the http rules in the descriptor
that way if someone wants to reuse those routes, let them just work. note: for this case,
we DO need to allow recalculation of the route. hm.

 - need to create envoy routes for the grpc service. this is how we map the functions to the cluster
   like so:
              - match: { prefix: "/bookstore.Bookstore" }
                route: { cluster: service_google }

          http_filters:
          - name: envoy.grpc_json_transcoder
            config:
              proto_descriptor: ./proto.pb
              services: [bookstore.Bookstore]
          - name: envoy.router

for every service in the file, create a d

*/

func (p *Plugin) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func isOurs(in *v1.Upstream) bool {
	if in.Metadata == nil || in.Metadata.Annotations == nil {
		return false
	}
	return in.Metadata.Annotations[annotations.ServiceType] == ServiceTypeGRPC
}

func (p *Plugin) ProcessUpstream(params *plugin.UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error {
	if !isOurs(in) {
		return nil
	}

	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	common.InitFilterMetadata(filterName, out.Metadata)
	out.Metadata.FilterMetadata[filterName] = &types.Struct{
		Fields: make(map[string]*types.Value),
	}

	return nil
}

func ConvertProto(b []byte) (*descriptor.FileDescriptorSet, error) {
	var fileDescriptor descriptor.FileDescriptorSet
	err := proto.Unmarshal(b, &fileDescriptor)
	return &fileDescriptor, err
}

func AddHttpRulesToProto(upstreamName, serviceName string, set *descriptor.FileDescriptorSet) {
	for _, file := range set.File {
		for _, svc := range file.Service {
			if svc.Name == nil || *svc.Name != serviceName {
				continue
			}
			for _, method := range svc.Method {
				proto.SetExtension(method.Options, api.E_Http, &api.HttpRule{
					Pattern: &api.HttpRule_Post{
						Post: generateHttpPath(upstreamName, serviceName, *method.Name),
					},
				})
			}
		}
	}
}

func generateHttpPath(upstreamName, serviceName, methodName string) string {
	h := sha1.New()
	h.Write([]byte(upstreamName + serviceName))
	return "/" + fmt.Sprintf("%x", h.Sum(nil))[:8] + "/" + upstreamName + "/" + serviceName + "/" + methodName
}

func FuncsForProto(serviceName string, set *descriptor.FileDescriptorSet) []*v1.Function {
	var funcs []*v1.Function
	for _, file := range set.File {
		for _, svc := range file.Service {
			if svc.Name == nil || *svc.Name != serviceName {
				continue
			}
			for _, method := range svc.Method {
				g, err := proto.GetExtension(method.Options, api.E_Http)
				if err != nil {
					log.Printf("missing http option on the extensions, skipping: %v", *method.Name)
					continue
				}
				httpRule, ok := g.(*api.HttpRule)
				if !ok {
					panic(g)
				}
				log.Printf("rule: %v", httpRule)
				verb, path := verbAndPathForRule(httpRule)
				fn := &v1.Function{
					Name: *method.Name,
					Spec: transformation.EncodeFunctionSpec(transformation.Template{
						Path:            toInjaTemplateFormat(path),
						Header:          map[string]string{":method": verb},
						PassthroughBody: true,
					}),
				}
				funcs = append(funcs, fn)
			}
		}
		log.Printf("%v", file.MessageType)
	}
	return funcs
}

func toInjaTemplateFormat(in string) string {
	in = strings.Replace(in, "{", "{{", -1)
	return strings.Replace(in, "}", "}}", -1)
}

func verbAndPathForRule(httpRule *api.HttpRule) (string, string) {
	switch rule := httpRule.Pattern.(type) {
	case *api.HttpRule_Get:
		return "GET", rule.Get
	case *api.HttpRule_Custom:
		return rule.Custom.Kind, rule.Custom.Path
	case *api.HttpRule_Delete:
		return "DELETE", rule.Delete
	case *api.HttpRule_Patch:
		return "PATCH", rule.Patch
	case *api.HttpRule_Post:
		return "POST", rule.Post
	case *api.HttpRule_Put:
		return "PUT", rule.Put
	}
	panic("unknown rule type")
}

func lookupMessageType(inputType string, messageTypes []*descriptor.DescriptorProto) *descriptor.DescriptorProto {
	for _, msg := range messageTypes {
		if *msg.Name == inputType {
			return msg
		}
	}
	return nil
}
