package main

import (
	"github.com/solo-io/solo-kit/pkg/code-generator/cmd"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-projects/pkg/version"
)

//go:generate go run generate.go

func main() {
	err := version.CheckVersions()
	if err != nil {
		log.Fatalf("generate failed!: %v", err)
	}
	log.Printf("Starting generate...")
	if err := cmd.Run(".", true, true, nil, []string{"./projects/gloo/pkg/"}); err != nil {
		log.Fatalf("generate failed!: %v", err)
	}
}
