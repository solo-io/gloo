package validator

import (
	"fmt"

	"github.com/solo-io/glue/pkg/api/types"
)

type Validator struct{}

func (v *Validator) Validate(cfg types.Config) error {
	return fmt.Errorf("not implemented")
}
