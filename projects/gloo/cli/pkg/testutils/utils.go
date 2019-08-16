package testutils

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/gogo/protobuf/types"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd"
)

func GlooctlEE(args string) error {
	app := cmd.App("test")
	app.SetArgs(strings.Split(args, " "))
	return app.Execute()
}

func GlooctlEEOut(args string) (string, error) {
	stdOut := os.Stdout
	stdErr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w
	os.Stderr = w

	app := cmd.App("test")
	app.SetArgs(strings.Split(args, " "))
	err = app.Execute()

	outC := make(chan string)

	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = stdOut // restoring the real stdout
	os.Stderr = stdErr
	out := <-outC

	return strings.TrimSuffix(out, "\n"), err
}

func GetTestSettings() *gloov1.Settings {
	return &gloov1.Settings{
		Metadata: core.Metadata{
			Name:      "default",
			Namespace: defaults.GlooSystem,
		},
		BindAddr:        "test:80",
		ConfigSource:    &gloov1.Settings_DirectoryConfigSource{},
		DevMode:         true,
		SecretSource:    &gloov1.Settings_KubernetesSecretSource{},
		WatchNamespaces: []string{"default"},
		Extensions: &gloov1.Extensions{
			Configs: map[string]*types.Struct{
				"someotherextension": {},
			},
		},
	}
}
