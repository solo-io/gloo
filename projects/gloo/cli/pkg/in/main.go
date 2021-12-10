package main

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/bug-report/pkg"
	"log"
)

func main() {
	err := pkg.RootCmd(nil).Execute()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
