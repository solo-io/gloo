package main

import (
	"context"
	"fmt"
	"log"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/file"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/syncer"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	printUsage()
	ctx := contextutils.WithLogger(context.Background(), "vcs")
	fc, err := file.NewFileClient("testRoot")
	if err != nil {
		return err
	}
	kc, err := file.NewKubeClient()
	if err != nil {
		return err
	}

	fmt.Println("begin applying changes")
	syncer.ApplyVcsToDeployment(ctx, fc, kc)
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
