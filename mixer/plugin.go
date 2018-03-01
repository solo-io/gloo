package gfunc

import (
	"bytes"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/types"

	"text/template"

	"github.com/gogo/protobuf/jsonpb"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/plugin"
)

func init() {
	plugin.Register(&Plugin{}, nil)
}

type MixerConfigTemplateValues struct {
	PodName, Namespace, ServiceName, MixerClusterName string
}

type Plugin struct {
	template MixerConfigTemplateValues
}

const (
	ServiceTypeNatsStreaming = "nats-streaming"

	// generic plugin info
	filterName          = "mixer"
	pluginStage         = plugin.InAuth
	mixerTemplateConfig = `
	{
		"mixer_attributes": {
		 "destination.ip": "",
		 "destination.service": "{{.ServiceName}}.{{.Namespace}}.svc.cluster.local",
		 "destination.uid": "kubernetes://{{.PodName}}.{{.Namespace}}"
		},
		"forward_attributes": {
		 "source.ip": "",
		 "source.uid": "kubernetes://{{.PodName}}.{{.Namespace}}"
		},
		"quota_name": "RequestCount",
		"v2": {
		 "transport" : {
			 "check_cluster" : "{{.MixerClusterName}}",
			 "report_cluster" : "{{.MixerClusterName}}"
		 },
		 "defaultDestinationService": "{{.ServiceName}}.{{.Namespace}}.svc.cluster.local",
		 "forwardAttributes": {
		  "attributes": {
		   "source.ip": {
			"bytesValue": null
		   },
		   "source.uid": {
			"stringValue": "kubernetes://{{.PodName}}.{{.Namespace}}"
		   }
		  }
		 },
		 "mixerAttributes": {
		  "attributes": {
		   "destination.ip": {
			"bytesValue": null
		   },
		   "destination.uid": {
			"stringValue": "kubernetes://{{.PodName}}.{{.Namespace}}"
		   }
		  }
		 },
		 "serviceConfigs": {
		  "{{.ServiceName}}.{{.Namespace}}.svc.cluster.local": {
		   "mixerAttributes": {
			"attributes": {
			 "destination.service": {
			  "stringValue": "{{.ServiceName}}.{{.Namespace}}.svc.cluster.local"
			 }
			}
		   }
		  }
		 }
		}
	   }
	
	`
)

func (p *Plugin) GetDependencies(cfg *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *Plugin) ProcessUpstream(params *plugin.UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error {

	/*
		TODO: find the mixer upstream;
		1. if we have the mixer upstream, save its cluster name
		2. if we have the gloo ingress upstream, get the gloo service name from the annotation of the gloo ingress upstream
		3. on filter creation, create a mixer filter from template with above values. pray.
	*/

	return nil
}

func anyempty(strs ...string) bool {
	for _, a := range strs {
		if len(a) == 0 {
			return true
		}
	}
	return false
}

func (p *Plugin) HttpFilters(params *plugin.FilterPluginParams) []plugin.StagedFilter {
	if anyempty(p.template.MixerClusterName, p.template.Namespace, p.template.PodName, p.template.ServiceName) {
		return nil
	}
	tmpl := template.Must(template.New("mixer").Parse(mixerTemplateConfig))
	var b bytes.Buffer
	tmpl.Execute(&b, p.template)
	jsonstring := b.String()

	var pbconfig types.Struct
	err := jsonpb.UnmarshalString(jsonstring, &pbconfig)
	if err != nil {
		log.Warnf("you found a bug: %v", err)
		return nil
	}
	filter := &envoyhttp.HttpFilter{
		Name:   filterName,
		Config: &pbconfig,
	}
	return []plugin.StagedFilter{{HttpFilter: filter, Stage: pluginStage}}
}
