package upstream

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Spec struct {
	Region string
	Secret string // TODO should this be a secret ref?
}

func FromMap(m map[string]interface{}) (*Spec, error) {
	s := new(Spec)
	mapstructure.Decode(m, s)
	return s, s.validateLambda()
}

func (s *Spec) validateLambda() error {

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

func (s *Spec) GetLambdaHostname() string {
	return fmt.Sprintf("lambda.%s.amazonaws.com", s.Region)
}
