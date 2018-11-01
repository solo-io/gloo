package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/file"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	action := flag.String("a", "vs", "action to take")
	fmt.Printf("running, action: %v\n", *action)
	ctx := contextutils.WithLogger(context.Background(), "vcs")
	dc, err := file.NewDualClient("kube")
	if err != nil {
		return err
	}
	file.GenerateFilesystem(ctx, "gloo-system", dc)
	return nil
}
