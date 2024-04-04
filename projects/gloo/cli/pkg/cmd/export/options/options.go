package options

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
)

type Options struct {
	*options.Options

	// OutputDir is the name of the directory where the exported artifacts will be stored
	// If one is not specified, a tmp directory will be chosen
	OutputDir string
}
