package function

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

type Spec struct {
	FunctionName string
	Qualifier    string
}

func FromMap(m map[string]interface{}) (*Spec, error) {
	s := new(Spec)
	mapstructure.Decode(m, s)
	return s, s.ValidateLambda()
}

func (s *Spec) ValidateLambda() error {
	if s.FunctionName == "" {
		return errors.New("invalid function name")
	}
	return nil
}
