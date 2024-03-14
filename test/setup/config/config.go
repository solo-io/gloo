package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/gloo/test/setup/defaults"
	"github.com/solo-io/gloo/test/setup/types"
)

// Load a configuration from a string.
func Load(ctx context.Context, configPath string, reader io.Reader) (*types.Config, error) {
	var (
		err    error
		config types.Config
		logger = contextutils.LoggerFrom(ctx)
	)

	if os.Getenv("VERSION") == "" {
		// if version is not set, default to the dev version
		os.Setenv("VERSION", "1.0.1-dev")
	}

	// Can set defaults here for the process if ran directly from the command line.
	//
	// These will likely be set in make targets, but can't set defaults in the config
	// file itself.
	if os.Getenv("ISTIO_VERSION") == "" {
		latestVersion := strings.Split(defaults.DefaultIstioTag, "-")
		if len(latestVersion) == 0 {
			panic("invalid default istio tag setting")
		}

		os.Setenv("ISTIO_VERSION", latestVersion[0])
	}

	if os.Getenv("ISTIO_HUB") == "" {
		os.Setenv("ISTIO_HUB", defaults.DefaultIstioImageRegistry)
	}

	if os.Getenv("IMAGE_REGISTRY") == "" {
		os.Setenv("IMAGE_REGISTRY", defaults.DefaultGlooImageRegistry)
	}

	//if os.Getenv("VERSION") == "" {
	//	os.Setenv("VERSION", helpers.Version())
	//}

	logger.Infof("Using Istio hub: %s version: %s", os.Getenv("ISTIO_HUB"), os.Getenv("ISTIO_VERSION"))

	contents, err := readConfig(ctx, configPath, reader)
	if err != nil {
		return nil, err
	}

	stringWithEnvs := os.Expand(string(contents), func(s string) string {
		// Escape case
		if strings.HasPrefix(s, "$") {
			return s
		}
		// Allow special syntax through e.g: ${env:MY_POD_IP}
		if strings.Contains(s, ":") || strings.Contains(s, ":-") {
			return fmt.Sprintf("${%s}", s)
		}
		return os.Getenv(s)
	})

	if err = yaml.Unmarshal([]byte(stringWithEnvs), &config); err != nil {
		return nil, err
	}

	if len(config.Clusters) == 0 {
		return nil, fmt.Errorf("no clusters defined in config")
	}

	for _, cluster := range config.Clusters {
		for _, chart := range cluster.Charts {
			// Nice to have these values in a struct for things like loading images
			// from source.
			if chart.Name == "gloo" {
				valuesData, err := yaml.Marshal(chart.Values)
				if err != nil {
					return nil, err
				}

				if err = yaml.Unmarshal(valuesData, &cluster.GlooEdge); err != nil {
					return nil, err
				}
			}
		}
	}
	return &config, nil
}

func readConfig(ctx context.Context, configPath string, reader io.Reader) ([]byte, error) {
	defer contextutils.LoggerFrom(ctx).Info("finished reading config")

	if configPath == "-" {
		contextutils.LoggerFrom(ctx).Info("reading config from stdin", configPath)
		return io.ReadAll(reader)
	}
	return os.ReadFile(configPath)
}
