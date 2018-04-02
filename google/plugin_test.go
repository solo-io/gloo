package gfunc

import (
	"strings"
	"testing"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugin"
)

func TestProcessingNonGoogleUpstreamReturnsEmpty(t *testing.T) {
	upstreams := []*v1.Upstream{&v1.Upstream{}, &v1.Upstream{Type: "some-upstream"}}
	p := Plugin{}
	for _, u := range upstreams {
		err := p.ProcessUpstream(nil, u, nil)
		if err != nil {
			t.Errorf("Parsing non-Google upstream should not have errors %v", err)
		}
	}
}

func TestProcessingGoogleUpstream(t *testing.T) {
	upstream := &v1.Upstream{
		Type: UpstreamTypeGoogle,
		Spec: upstreamSpec("us-east1", "project-x"),
	}
	out := &envoyapi.Cluster{}
	p := Plugin{}
	err := p.ProcessUpstream(nil, upstream, out)
	if err != nil {
		t.Errorf("parsing valid upstream shouldn't have errors %v", err)
	}
	if !p.isNeeded {
		t.Errorf("google plugin should be required after valid google upstream")
	}
	if len(out.Hosts) < 1 {
		t.Errorf("was expecting hosts")
	}
	address, ok := out.Hosts[0].Address.(*envoycore.Address_SocketAddress)
	if !ok {
		t.Errorf("was expecting address")
	}
	if !strings.Contains(address.SocketAddress.Address, "us-east1") {
		t.Error("was expecting region in hostname")
	}
	if !strings.Contains(address.SocketAddress.Address, "project-x") {
		t.Error("was expecting project in hostname")
	}
}

func TestProcessingFunctionSpecForNonGoogle(t *testing.T) {
	p := Plugin{}
	nonGoogle := &plugin.FunctionPluginParams{}
	out, err := p.ParseFunctionSpec(nonGoogle, funcSpec("https://host.io/func"))
	if err != nil {
		t.Errorf("non google params should not result in error")
	}
	if out != nil {
		t.Errorf("non google params should return nil")
	}
}

func TestProcessingFunctionSpecForGoogle(t *testing.T) {
	p := Plugin{}
	nonGoogle := &plugin.FunctionPluginParams{UpstreamType: UpstreamTypeGoogle}
	out, err := p.ParseFunctionSpec(nonGoogle, funcSpec("https://host.io/func"))
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if host := get(out, functionHost); host != "host.io" {
		t.Error("host didn't match", host)
	}
	if path := get(out, functionPath); path != "/func" {
		t.Error("path didn't match", path)
	}
}

func get(s *types.Struct, key string) string {
	v, ok := s.Fields[key]
	if !ok {
		return ""
	}
	sv, ok := v.Kind.(*types.Value_StringValue)
	if !ok {
		return ""
	}
	return sv.StringValue
}
