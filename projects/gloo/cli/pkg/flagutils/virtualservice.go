package flagutils

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/pflag"
)

func AddVirtualServiceFlags(set *pflag.FlagSet, vs *options.InputVirtualService) {
	addDomainsFlag(set, &vs.Domains)
}

func addDomainsFlag(set *pflag.FlagSet, ptr *[]string) {
	set.StringSliceVar(ptr, "domains", []string{}, "comma seperated list of domains")
}
