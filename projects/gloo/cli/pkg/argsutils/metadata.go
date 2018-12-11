package argsutils

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/surveyutils"
)

// TODO: split the validation piece from the interactive piece
func MetadataArgsParse(opts *options.Options, args []string) error {
	if !opts.Top.Interactive {
		// handle static args
		if opts.Metadata.Namespace == "" {
			return errors.Errorf("namespace must be specified")
		}
		if opts.Metadata.Name == "" {
			if len(args) == 0 {
				return errors.Errorf("name must be specified in flag (--name) or via first arg")
			} else {
				opts.Metadata.Name = args[0]
			}
		}
	} else {
		// handle dynamic args
		metadata := core.Metadata{}
		surveyutils.MetadataSurvey(&metadata)
		opts.Metadata = metadata
	}
	return nil
}
