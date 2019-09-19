package main

import (
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
)

var glooEGenerationFiles = &generate.GenerationFiles{
	Artifact:             generate.GlooE,
	ValuesTemplate:       "install/helm/gloo-ee/values-template.yaml",
	ValuesOutput:         "install/helm/gloo-ee/values.yaml",
	ChartTemplate:        "install/helm/gloo-ee/Chart-template.yaml",
	ChartOutput:          "install/helm/gloo-ee/Chart.yaml",
	RequirementsTemplate: "install/helm/gloo-ee/requirements-template.yaml",
	RequirementsOutput:   "install/helm/gloo-ee/requirements.yaml",
}

var glooOsWithReadOnlyUiGenerationFiles = &generate.GenerationFiles{
	Artifact:             generate.GlooWithRoUi,
	ValuesTemplate:       "install/helm/gloo-os-with-ui/values-template.yaml",
	ValuesOutput:         "install/helm/gloo-os-with-ui/values.yaml",
	ChartTemplate:        "install/helm/gloo-os-with-ui/Chart-template.yaml",
	ChartOutput:          "install/helm/gloo-os-with-ui/Chart.yaml",
	RequirementsTemplate: "install/helm/gloo-os-with-ui/requirements-template.yaml",
	RequirementsOutput:   "install/helm/gloo-os-with-ui/requirements.yaml",
}

func main() {
	args := &generate.GenerationArguments{}
	if err := generate.GetArguments(args); err != nil {
		log.Fatalf("unable to get valid generation arguments: %v", err)
	}
	err := generate.Run(args, glooEGenerationFiles, glooOsWithReadOnlyUiGenerationFiles)
	if err != nil {
		log.Fatalf("error while running generation: %v", err)
	}
}
