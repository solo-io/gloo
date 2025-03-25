package envoy

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/pflag"
)

type Options struct {
	*options.Options
	InputFile          string
	OutputDir          string
	FolderPerNamespace bool
	Stats              bool
}

func (o *Options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.InputFile, "input-file", "", "File to convert")
	flags.BoolVar(&o.Stats, "stats", false, "Print stats about the conversion")
	flags.StringVar(&o.OutputDir, "_output", "./_output",
		"Where to write files")
	flags.BoolVar(&o.FolderPerNamespace, "folder-per-namespace", false,
		"Organize files by namespace")
}
