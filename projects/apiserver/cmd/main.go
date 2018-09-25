package main

import (
	"flag"
	"log"
	"context"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/setup"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/setup"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/setup"
	"github.com/solo-io/solo-kit/test/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	port := flag.Int("p", 8082, "port to bind")
	dev := flag.Bool("dev", false, "use memory instead of connecting to real gloo storage")
	flag.Parse()
	glooOpts, err := config.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	gatewayOpts, err := gatewaysetup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	sqoopOpts, err := sqoopsetup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}

	ctx := contextutils.WithLogger(context.Background(), "apiserver")

	contextutils.LoggerFrom(ctx).Infof("listening on :%v", *port)
	if err := setup.Setup(*port, *dev, glooOpts, gatewayOpts, sqoopOpts); err != nil {
		return err
	}
	return nil
}
