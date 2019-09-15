package generate

import (
	"os"

	v1 "k8s.io/api/core/v1"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/log"
)

const (
	gopkgToml    = "Gopkg.toml"
	constraint   = "constraint"
	glooPkg      = "github.com/solo-io/gloo"
	nameConst    = "name"
	versionConst = "version"

	devPullPolicy          = string(v1.PullAlways)
	distributionPullPolicy = string(v1.PullIfNotPresent)
)

// We produce two helm artifacts: GlooE and Gloo with a read-only version of the GlooE UI

// these arguments are provided on the command line during make or CI. They are shared by all artifacts
type GenerationArguments struct {
	Version            string
	RepoPrefixOverride string
}

// GenerationConfig represents all the artifact-specific config
type GenerationConfig struct {
	Arguments            *GenerationArguments
	OsGlooVersion        string
	GenerationFiles      *GenerationFiles
	PullPolicyForVersion string
}

// GenerationFiles specify the files read or created while producing a given artifact
type GenerationFiles struct {
	// not a file, just a way of identifying the purpose of the fileset
	Artifact             Artifact
	ValuesTemplate       string
	ValuesOutput         string
	DistributionOutput   string
	ChartTemplate        string
	ChartOutput          string
	RequirementsTemplate string
	RequirementsOutput   string
}

type Artifact int

const (
	GlooE Artifact = iota
	GlooWithRoUi
)

func ArtifactName(artifact Artifact) string {
	switch artifact {
	case GlooE:
		return "GlooE"
	case GlooWithRoUi:
		return "Gloo OS with read-only UI"
	default:
		return "unknown artifact"
	}
}

// Run generates the helm artifacts for the corresponding file sets
func Run(args *GenerationArguments, fileSets ...*GenerationFiles) error {
	osGlooVersion, err := GetVersionFromToml(gopkgToml, glooPkg)
	if err != nil {
		return errors.Wrapf(err, "failed to determine open source Gloo version")
	}

	for _, fileSet := range fileSets {
		genConfig := GetGenerationConfig(args, osGlooVersion, fileSet)

		if err := genConfig.runGeneration(); err != nil {
			return errors.Wrapf(err, "unable to Run generation for glooE")
		}
	}
	return nil
}

func GetGenerationConfig(args *GenerationArguments, osGlooVersion string, generationFiles *GenerationFiles) *GenerationConfig {
	pullPolicyForVersion := distributionPullPolicy
	if args.Version == "dev" {
		pullPolicyForVersion = devPullPolicy
	}
	return &GenerationConfig{
		Arguments:            args,
		OsGlooVersion:        osGlooVersion,
		PullPolicyForVersion: pullPolicyForVersion,
		GenerationFiles:      generationFiles,
	}
}

func GetArguments(args *GenerationArguments) error {
	if len(os.Args) < 2 {
		return errors.New("Must provide version as argument")
	} else {
		args.Version = os.Args[1]

		if len(os.Args) == 3 {
			args.RepoPrefixOverride = os.Args[2]
		}
	}
	return nil
}

func (gc *GenerationConfig) runGeneration() error {
	osGlooVersion, err := GetVersionFromToml(gopkgToml, glooPkg)
	if err != nil {
		log.Fatalf("failed to determine open source Gloo version. Cause: %v", err)
	}
	log.Printf("Open source gloo version is: %v", osGlooVersion)

	log.Printf("Generating helm files.")
	if err := gc.generateValuesYamls(); err != nil {
		return errors.Wrapf(err, "generating values.yaml failed")
	}
	if err := gc.generateChartYaml(gc.GenerationFiles.ChartTemplate, gc.GenerationFiles.ChartOutput, gc.Arguments.Version); err != nil {
		return errors.Wrapf(err, "generating Chart.yaml failed")
	}
	if err := generateRequirementsYaml(gc.GenerationFiles.RequirementsTemplate, gc.GenerationFiles.RequirementsOutput, gc.OsGlooVersion); err != nil {
		return errors.Wrapf(err, "unable to parse Gopkg.toml for proper gloo version")
	}
	return nil
}

func (gc *GenerationConfig) generateValuesYamls() error {
	switch gc.GenerationFiles.Artifact {
	case GlooE:
		return gc.generateValuesYamlForGlooE()
	case GlooWithRoUi:
		return gc.generateValuesYamlForGlooOsWithRoUi()
	default:
		return errors.New("unknown artifact specified")
	}
}

