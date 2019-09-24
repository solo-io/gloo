package version

import (
	"github.com/pkg/errors"
	glooversion "github.com/solo-io/gloo/pkg/version"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/versionutils"
	"github.com/solo-io/go-utils/versionutils/dep"
	"github.com/solo-io/go-utils/versionutils/git"
)

const GlooPkg = "github.com/solo-io/gloo"

var attributeTypes = map[dep.VersionType]string{
	dep.Version:  "version",
	dep.Branch:   "branch",
	dep.Revision: "revision",
}

func CheckVersions() error {
	tomlTree, err := versionutils.ParseFullToml()
	if err != nil {
		return err
	}

	log.Printf("Checking expected solo kit and gloo versions...")
	expectedGlooVersion, err := versionutils.GetDependencyVersionInfo(GlooPkg, tomlTree)
	if err != nil {
		return err
	}
	log.Printf("Expecting gloo with %s [%s]", attributeTypes[expectedGlooVersion.Type], expectedGlooVersion.Version)

	expectedSoloKitVersion, err := versionutils.GetDependencyVersionInfo(glooversion.SoloKitPkg, tomlTree)
	if err != nil {
		return err
	}
	log.Printf("Expecting solo-kit with %s [%s]", attributeTypes[expectedSoloKitVersion.Type], expectedSoloKitVersion.Version)

	log.Printf("Checking repo versions...")
	actualGlooVersion, err := git.GetGitRefInfo("../gloo")
	if err != nil {
		return err
	}
	log.Printf("Found gloo ref. Tag [%s], Branch [%s], Commit [%s]",
		actualGlooVersion.Tag, actualGlooVersion.Branch, actualGlooVersion.Hash)

	actualSoloKitVersion, err := git.GetGitRefInfo("../solo-kit")
	if err != nil {
		return err
	}
	log.Printf("Found solo-kit ref. Tag [%s], Branch [%s], Commit [%s]",
		actualSoloKitVersion.Tag, actualSoloKitVersion.Branch, actualSoloKitVersion.Hash)

	if err := compare(expectedGlooVersion, actualGlooVersion, "gloo"); err != nil {
		return err
	}

	if err := compare(expectedSoloKitVersion, actualSoloKitVersion, "solo-kit"); err != nil {
		return err
	}

	log.Printf("Versions are pinned correctly.")
	return nil
}

func compare(expected *dep.VersionInfo, actual *git.RefInfo, packageName string) error {
	switch expected.Type {
	case dep.Version:
		expectedTaggedVersion := git.AppendTagPrefix(expected.Version)
		if actual.Tag != expectedTaggedVersion {
			return errors.Errorf("Expected %s tag [%s], found tag [%s] in repo. "+
				"Run 'make pin-repos' or fix manually.", packageName, expectedTaggedVersion, actual.Tag)
		}
	case dep.Branch:
		if actual.Branch != expected.Version {
			return errors.Errorf("Expected %s branch [%s], found branch [%s] in repo. "+
				"Run 'make pin-repos' or fix manually.", packageName, expected.Version, actual.Branch)
		}
	case dep.Revision:
		if actual.Hash != expected.Version {
			return errors.Errorf("Expected %s revision [%s], found commit [%s] in repo. "+
				"Run 'make pin-repos' or fix manually.", packageName, expected.Version, actual.Hash)
		}
	default:
		return errors.Errorf("Unexpected dep version attribute type: [%d]", expected.Type)
	}
	return nil
}
