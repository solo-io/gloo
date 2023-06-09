package envoy

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo/v2"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/test/services/envoy"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/gloo/test/testutils/version"
	"github.com/solo-io/skv2/codegen/util"
)

func NewFactory() envoy.Factory {
	var err error

	// if an envoy binary is explicitly specified, use it
	envoyPath := os.Getenv(testutils.EnvoyBinary)
	if envoyPath != "" {
		log.Printf("Using envoy from environment variable: %s", envoyPath)
		return envoy.NewLinuxFactory(bootstrapTemplate, envoyPath, "")
	}

	switch runtime.GOOS {
	case "darwin":
		log.Printf("Using docker to Run envoy")

		image := fmt.Sprintf("quay.io/solo-io/gloo-ee-envoy-wrapper:%s", mustGetEnvoyWrapperTag())
		return envoy.NewDockerFactory(bootstrapTemplate, image)

	case "linux":
		var tmpDir string

		// try to grab one form docker...
		tmpDir, err = os.MkdirTemp(os.Getenv("HELPER_TMP"), "envoy")
		Expect(err).NotTo(HaveOccurred())

		envoyImage := fmt.Sprintf("gcr.io/gloo-ee/envoy-gloo-ee:%s", mustGetEnvoyGlooTag())
		log.Printf("Using envoy docker image: %s", envoyImage)

		bash := fmt.Sprintf(`
set -ex
CID=$(docker run -d  %s /bin/bash -c exit)

# just print the image sha for repoducibility
echo "Using Envoy Image:"
docker inspect %s -f "{{.RepoDigests}}"

docker cp $CID:/usr/local/bin/envoy .
docker rm $CID
    `, envoyImage, envoyImage)
		scriptfile := filepath.Join(tmpDir, "getenvoy.sh")

		os.WriteFile(scriptfile, []byte(bash), 0755)

		cmd := exec.Command("bash", scriptfile)
		cmd.Dir = tmpDir
		cmd.Stdout = GinkgoWriter
		cmd.Stderr = GinkgoWriter
		err = cmd.Run()
		Expect(err).NotTo(HaveOccurred())

		return envoy.NewLinuxFactory(bootstrapTemplate, filepath.Join(tmpDir, "envoy"), tmpDir)

	default:
		Fail("Unsupported OS: " + runtime.GOOS)
	}
	return nil
}

// mustGetEnvoyGlooTag returns the tag of the envoy-gloo-ee image which will be executed
// The tag is chosen using the following process:
//  1. If ENVOY_IMAGE_TAG is defined, use that tag
//  2. If not defined, use the ENVOY_GLOO_IMAGE_VERSION tag defined in the Makefile
func mustGetEnvoyGlooTag() string {
	eit := os.Getenv(testutils.EnvoyImageTag)
	if eit != "" {
		return eit
	}

	makefile := filepath.Join(util.GetModuleRoot(), "Makefile")
	inFile, err := os.Open(makefile)
	Expect(err).NotTo(HaveOccurred())

	defer inFile.Close()

	const prefix = "ENVOY_GLOO_IMAGE_VERSION ?= "

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}

	Fail("Could not determine envoy-gloo-ee tag. Find valid tag names here https://console.cloud.google.com/gcr/images/gloo-ee/GLOBAL/envoy-gloo-ee and update your make target or use the variable ENVOY_IMAGE_TAG to set it.")
	return ""
}

// mustGetEnvoyWrapperTag returns the tag of the gloo-ee-envoy-wrapper image which will be executed
// The tag is chosen using the following process:
//  1. If ENVOY_IMAGE_TAG is defined, use that tag
//  2. If not defined, use the latest released tag of that image
func mustGetEnvoyWrapperTag() string {
	eit := os.Getenv(testutils.EnvoyImageTag)
	if eit != "" {
		return eit
	}

	latestPatchVersion, err := version.GetLastReleaseOfCurrentBranch()
	if err != nil {
		Fail(errors.Wrap(err, "Failed to extract the latest release of current minor").Error())
	}

	return strings.TrimPrefix(latestPatchVersion.String(), "v")
}
