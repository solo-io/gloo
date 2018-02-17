package aws

import (
	"errors"
	"fmt"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/protoutil"
)

type UpstreamSpec struct {
	Region    string `json:"region"`
	SecretRef string `json:"secret_ref"`
}

func DecodeUpstreamSpec(generic v1.UpstreamSpec) (*UpstreamSpec, error) {
	s := new(UpstreamSpec)
	if err := protoutil.UnmarshalStruct(generic, s); err != nil {
		return nil, err
	}
	return s, s.validateLambda()
}

func (s *UpstreamSpec) validateLambda() error {
	switch s.Region {
	case "us-east-2":
		return nil
	case "us-east-1":
		return nil
	case "us-west-1":
		return nil
	case "us-west-2":
		return nil
	case "ap-northeast-1":
		return nil
	case "ap-northeast-2":
		return nil
	case "ap-south-1":
		return nil
	case "ap-southeast-1":
		return nil
	case "ap-southeast-2":
		return nil
	case "ca-central-1":
		return nil
	case "cn-north-1":
		return nil
	case "eu-central-1":
		return nil
	case "eu-west-1":
		return nil
	case "eu-west-2":
		return nil
	case "eu-west-3":
		return nil
	case "sa-east-1":
		return nil
	}
	return errors.New("no such region")
}

func (s *UpstreamSpec) GetLambdaHostname() string {
	return fmt.Sprintf("lambda.%s.amazonaws.com", s.Region)
}

type FunctionSpec struct {
	FunctionName string `json:"function_name"`
	Qualifier    string `json:"qualifier"`
}

func DecodeFunctionSpec(generic v1.FunctionSpec) (*FunctionSpec, error) {
	s := new(FunctionSpec)
	if err := protoutil.UnmarshalStruct(generic, s); err != nil {
		return nil, err
	}
	return s, s.ValidateLambda()
}

func (s *FunctionSpec) ValidateLambda() error {
	if s.FunctionName == "" {
		return errors.New("invalid function name")
	}
	return nil
}
