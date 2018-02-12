package main

import (
	"flag"

	//register CRDs
	"github.com/solo-io/glue/internal/bootstrap"
	_ "github.com/solo-io/glue/internal/install"
)

func main() {
	opts := bootstrap.Options{}
	flag.StringVar()
}
