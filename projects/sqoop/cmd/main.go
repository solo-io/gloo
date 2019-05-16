package main

import (
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-projects/projects/sqoop/pkg/setup"
)

func main() {
	if err := setup.Main(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
