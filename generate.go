package main

import (
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-kit/pkg/code-generator/cmd"
	"github.com/solo-io/solo-kit/pkg/code-generator/sk_anyvendor"
)

const (
	GlooPkg     = "github.com/solo-io/gloo"
	SoloAPIsPkg = "github.com/solo-io/solo-apis"
	Skv2Pkg     = "github.com/solo-io/skv2"
)

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
	imports.External[GlooPkg] = []string{
		"projects/**/*.proto",
	}
	// Import rate limit API
	imports.External[SoloAPIsPkg] = []string{
		"api/rate-limiter/**/*.proto",
	}
	// Import skv2 APIs
	imports.External[Skv2Pkg] = []string{
		"api/**/*.proto",
	}

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
