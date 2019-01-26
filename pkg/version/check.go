package version

import (
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/log"
)

func CheckVersions() error {
	log.Printf("Checking expected solo kit and gloo versions...")
	tomlTree, err := ParseToml()
	if err != nil {
		return err
	}

	expectedGlooVersion, err := GetVersion(GlooPkg, tomlTree)
	if err != nil {
		return err
	}

	expectedSoloKitVersion, err := GetVersion(SoloKitPkg, tomlTree)
	if err != nil {
		return err
	}

	log.Printf("Checking repo versions...")
	actualGlooVersion, err := GetGitVersion("../gloo")
	if err != nil {
		return err
	}
	expectedTaggedGlooVersion := GetTag(expectedGlooVersion)
	if expectedTaggedGlooVersion != actualGlooVersion {
		return errors.Errorf("Expected gloo version %s, found gloo version %s in repo. Run 'make pin-repos' or fix manually.", expectedTaggedGlooVersion, actualGlooVersion)
	}

	actualSoloKitVersion, err := GetGitVersion("../solo-kit")
	if err != nil {
		return err
	}
	expectedTaggedSoloKitVersion := GetTag(expectedSoloKitVersion)
	if expectedTaggedSoloKitVersion != actualSoloKitVersion {
		return errors.Errorf("Expected solo kit version %s, found solo kit version %s in repo. Run 'make pin-repos' or fix manually.", expectedTaggedSoloKitVersion, actualSoloKitVersion)
	}
	log.Printf("Versions are pinned correctly.")
	return nil
}
