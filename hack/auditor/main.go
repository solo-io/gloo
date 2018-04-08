package main

import (
	"flag"
	"fmt"
	"net/http"
	"github.com/solo-io/gloo/hack/auditor/audit"
	"github.com/solo-io/gloo/pkg/log"
)

func main() {
	port := flag.Uint("port", 8080, "listener port")
	flag.Parse()
	m := audit.NewServeMux()
	log.Printf("listening")
	log.Fatalf("%v", http.ListenAndServe(fmt.Sprintf(":%d", *port), m))
}
