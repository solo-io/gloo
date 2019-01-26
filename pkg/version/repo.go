package version

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

const (
	gopkgToml    = "Gopkg.toml"
	constraint   = "constraint"
	nameConst    = "name"
	versionConst = "version"

	GlooPkg    = "github.com/solo-io/gloo"
	SoloKitPkg = "github.com/solo-io/solo-kit"
)

func PinGitVersion(relativeRepoDir string, version string) error {
	tag := GetTag(version)
	cmd := exec.Command("git", "checkout", tag)
	cmd.Dir = relativeRepoDir
	buf := &bytes.Buffer{}
	out := io.MultiWriter(buf, os.Stdout)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "%v failed: %s", cmd.Args, buf.String())
	}
	return nil
}

func GetGitVersion(relativeRepoDir string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--dirty")
	cmd.Dir = relativeRepoDir
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrapf(err, "%v failed: %s", cmd.Args, output)
	}
	return strings.TrimSpace(string(output)), nil
}

func GetTag(version string) string {
	if strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}

func GetVersion(pkgName string, tomlTree []*toml.Tree) (string, error) {
	for _, v := range tomlTree {
		if v.Get(nameConst) == pkgName && v.Get(versionConst) != "" {
			return v.Get(versionConst).(string), nil
		}
	}
	return "", fmt.Errorf("unable to find version for %s in toml", pkgName)
}

func ParseToml() ([]*toml.Tree, error) {
	config, err := toml.LoadFile(gopkgToml)
	if err != nil {
		return nil, err
	}

	tomlTree := config.Get(constraint)

	switch typedTree := tomlTree.(type) {
	case []*toml.Tree:
		return typedTree, nil
	default:
		return nil, fmt.Errorf("unable to parse toml tree")
	}
}
