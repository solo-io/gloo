package main

import (
	"github.com/solo-io/gloo-function-discovery/cmd"
)

func main() {
	root := cmd.RootCmd()
	root.Execute()
}
