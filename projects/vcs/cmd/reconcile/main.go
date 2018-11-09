package main

import (
	"context"
	"fmt"
	"log"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/file"
	"github.com/solo-io/solo-kit/projects/vcs/pkg/syncer"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	printUsage()
	ctx := contextutils.WithLogger(context.Background(), "vcs")
	dc, err := file.NewDualClient("kube")
	if err != nil {
		return err
	}

	fmt.Println("begin applying changes")
	syncer.ApplyVcsToDeployment(ctx, dc)
	fmt.Println("done applying changes")

	return nil
}

func printUsage() {
	fmt.Printf(`Usage of this demo script
Run the dualClient script to create a local file system
Edit the filesystem
Run this script to sync your local changes to kubernetes
(this script will be moved to an e2e eventually)
----------------------------
`)
}
