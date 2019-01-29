package testutils

import (
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
				"someotherextension": &types.Struct{},
			},
		},
	}
}
