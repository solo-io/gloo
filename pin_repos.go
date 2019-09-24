package main

import (
	glooversion "github.com/solo-io/gloo/pkg/version"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/versionutils"
	"github.com/solo-io/go-utils/versionutils/dep"
	"github.com/solo-io/go-utils/versionutils/git"
	"github.com/solo-io/solo-projects/pkg/version"
)

func main() {
	tomlTree, err := versionutils.ParseFullToml()
	fatalCheck(err, "parsing error")

	glooVersion, err := versionutils.GetDependencyVersionInfo(version.GlooPkg, tomlTree)
	fatalCheck(err, "getting gloo version")

	soloKitVersion, err := versionutils.GetDependencyVersionInfo(glooversion.SoloKitPkg, tomlTree)
	fatalCheck(err, "getting solo-kit version")

	fatalCheck(git.PinDependencyVersion("../gloo", handleTag(glooVersion)), "consider git fetching in gloo repo")

	fatalCheck(git.PinDependencyVersion("../solo-kit", handleTag(soloKitVersion)), "consider git fetching in solo-kit repo")
}

func handleTag(versionInfo *dep.VersionInfo) string {
	if versionInfo.Type == dep.Version {
		// If the toml version attribute is "version", we are looking for a tag
		return git.AppendTagPrefix(versionInfo.Version)
	}
	return versionInfo.Version
}

func fatalCheck(err error, msg string) {
	if err != nil {
		log.Fatalf("Error (%v) unable to pin repos!: %v", msg, err)
	}
}
