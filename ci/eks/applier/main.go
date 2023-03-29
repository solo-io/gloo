package main

import (
	"github.com/solo-io/solo-projects/ci/eks/applier/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	cmd.Execute()
}
