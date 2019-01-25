package testutils

import (
	"strings"

	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd"
)

func GlooctlEE(args string) error {
	app := cmd.App("test")
	app.SetArgs(strings.Split(args, " "))
	return app.Execute()
}
