package main

import (
	"flag"
	"fmt"
	"net/http"
	"github.com/golang/glog"
)

const (
	hello = "Hi, there!\n"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	if _, err := w.Write([]byte(hello)); err != nil {
		glog.Error(err)
	}
}

// RunHTTP opens a simple listener on the port.
func RunHTTP(port uint) {
	glog.Infof("upstream listening HTTP1.1 on %d", port)
	h := handler{}
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: h}
	if err := server.ListenAndServe(); err != nil {
		glog.Fatal(err)
	}
}

func main() {
	port := flag.Uint("port", 8080, "listener port")
	flag.Parse()
	RunHTTP(*port)
}
