package main

import (
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-kit/pkg/code-generator/cmd"
	"github.com/solo-io/solo-kit/pkg/code-generator/sk_anyvendor"
)

const GlooPkg = "github.com/solo-io/gloo"

//go:generate go run generate.go

func main() {
	log.Printf("Starting generate...")

	imports := sk_anyvendor.CreateDefaultMatchOptions(
		[]string{
			"projects/observability/**/*.proto",
			"projects/grpcserver/**/*.proto",
			sk_anyvendor.SoloKitMatchPattern,
		},
	)
	imports.External[GlooPkg] = []string{"projects/**/*.proto"}

	generateOptions := cmd.GenerateOptions{
		SkipGenMocks:    true,
		SkipDirs:        []string{"./projects/gloo/pkg/", "./projects/gloo-ui/"},
		RelativeRoot:    ".",
		CompileProtos:   true,
		ExternalImports: imports,
	}
	if err := cmd.Generate(generateOptions); err != nil {
		log.Fatalf("generate failed!: %v", err)
	}
	log.Printf("Finished generating code")
}
