package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/plugins"
	"go.uber.org/zap"
)

var pluginManifestFlagUsage = `A .yaml file containing information required to load the ext auth plugins. Must have the following format: 

	plugins:
	- name: MyPlugin
	  pluginFileName: Plugin.so
	  exportedSymbolName: MyPlugin
	- name: AnotherPlugin
	  plugin_file_name: AnotherFile.so
	  exported_symbol_name: AnotherSymbol
`

func main() {
	ctx := contextutils.WithLogger(context.Background(), "verify-plugins")
	logger := contextutils.LoggerFrom(ctx)

	// Get command line flags
	pluginDir := flag.String("pluginDir", "", "directory containing the compiled plugin files")
	pluginManifestFile := flag.String("manifest", "", pluginManifestFlagUsage)
	debug := flag.Bool("debug", false, "if present sets log level to debug")
	flag.Parse()

	if *pluginManifestFile == "" {
		fmt.Println("Missing required -f option")
		flag.Usage()
		os.Exit(1)
	}

	if *pluginDir == "" {
		fmt.Println("Missing required -pluginDir option")
		flag.Usage()
		os.Exit(1)
	}

	if *debug {
		contextutils.SetLogLevel(zap.DebugLevel)
	}

	if err := verifyPlugins(ctx, *pluginDir, *pluginManifestFile); err != nil {
		logger.Errorw("Plugin(s) cannot be loaded by Gloo", zap.Any("error", err))
		os.Exit(1)
	}

	logger.Info("Successfully verified that plugins can be loaded by Gloo!")
	os.Exit(0)
}

func verifyPlugins(ctx context.Context, pluginDir, pluginManifestFile string) error {
	pluginConfig, err := parseManifestFile(pluginManifestFile)
	if err != nil {
		return errors.Wrapf(err, "failed to parse plugin manifest file")
	}

	return loadPlugins(ctx, pluginDir, pluginConfig)
}

func parseManifestFile(filePath string) (*extauth.PluginAuth, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	bytes, err = yaml.YAMLToJSON(bytes)
	if err != nil {
		return nil, err
	}

	into := &extauth.PluginAuth{}
	if err = protoutils.UnmarshalBytes(bytes, into); err != nil {
		return nil, err
	}
	return into, nil
}

func loadPlugins(ctx context.Context, pluginDir string, pluginConfig *extauth.PluginAuth) error {
	sanitizeConfig(pluginConfig)
	_, err := plugins.NewPluginLoader(pluginDir).Load(ctx, pluginConfig)
	return err
}

// Loader will fail if proto is nil
func sanitizeConfig(pluginConfig *extauth.PluginAuth) {
	for _, authPlugin := range pluginConfig.Plugins {
		if authPlugin.Config == nil {
			authPlugin.Config = &types.Struct{
				Fields: map[string]*types.Value{},
			}
		}
	}
}
