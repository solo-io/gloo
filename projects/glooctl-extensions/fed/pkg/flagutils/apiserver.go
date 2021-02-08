package flagutils

import (
	"github.com/solo-io/solo-projects/projects/glooctl-extensions/fed/pkg/cmd/options"
	"github.com/spf13/pflag"
)

func AddApiserverFlags(set *pflag.FlagSet, opt *options.Options) {
	set.StringVarP(&opt.ApiserverPort, "port", "p", "10101", "gloo fed console apiserver port")
}
