package function

import (
	"github.com/mitchellh/mapstructure"
)

type Spec struct {
	Async bool
}

func FromMap(m interface{}) (*Spec, error) {
	s := new(Spec)
	err := mapstructure.Decode(m, s)
	return s, err
}
