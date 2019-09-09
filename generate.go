package main

import (
	"os"

	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-kit/pkg/code-generator/cmd"
	"github.com/solo-io/solo-kit/pkg/code-generator/docgen/options"
	"github.com/solo-io/solo-projects/pkg/version"
)

//go:generate go run generate.go

func main() {
	err := version.CheckVersions()
	if err != nil {
		log.Fatalf("generate failed!: %v", err)
	}
	log.Printf("Starting generate...")

	generateOptions := cmd.GenerateOptions{
		SkipGenMocks:  true,
		CustomImports: []string{os.ExpandEnv("$GOPATH/src/github.com/solo-io/gloo/projects/gloo/api/external")},
		SkipDirs:      []string{"./projects/gloo/pkg/", "./projects/gloo-ui/"},
		RelativeRoot:  ".",
		CompileProtos: true,
		GenDocs: &cmd.DocsOptions{
			Output: options.Hugo,
		},
	}
	if err := cmd.Generate(generateOptions); err != nil {
		log.Fatalf("generate failed!: %v", err)
	}
}
