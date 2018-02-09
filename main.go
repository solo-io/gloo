package main

import (
	"github.com/solo-io/glue-discovery/cmd"
)

func main() {
	root := cmd.RootCmd()
	root.Execute()
}
