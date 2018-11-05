package main

import (
	"log"
	"os"

	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/code-generator/protoc"
)

func main() {
	outputDescriptors := os.Getenv("OUTPUT_DESCRIPTORS") == "1"
	plugin := &protoc.Plugin{OutputDescriptors: outputDescriptors}
	// use this to debug without running protoc
	if os.Getenv("DEBUG") == "1" {
		f, err := os.Open("projects/supergloo/api/v1/project.json.descriptors")
		if err != nil {
			log.Fatal(err)
		}
		return
		if err := protokit.RunPluginWithIO(plugin, f, os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
	if err := protokit.RunPlugin(plugin); err != nil {
		log.Fatal(err)
	}
}
