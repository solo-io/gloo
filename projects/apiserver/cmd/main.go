package main

import (
	"flag"
	"log"

	"github.com/solo-io/solo-kit/projects/apiserver/pkg/setup"
)

func main() {
	port := flag.Int("p", 8082, "port to bind")
	flag.Parse()
	log.Printf("listening on :%v", *port)
	if err := setup.Setup(*port); err != nil {
		log.Fatalf("%v", err)
	}
}
