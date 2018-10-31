package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/setup"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	port := flag.Int("p", 8083, "port to bind")
	fmt.Println("running")
	ctx := contextutils.WithLogger(context.Background(), "vcs")
	return setup.Setup(ctx, *port)
}
