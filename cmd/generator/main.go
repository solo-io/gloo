package main

import (
	"log"
	"os"

	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/code-generator/protoc"
)

func main() {
	// use this to debug without running
	if os.Getenv("REAL") != "1" {
		f, err := os.Open("projects/supergloo/api/v1/project.json.descriptors")
		if err != nil {
			log.Fatal(err)
		}
		if err := protokit.RunPluginWithIO(new(protoc.Plugin), f, os.Stdout); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := protokit.RunPlugin(new(protoc.Plugin)); err != nil {
		log.Fatal(err)
	}
}
