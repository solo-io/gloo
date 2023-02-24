package kube2e

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	runtimeCheck "runtime"
	"strings"

	"github.com/solo-io/go-utils/stringutils"
	yamlHelper "gopkg.in/yaml.v2"

	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/test/services"
)

// GetYamlData gets the kind and metadata of a yaml file - used for creation of resources
func GetYamlData(yamlAssetsPath string, filename string) (kind string, name string, namespace string) {
	type KubernetesStruct struct {
		Kind     string `yaml:"kind"`
		Metadata struct {
			Name      string `yaml:"name"`
			Namespace string `yaml:"namespace"`
		} `yaml:"metadata"`
	}

	var kubernetesValues KubernetesStruct

	file, err := ioutil.ReadFile(filepath.Join(yamlAssetsPath, filename))
	if err != nil {
		fmt.Println(err.Error())
	}

	err = yamlHelper.Unmarshal(file, &kubernetesValues)
	if err != nil {
		fmt.Println(err.Error())
	}
	return kubernetesValues.Kind, kubernetesValues.Metadata.Name, kubernetesValues.Metadata.Namespace
}

// ApplyK8sResources gets all yaml files from a specified directory in upgrade/assets
// sort these files by kubernetes resource type and determine if there are architecture (arm vs amd) specific files
// Resource type is determined from the kubectl yaml config for all resources other than Deploy
// Deployment names can be specific to arm or amd - if there is a specific arm deployment include a deploy yaml with _arm_ as part of its name to only use this deploy on arm machines
// Then apply resources in a specific order
func ApplyK8sResources(yamlAssetsPath string, directory string) {
	fmt.Printf("\n=============== %s Directory Resources ===============\n", directory)
	files, err := ioutil.ReadDir(filepath.Join(yamlAssetsPath, directory))
	Expect(err).NotTo(HaveOccurred())

	// final map of files for resources
	fileMap := make(map[string][]string)

	// map of deployments
	deploymentMap := make(map[string][]string)
	for _, f := range files {

		// check filetype
		fileEnding := strings.Split(f.Name(), ".")
		if fileEnding[len(fileEnding)-1] != "yaml" {
			continue
		}
		kind, _, _ := GetYamlData(yamlAssetsPath, filepath.Join(directory, f.Name()))

		if kind == "Deployment" {
			splitName := strings.Split(f.Name(), "_")
			if deploymentMap[splitName[0]] == nil {
				deploymentMap[splitName[0]] = []string{f.Name()}
			} else {
				deploymentMap[splitName[0]] = append(deploymentMap[splitName[0]], f.Name())
			}
			// create a placeholder for the deployments
			if fileMap[kind] == nil {
				fileMap[kind] = []string{}
			}
		} else {
			if fileMap[kind] == nil {
				fileMap[kind] = []string{f.Name()}
			} else {
				fileMap[kind] = append(fileMap[kind], f.Name())
			}
		}
	}
	// handle deployment map
	// there may be arm specific deployments - if arm, use those
	// If there are not, use all deployments
	for _, filesList := range deploymentMap {
		var filename string
		for _, file := range filesList {
			splitName := strings.Split(file, "_")
			if stringutils.ContainsString("arm", splitName) {
				if runtimeCheck.GOARCH == "arm64" {
					filename = file
					break // return early as you only need one deployment, and we found an arm deployment while running on an arm64 machine
				}
			} else {
				filename = file
			}
		}
		fileMap["Deployment"] = append(fileMap["Deployment"], filename)
	}

	// order of resource creation
	creationOrder := []string{"Deployment", "Service", "Upstream", "AuthConfig", "RateLimitConfig", "VirtualService"}
	for _, kind := range creationOrder {
		filePaths := fileMap[kind]
		if nil != filePaths {
			for _, filepath := range filePaths {
				ApplyK8sResource(yamlAssetsPath, directory, filepath)
			}
		}
	}
}

// ApplyK8sResource creates resources from yaml files in the upgrade/assets folder and validates they have been created
func ApplyK8sResource(yamlAssetsPath string, subDirectory string, filename string) {
	filePath := filepath.Join(subDirectory, filename)
	kind, name, namespace := GetYamlData(yamlAssetsPath, filePath)
	fmt.Printf("Creating %s %s in namespace: %s", kind, name, namespace)
	RunAndCleanCommand("kubectl", "apply", "-f", filepath.Join(yamlAssetsPath, filePath))

	// validate resource creation
	switch kind {
	case "Service":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "svc/"+name, "-n", namespace)
		}, "20s", "1s").ShouldNot(BeEmpty())
		fmt.Printf(" (✓)\n")
	case "Deployment":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "deploy/"+name, "-n", namespace)
		}, "20s", "1s").ShouldNot(BeEmpty())
		fmt.Printf(" (✓)\n")
	case "Upstream":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "us/"+name, "-n", namespace)
		}, "20s", "1s").ShouldNot(BeEmpty())
		fmt.Printf(" (✓)\n")
	case "RateLimitConfig":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "ratelimitconfig/"+name, "-n", namespace, "-o", "jsonpath={.status.state}")
		}, "20s", "1s").Should(Equal("ACCEPTED"))
		fmt.Printf(" (✓)\n")
	case "VirtualService":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "vs/"+name, "-n", namespace, "-o", "jsonpath={.status.statuses."+namespace+".state}")
		}, "20s", "1s").Should(Equal("Accepted"))
		fmt.Printf(" (✓)\n")
	case "AuthConfig":
		Eventually(func() (string, error) {
			return services.KubectlOut("get", "AuthConfig/"+name, "-n", namespace, "-o", "jsonpath={.status.statuses."+namespace+".state}")
		}, "20s", "1s").Should(Equal("Accepted"))
		fmt.Printf(" (✓)\n")
	default:
		fmt.Printf(" : No validation found for yaml kind: %s\n", kind)
	}
}
