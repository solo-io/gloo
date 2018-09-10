package main

import (
	"log"
	"os"
	"syscall"

	"github.com/solo-io/envoy-operator/pkg/downward"
)

func main() {
	inputfile := "/etc/envoy/envoy.yaml"
	outfile := "/tmp/envoy.yaml"

	transformer := downward.NewTransformer()
	err := transformer.TransformFiles(inputfile, outfile)
	if err != nil {
		log.Fatalf("initializer failed: %v", err)
	}
	env := os.Environ()
	args := []string{
		"/usr/local/bin/envoy", "-c", outfile, "--v2-config-only",
	}
	if err := syscall.Exec(args[0], args, env); err != nil {
		panic(err)
	}

}
