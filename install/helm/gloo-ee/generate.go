package main

import (
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
)

var glooEGenerationFiles = &generate.GenerationFiles{
	Artifact:             generate.GlooE,
	ValuesTemplate:       "install/helm/gloo-ee/values-template.yaml",
	ValuesOutput:         "install/helm/gloo-ee/values.yaml",
	DistributionOutput:   "install/distribution/values.yaml",
	ChartTemplate:        "install/helm/gloo-ee/Chart-template.yaml",
	ChartOutput:          "install/helm/gloo-ee/Chart.yaml",
	RequirementsTemplate: "install/helm/gloo-ee/requirements-template.yaml",
	RequirementsOutput:   "install/helm/gloo-ee/requirements.yaml",
}

// TODO ---- Use real values for Gloo w/ UI artifact
var glooOsWithReadOnlyUiGenerationFiles = &generate.GenerationFiles{
	Artifact:             generate.GlooWithRoUi,
	ValuesTemplate:       "install/helm/gloo-ee/values-template.yaml",
	ValuesOutput:         "install/helm/gloo-ee/values.yaml",
	DistributionOutput:   "install/distribution/values.yaml",
	ChartTemplate:        "install/helm/gloo-ee/Chart-template.yaml",
	ChartOutput:          "install/helm/gloo-ee/Chart.yaml",
	RequirementsTemplate: "install/helm/gloo-ee/requirements-template.yaml",
	RequirementsOutput:   "install/helm/gloo-ee/requirements.yaml",
}

func main() {

	// TODO: pass config, not just filesets
	args := &generate.GenerationArguments{}
	if err := generate.GetArguments(args); err != nil {
		log.Fatalf("unable to get valid generation arguments: %v", err)
	}
	//err := generate.Run(args, glooEGenerationFiles, glooOsWithReadOnlyUiGenerationFiles)
	err := generate.Run(args, glooEGenerationFiles)
	if err != nil {
		log.Fatalf("error while running generation: %v", err)
	}
}
