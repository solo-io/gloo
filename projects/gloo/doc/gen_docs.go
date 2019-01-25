package main

import (
	"log"

	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	app := cmd.App(version.Version)
	disableAutoGenTag(app)
	err := doc.GenMarkdownTree(app, "./docs/cli")
	if err != nil {
		log.Fatal(err)
	}
}

func disableAutoGenTag(c *cobra.Command) {
	c.DisableAutoGenTag = true
	for _, c := range c.Commands() {
		disableAutoGenTag(c)
	}
}
