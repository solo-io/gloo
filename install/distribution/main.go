package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/ghodss/yaml"
	"k8s.io/api/apps/v1"
	"sigs.k8s.io/kind/pkg/docker"
)

const (
	glooe           = "glooe"
	glooeYaml       = glooe + "-distribution.yaml"
	manifest        = "install/manifest/" + glooeYaml
	scriptsDir      = "install/distribution/scripts"
	output          = "_output"
	distribution    = output + "/distribution"
	tarDir          = output + "/distribution_tar"
	deployment      = "Deployment"
	glooctl         = "glooctl"
	tgzExt          = ".tgz"
	tarExt          = ".tar"
	setupScriptName = "setup"

	googleEnv = "GOOGLE_APPLICATION_CREDENTIALS"
)

var (
	version               string
	outputDistributionDir string
	id                    uuid.UUID

	logger *zap.SugaredLogger

	projectId = os.ExpandEnv("PROJECT_ID")
)

func init() {
	devLogger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger = devLogger.Sugar()
}

func main() {
	defer logger.Sync()
	if len(os.Args) < 2 {
		panic("Must provide version as argument")
	} else {
		version = os.Args[1]
		outputDistributionDir = filepath.Join(distribution, version)
	}
	if err := prepareWorkspace(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := prepareFiles(); err != nil {
		logger.Fatal(err.Error())
	}

	if os.Getenv(googleEnv) == "" {
		if err := localDistributionTarball(); err != nil {
			logger.Fatal(err.Error())
		}
		logger.Fatal("google cloud credentials were not present, therefore exiting")
	}

	distBucketCli, err := newDistributionBucketClient()
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := syncDataToBucket(distBucketCli); err != nil {
		logger.Fatal(err.Error())
	}

}

func prepareFiles() error {
	specs, err := readManifestIntoParts()
	if err != nil {
		return err
	}
	deployments, err := extractContainerImagesFromSpecs(specs)
	if err != nil {
		return err
	}

	if err := saveImages(deployments); err != nil {
		return err
	}

	if err := copyManifest(); err != nil {
		return err
	}

	if err := copyBinaries(); err != nil {
		return err
	}

	if err := copySetupScripts(); err != nil {
		return err
	}
	return nil
}

func prepareWorkspace() error {
	directories := []string{tarDir, outputDistributionDir}
	for _, v := range directories {
		if _, err := os.Stat(v); os.IsNotExist(err) {
			err = os.MkdirAll(v, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readManifestIntoParts() ([]string, error) {
	bytes, err := ioutil.ReadFile(manifest)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), "---"), nil
}

func extractContainerImagesFromSpecs(specs []string) ([]v1.Deployment, error) {
	deployments := make([]v1.Deployment, 0, len(specs))
	for _, spec := range specs {
		var dpl v1.Deployment
		err := yaml.Unmarshal([]byte(spec), &dpl)
		if err != nil {
			return nil, err
		}
		if dpl.Kind == deployment {
			deployments = append(deployments, dpl)
		}
	}
	return deployments, nil
}

func savedImageName(imageName string) string {
	splitImage := strings.Split(imageName, "/")
	savedImage := filepath.Join(outputDistributionDir, splitImage[len(splitImage)-1]) + tarExt
	return savedImage
}

func saveImages(deployments []v1.Deployment) error {
	for _, dpl := range deployments {
		containers := dpl.Spec.Template.Spec.Containers
		for _, image := range containers {
			var err error
			_, err = docker.PullIfNotPresent(image.Image, 0)
			if err != nil {
				return err
			}
			err = docker.Save(image.Image, savedImageName(image.Image))
			if err != nil {
				return err
			}
			logger.Infof("Successfully saved image: (%s)", image.Image)
		}
	}

	return nil
}

func copyFile(source, dest string) error {
	cmd := exec.Command("cp", source, dest)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func copyManifest() error {
	destinationManifest := filepath.Join(outputDistributionDir, glooeYaml)
	if err := copyFile(manifest, destinationManifest); err != nil {
		return err
	}
	logger.Infof("Successfully copied local manifest to distribution directory")
	return nil
}

func copyBinaries() error {
	info, err := ioutil.ReadDir(output)
	if err != nil {
		return err
	}
	for _, file := range info {
		if strings.Contains(file.Name(), glooctl) {
			source := filepath.Join(output, file.Name())
			dest := filepath.Join(outputDistributionDir, file.Name())
			if err := copyFile(source, dest); err != nil {
				return err
			}
			logger.Infof("Successfully copied binary: (%s) to distribution directory", file.Name())
		}
	}
	return nil
}

// Copy setup scripts from install/distribution/scripts to output/distribution
func copySetupScripts() error {
	for _, extension := range []string{"sh", "bat"} {
		filename := strings.Join([]string{setupScriptName, extension}, ".")
		source := filepath.Join(scriptsDir, filename)
		dest := filepath.Join(outputDistributionDir, filename)
		if err := copyFile(source, dest); err != nil {
			return err
		}
		logger.Infof("Successfully copied setup script (%s) to distribution directory", filename)
	}
	return nil
}

func tarballFileName() (string, error) {
	var err error
	id, err = uuid.NewRandom()
	if err != nil {
		return "", err
	}

	tarFile := fmt.Sprintf("%s%s%s", glooe, version, tgzExt)
	tarFilepath := filepath.Join(id.String(), tarFile)
	return tarFilepath, nil
}

func localDistributionTarball() error {
	tarFilepath, err := tarballFileName()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(tarDir, id.String()), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(tarDir, tarFilepath))
	if err != nil {
		return err
	}
	defer file.Close()

	if err := writeDistributionTarball(file); err != nil {
		return err
	}

	return nil
}

func writeDistributionTarball(wr io.Writer) error {
	cmd := exec.Command("tar", "-C", distribution, "-cz", ".")
	cmd.Stderr = os.Stderr

	cmd.Stdout = wr

	logger.Infof("creating tar.gz file")
	if err := cmd.Run(); err != nil {
		return err
	}
	logger.Info("successfully created tar.gz file")
	return nil
}