func (gc *GenerationConfig) generateValuesYamlForGlooE() error {
	// generate two forms of artifacts for GlooE // TODO- document why we are doing this
	type glooEArtifactSpecialization struct {
		outputFile       string
		repositoryPrefix string
	}
	specializations := []glooEArtifactSpecialization{{
		outputFile:       gc.GenerationFiles.ValuesOutput,
		repositoryPrefix: gc.Arguments.RepoPrefixOverride,
	}, {
		outputFile:       gc.GenerationFiles.DistributionOutput,
		repositoryPrefix: "",
	}}
	for i, specialization := range specializations {
		outputFile := specialization.outputFile
		repositoryPrefix := specialization.repositoryPrefix
		version := gc.Arguments.Version
		pullPolicy := gc.PullPolicyForVersion
		config, err := readConfig(gc.GenerationFiles.ValuesTemplate)
		if err != nil {
			return err
		}

		config.Gloo.Gloo.Deployment.Image.Tag = version
		for _, v := range config.Gloo.GatewayProxies {
			v.PodTemplate.Image.Tag = version
		}
		if config.Gloo.IngressProxy != nil {
			config.Gloo.IngressProxy.Deployment.Image.Tag = version
		}
		// Use open source gloo version for discovery and gateway
		config.Gloo.Discovery.Deployment.Image.Tag = gc.OsGlooVersion
		config.Gloo.Gateway.Deployment.Image.Tag = gc.OsGlooVersion
		config.RateLimit.Deployment.Image.Tag = version
		config.Observability.Deployment.Image.Tag = version
		config.ApiServer.Deployment.Server.Image.Tag = version
		config.ApiServer.Deployment.Envoy.Image.Tag = version
		config.ApiServer.Deployment.Ui.Image.Tag = version

		config.Gloo.Gloo.Deployment.Image.PullPolicy = pullPolicy
		for _, v := range config.Gloo.GatewayProxies {
			v.PodTemplate.Image.PullPolicy = pullPolicy
		}
		if config.Gloo.IngressProxy != nil {
			config.Gloo.IngressProxy.Deployment.Image.PullPolicy = pullPolicy
		}
		config.Gloo.Discovery.Deployment.Image.PullPolicy = pullPolicy
		config.Gloo.Gateway.Deployment.Image.PullPolicy = pullPolicy
		config.RateLimit.Deployment.Image.PullPolicy = pullPolicy
		config.Observability.Deployment.Image.PullPolicy = pullPolicy
		config.Redis.Deployment.Image.PullPolicy = pullPolicy
		config.ApiServer.Deployment.Ui.Image.PullPolicy = pullPolicy
		config.ApiServer.Deployment.Server.Image.PullPolicy = pullPolicy
		config.ApiServer.Deployment.Envoy.Image.PullPolicy = pullPolicy

		if err = updateExtensionsImageVersionAndPullPolicy(config, version, pullPolicy); err != nil {
			return err
		}

		if repositoryPrefix != "" {
			config.Global.Image.Registry = repositoryPrefix
		}

		if err := writeYaml(&config, outputFile); err != nil {
			return errors.Wrapf(err, "unable to generate GlooE specialization %v", i)
		}
	}
	return nil
}

func updateExtensionsImageVersionAndPullPolicy(config HelmConfig, version, pullPolicy string) (err error) {
	bytes, err := yaml.Marshal(config.Global.Extensions)
	if err != nil {
		return err
	}
	var glooEeExtensions GlooEeExtensions
	err = yaml.Unmarshal(bytes, &glooEeExtensions)
	if err != nil {
		return err
	}
	glooEeExtensions.ExtAuth.Deployment.Image.Tag = version
	glooEeExtensions.ExtAuth.Deployment.Image.PullPolicy = pullPolicy
	config.Global.Extensions = glooEeExtensions
	return nil
}

func (gc *GenerationConfig) generateValuesYamlForGlooOsWithRoUi() error {
	outputFile := gc.GenerationFiles.ValuesOutput
	repositoryPrefix := gc.Arguments.RepoPrefixOverride
	version := gc.Arguments.Version
	pullPolicy := gc.PullPolicyForVersion
	config, err := readConfig(gc.GenerationFiles.ValuesTemplate)
	if err != nil {
		return err
	}

	config.Gloo.Gloo.Deployment.Image.Tag = version
	for _, v := range config.Gloo.GatewayProxies {
		v.PodTemplate.Image.Tag = version
	}
	if config.Gloo.IngressProxy != nil {
		config.Gloo.IngressProxy.Deployment.Image.Tag = version
	}
	// Use open source gloo version for discovery and gateway
	config.Gloo.Discovery.Deployment.Image.Tag = gc.OsGlooVersion
	config.Gloo.Gateway.Deployment.Image.Tag = gc.OsGlooVersion
	config.ApiServer.Deployment.Server.Image.Tag = version
	config.ApiServer.Deployment.Envoy.Image.Tag = version
	config.ApiServer.Deployment.Ui.Image.Tag = version

	config.Gloo.Gloo.Deployment.Image.PullPolicy = pullPolicy
	for _, v := range config.Gloo.GatewayProxies {
		v.PodTemplate.Image.PullPolicy = pullPolicy
	}
	if config.Gloo.IngressProxy != nil {
		config.Gloo.IngressProxy.Deployment.Image.PullPolicy = pullPolicy
	}
	config.Gloo.Discovery.Deployment.Image.PullPolicy = pullPolicy
	config.Gloo.Gateway.Deployment.Image.PullPolicy = pullPolicy
	config.ApiServer.Deployment.Ui.Image.PullPolicy = pullPolicy
	config.ApiServer.Deployment.Server.Image.PullPolicy = pullPolicy
	config.ApiServer.Deployment.Envoy.Image.PullPolicy = pullPolicy

	if repositoryPrefix != "" {
		config.Global.Image.Registry = repositoryPrefix
	}

	return writeYaml(&config, outputFile)
}
