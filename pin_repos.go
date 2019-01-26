package main

import (
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-projects/pkg/version"
)

func main() {
	tomlTree, err := version.ParseToml()
	if err != nil {
		fatal(err)
	}

	glooVersion, err := version.GetVersion(version.GlooPkg, tomlTree)
	if err != nil {
		fatal(err)
	}

	soloKitVersion, err := version.GetVersion(version.SoloKitPkg, tomlTree)
	if err != nil {
		fatal(err)
	}

	err = version.PinGitVersion("../gloo", glooVersion)
	if err != nil {
		fatal(err)
	}
	err = version.PinGitVersion("../solo-kit", soloKitVersion)
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	log.Fatalf("unable to pin repos!: %v", err)
}