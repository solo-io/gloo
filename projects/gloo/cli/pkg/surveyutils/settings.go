package surveyutils

import (
	"github.com/solo-io/solo-projects/pkg/cliutil"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
)

func AddSettingsExtAuthFlagsInteractive(opts *options.OIDCSettings) error {

	if err := cliutil.GetStringInput("name of the extauth server upstream: ", &opts.ExtAtuhServerUpstreamRef.Name); err != nil {
		return err
	}
	if err := cliutil.GetStringInput("namespace of the extauth server upstream: ", &opts.ExtAtuhServerUpstreamRef.Namespace); err != nil {
		return err
	}
	return nil
}
