package gfunc

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/protoutil"
)

type UpstreamSpec struct {
	Region    string `json:"region"`
	ProjectId string `json:"project_id"`
}

func DecodeUpstreamSpec(generic v1.UpstreamSpec) (*UpstreamSpec, error) {
	s := new(UpstreamSpec)
	if err := protoutil.UnmarshalStruct(generic, s); err != nil {
		return nil, err
	}
	return s, s.validateUpstream()
}

func (s *UpstreamSpec) validateUpstream() error {
	switch s.Region {

	case "northamerica-northeast1":
		return nil
	case "us-central1":
		return nil
	case "us-west1":
		return nil
	case "us-east4":
		return nil
	case "us-east1":
		return nil
	case "southamerica-east1":
		return nil
	case "europe-west1":
		return nil
	case "europe-west2":
		return nil
	case "europe-west3":
		return nil
	case "europe-west4":
		return nil
	case "asia-south1":
		return nil
	case "asia-southeast1":
		return nil
	case "asia-east1":
		return nil
	case "asia-northeast1":
		return nil
	case "australia-southeast1":
		return nil

	}
	return errors.New("no such region")
}

func (s *UpstreamSpec) GetGFuncHostname() string {
	return fmt.Sprintf("%s-%s.cloudfunctions.net", s.Region, s.ProjectId)
}

type FunctionSpec struct {
	URL string `json:"URL"`

	path string
	host string
}

func DecodeFunctionSpec(generic v1.FunctionSpec) (*FunctionSpec, error) {
	s := new(FunctionSpec)
	if err := protoutil.UnmarshalStruct(generic, s); err != nil {
		return nil, err
	}
	err := s.ValidateGFunc()
	if err != nil {
		return s, err
	}

	parsedUrl, err := url.Parse(s.URL)
	if err != nil {
		return s, err
	}

	s.path = parsedUrl.Path
	s.host = parsedUrl.Host
	return s, nil
}

func EncodeFunctionSpec(spec FunctionSpec) *types.Struct {
	v1Spec, err := protoutil.MarshalStruct(spec)
	if err != nil {
		panic(err)
	}
	return v1Spec
}

func (s *FunctionSpec) ValidateGFunc() error {
	if s.URL == "" {
		return errors.New("invalid function url")
	}
	return nil
}
