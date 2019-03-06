package flagutils

import (
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/pflag"
)

func AddInstallFlags(set *pflag.FlagSet, install *options.InstallExtended) {
	set.StringVar(&install.LicenseKey, "license-key", "", "License key to activate GlooE features")
}
