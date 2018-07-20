package main

import (
	"log"

	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/code-generator/protoc"
)

func main() {
	if err := protokit.RunPlugin(new(protoc.Plugin)); err != nil {
		log.Fatal(err)
	}
}
