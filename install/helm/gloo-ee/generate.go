package main

import (
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
)

var glooEGenerationFiles = &generate.GenerationFiles{
	Artifact:             generate.GlooE,
	ValuesTemplate:       "install/helm/gloo-ee/values-template.yaml",
	ValuesOutput:         "install/helm/gloo-ee/values.yaml",
	DocsOutput:           "install/helm/gloo-ee/reference/values.txt",
	ChartTemplate:        "install/helm/gloo-ee/Chart-template.yaml",
	ChartOutput:          "install/helm/gloo-ee/Chart.yaml",
	RequirementsTemplate: "install/helm/gloo-ee/requirements-template.yaml",
	RequirementsOutput:   "install/helm/gloo-ee/requirements.yaml",
}

func main() {
	args := &generate.GenerationArguments{}
	if err := generate.GetArguments(args); err != nil {
		log.Fatalf("unable to get valid generation arguments: %v", err)
	}
	log.Printf("Running generate with args: %v", args)
	err := generate.Run(args, glooEGenerationFiles)
	if err != nil {
		log.Fatalf("error while running generation: %v", err)
	}
}
