package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	v1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/kind/pkg/docker"
)

const (
	glooeYaml    = "glooe.yaml"
	manifest     = "install/distribution/" + glooeYaml
	output       = "_output"
	distribution = output + "/distribution"
	deployment   = "Deployment"
	glooctl      = "glooctl"
	tarFile      = ".tar"
)

var (
	version         string
	distributionDir string
)

func main() {
	if len(os.Args) < 2 {
		panic("Must provide version as argument")
	} else {
		version = os.Args[1]
		distributionDir = filepath.Join(distribution, version)
	}
	if err := prepareWorkspace(); err != nil {
		log.Fatal(err)
	}
	specs, err := readManifestIntoParts()
	if err != nil {
		log.Fatal(err)
	}
	deployments, err := extractContainerImagesFromSpecs(specs)
	if err != nil {
		log.Fatal(err)
	}

	if err := saveImages(deployments); err != nil {
		log.Fatal(err)
	}

	if err := copyManifest(); err != nil {
		log.Fatal(err)
	}

	if err := copyBinaries(); err != nil {
		log.Fatal(err)
	}
}

func prepareWorkspace() error {
	if _, err := os.Stat(distributionDir); os.IsNotExist(err) {
		err = os.MkdirAll(distributionDir, 0755)
		if err != nil {
			return err
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
	savedImage := filepath.Join(distributionDir, splitImage[len(splitImage)-1]) + tarFile
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
			log.Printf("Successfully saved image: (%s)", image.Image)
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
	destinationManifest := filepath.Join(distributionDir, glooeYaml)
	if err := copyFile(manifest, destinationManifest); err != nil {
		return err
	}
	log.Print("Successfully copied local manifest to distribution directory")
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
			dest := filepath.Join(distributionDir, file.Name())
			if err := copyFile(source, dest); err != nil {
				return err
			}
			log.Printf("Successfully copied binary: (%s) to distribution directory", file.Name())
		}
	}
	return nil
}
