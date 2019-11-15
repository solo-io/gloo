package main

import (
	"fmt"

	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/versionutils"
	"github.com/solo-io/solo-projects/pkg/version"
)

// script used to get Gloo version from Gopkg.toml file and print out on command line
//   - currently used in CI to determine version of glooctl binary to download
func main() {
	tomlTree, err := versionutils.ParseFullToml()
	if err != nil {
		log.Fatalf("Error parsing Gopkg.toml, unable to determine Gloo version: %v", err)
	}
	glooVersion, err := versionutils.GetDependencyVersionInfo(version.GlooPkg, tomlTree)
	if err != nil {
		log.Fatalf("Error getting gloo version, unable to determine Gloo version! %v", err)
	}
	fmt.Println("v" + glooVersion.Version)
}
