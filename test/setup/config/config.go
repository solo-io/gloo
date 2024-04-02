package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/test/testutils"
	"sigs.k8s.io/yaml"

	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/gloo/test/setup/types"
)

// Load a configuration from a string.
func Load(ctx context.Context, configPath string, reader io.Reader) (*types.Config, error) {
	var (
		err    error
		config types.Config
	)

	if os.Getenv(testutils.ClusterName) == "" {
		os.Setenv(testutils.ClusterName, "kind")
	}

	if os.Getenv(testutils.Version) == "" {
		// if version is not set, default to the dev version
		os.Setenv(testutils.Version, "1.0.1-dev")
	}

	// Can set defaults here for the process if ran directly from the command line.
	//
	// These will likely be set in make targets, but can't set defaults in the config
	// file itself.
	if os.Getenv(testutils.IstioVersion) == "" {
		latestVersion := strings.Split(kubeutils.DefaultIstioTag, "-")
		if len(latestVersion) == 0 {
			panic("invalid default istio tag setting")
		}

		os.Setenv(testutils.IstioVersion, latestVersion[0])
	}

	if os.Getenv(testutils.IstioHub) == "" {
		os.Setenv(testutils.IstioHub, kubeutils.DefaultIstioImageRegistry)
	}

	if os.Getenv(testutils.ImageRegistry) == "" {
		os.Setenv(testutils.ImageRegistry, kubeutils.DefaultGlooImageRegistry)
	}

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
