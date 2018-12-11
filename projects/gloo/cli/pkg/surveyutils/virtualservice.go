package surveyutils

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
)

func AddVirtualServiceFlagsInteractive(vs *options.InputVirtualService) error {
	if err := cliutil.GetStringSliceInput(
		fmt.Sprintf("Add another domain for this virtual service (empty to skip)? %v", vs.Domains),
		&vs.Domains,
	); err != nil {
		return err
	}

	return nil
}
