package rawgetter

import (
	"fmt"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	NoResourceVersionError = func(ref *core.ResourceRef) error {
		return errors.New(fmt.Sprintf("Must provide a resource version in the edited YAML for resource %v in namespace %v", ref.Name, ref.Namespace))
	}

	EditedRefError = func(ref *core.ResourceRef) error {
		return errors.New(fmt.Sprintf("Cannot change the resource's name or namespace in the edited YAML for resource %v in namespace %v", ref.Name, ref.Namespace))
	}

	FailedToReadCrdSpec = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to read spec in edited YAML for resource %v in namespace %v", ref.Name, ref.Namespace)
	}

	UnmarshalError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to parse the edited YAML for resource %v in namespace %v", ref.Name, ref.Namespace)
	}
)
