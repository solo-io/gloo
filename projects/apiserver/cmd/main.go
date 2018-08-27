package main

import (
	"flag"

	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/setup"
)

func main() {
	port := flag.Int("p", 8080, "port to bind")
	log.Printf("listening on :%v", *port)
	if err := setup.Setup(*port); err != nil {
		log.Fatalf("%v", err)
	}
}
