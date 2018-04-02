package gfunc

import (
	"strings"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
)

func funcSpec(u string) v1.FunctionSpec {
	return &types.Struct{
		Fields: map[string]*types.Value{
			"URL": &types.Value{
				Kind: &types.Value_StringValue{StringValue: u},
			},
		},
	}
}
func TestDecodeFuncSpecWithInvalidSpecFails(t *testing.T) {
	data := []*types.Struct{
		&types.Struct{
			Fields: map[string]*types.Value{},
		},
		funcSpec("apple"),
		funcSpec("http://apple.com"),
		funcSpec("http://solo.io:8433"),
		funcSpec("nowhere]:/apple"),
	}
	for _, d := range data {
		_, err := DecodeFunctionSpec(d)
		if err == nil {
			t.Errorf("invalid function spec should have returned error: %v", d)
		}
	}
}

type row struct {
	spec v1.FunctionSpec
	host string
	path string
}

func TestDecodeFuncSpecWithValidURL(t *testing.T) {
	data := []row{
		row{spec: funcSpec("http://test.com/apple"), host: "test.com", path: "/apple"},
		row{spec: funcSpec("http://solo.io/"), host: "solo.io", path: "/"},
	}
	for _, d := range data {
		f, err := DecodeFunctionSpec(d.spec)
		if err != nil {
			t.Errorf("error creating function spec %v", err)
		}
		if d.host != f.host {
			t.Errorf("functionspec created with wrong host. expected %s got %s", d.host, f.host)
		}
		if d.path != f.path {
			t.Errorf("functionspec created with wrong path. expected %s got %s", d.path, f.path)
		}
	}
}

func upstreamSpec(region, projectID string) v1.UpstreamSpec {
	return &types.Struct{
		Fields: map[string]*types.Value{
			"region": &types.Value{
				Kind: &types.Value_StringValue{StringValue: region},
			},
			"project_id": &types.Value{
				Kind: &types.Value_StringValue{StringValue: projectID},
			},
		},
	}
}

func TestDecodeUpstreamWithInvalidRegion(t *testing.T) {
	data := []*types.Struct{
		&types.Struct{
			Fields: map[string]*types.Value{},
		},
		upstreamSpec("solo", "gloo"),
		upstreamSpec("", "x23"),
	}
	for _, d := range data {
		_, err := DecodeUpstreamSpec(d)
		if err == nil {
			t.Errorf("invalid upstream spec should have returned error: %v", d)
		}
	}
}

func TestDecodeUpstreamWithInvalidProject(t *testing.T) {
	data := []*types.Struct{
		&types.Struct{
			Fields: map[string]*types.Value{},
		},
		upstreamSpec("us-east1", ""),
	}
	for _, d := range data {
		_, err := DecodeUpstreamSpec(d)
		if err == nil {
			t.Errorf("invalid upstream spec should have returned error: %v", d)
		}
	}
}

type upstreamTestCase struct {
	spec      v1.UpstreamSpec
	region    string
	projectID string
}

func TestDecodeUpstreamWithValidData(t *testing.T) {
	data := []upstreamTestCase{
		upstreamTestCase{spec: upstreamSpec("us-east1", "project-231x"), region: "us-east1", projectID: "project-231x"},
	}
	for _, d := range data {
		u, err := DecodeUpstreamSpec(d.spec)
		if err != nil {
			t.Errorf("error decoding upstream %v: %v", d.spec, err)
		}
		if d.region != u.Region {
			t.Errorf("decoding upstream failed. expected %s got %s", d.region, u.Region)
		}
		if d.projectID != u.ProjectId {
			t.Errorf("decoding upstream failed. expected %s got %s", d.projectID, u.ProjectId)
		}

		hostname := u.GetGFuncHostname()
		if !strings.Contains(hostname, d.region) {
			t.Errorf("hostname doesn't contain region")
		}
		if !strings.Contains(hostname, d.projectID) {
			t.Errorf("hostname doesn't contain project ID")
		}
	}
}
