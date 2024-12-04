package flagutils

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/pflag"
)

func AddDebugFlags(set *pflag.FlagSet, top *options.Top) {
	set.BoolVar(&top.Zip, "zip", false, "save logs to a tar file (specify location with -f)")
	set.BoolVar(&top.ErrorsOnly, "errors-only", false, "filter for error logs only")
}

func AddDirectoryFlag(set *pflag.FlagSet, strptr *string) {
	set.StringVarP(strptr, "directory", "d", "debug", "directory to write debug info to")
}

func AddNamespacesFlag(set *pflag.FlagSet, strptr *[]string) {
	set.StringArrayVarP(strptr, "namespaces", "N", []string{DefaultNamespace}, "namespaces from which to dump logs and resources (use flag multiple times to specify multiple namespaces, e.g. '-N gloo-system -N default')")
}
