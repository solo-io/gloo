package generate

import (
	"os"

	v1 "k8s.io/api/core/v1"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/log"
)

const (
	pathToGopkgTomlDir = "."

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

	osGlooVersion, err := GetGlooOsVersion(pathToGopkgTomlDir, fileSets...)
	if err != nil {
		return errors.Wrapf(err, "failed to determine open source Gloo version")
	}

	log.Printf("Open source gloo version is: %v", osGlooVersion)

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

////////////////////////////////////////////////////////////////////////////////
// generate Gloo-ee values file
////////////////////////////////////////////////////////////////////////////////

func (gc *GenerationConfig) generateValuesYamlForGlooE() error {
	config, err := readConfig(gc.GenerationFiles.ValuesTemplate)
	if err != nil {
		return err
	}

	version := gc.Arguments.Version
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
	config.Gloo.Gateway.CertGenJob.Image.Tag = gc.OsGlooVersion
	config.RateLimit.Deployment.Image.Tag = version
	config.Observability.Deployment.Image.Tag = version
	config.ApiServer.Deployment.Server.Image.Tag = version
	config.ApiServer.Deployment.Envoy.Image.Tag = version
	config.ApiServer.Deployment.Ui.Image.Tag = version

	pullPolicy := gc.PullPolicyForVersion
	config.Gloo.Gloo.Deployment.Image.PullPolicy = pullPolicy
	for _, v := range config.Gloo.GatewayProxies {
		v.PodTemplate.Image.PullPolicy = pullPolicy
	}
	if config.Gloo.IngressProxy != nil {
		config.Gloo.IngressProxy.Deployment.Image.PullPolicy = pullPolicy
	}
	config.Gloo.Discovery.Deployment.Image.PullPolicy = pullPolicy
	config.Gloo.Gateway.Deployment.Image.PullPolicy = pullPolicy
	config.Gloo.Gateway.CertGenJob.Image.PullPolicy = pullPolicy
	config.RateLimit.Deployment.Image.PullPolicy = pullPolicy
	config.Observability.Deployment.Image.PullPolicy = pullPolicy
	config.Redis.Deployment.Image.PullPolicy = pullPolicy
	config.ApiServer.Deployment.Ui.Image.PullPolicy = pullPolicy
	config.ApiServer.Deployment.Server.Image.PullPolicy = pullPolicy
	config.ApiServer.Deployment.Envoy.Image.PullPolicy = pullPolicy

	if err = updateExtensionsImageVersionAndPullPolicy(config, version, pullPolicy); err != nil {
		return err
	}

	if gc.Arguments.RepoPrefixOverride != "" {
		config.Global.Image.Registry = gc.Arguments.RepoPrefixOverride
	}

	if err := writeYaml(&config, gc.GenerationFiles.ValuesOutput); err != nil {
		return errors.Wrapf(err, "unable to generate GlooE")
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

////////////////////////////////////////////////////////////////////////////////
// generate Gloo-os with read-only ui values file
////////////////////////////////////////////////////////////////////////////////

func (gc *GenerationConfig) generateValuesYamlForGlooOsWithRoUi() error {
	config, err := readConfig(gc.GenerationFiles.ValuesTemplate)
	if err != nil {
		return err
	}

	glooEVersion := gc.Arguments.Version
	for _, v := range config.Gloo.GatewayProxies {
		v.PodTemplate.Image.Tag = gc.OsGlooVersion
	}
	if config.Gloo.IngressProxy != nil {
		config.Gloo.IngressProxy.Deployment.Image.Tag = gc.OsGlooVersion
	}
	// Use open source gloo version for gloo, discovery, and gateway
	config.Gloo.Gloo.Deployment.Image.Tag = gc.OsGlooVersion
	config.Gloo.Discovery.Deployment.Image.Tag = gc.OsGlooVersion
	config.Gloo.Gateway.Deployment.Image.Tag = gc.OsGlooVersion
	config.Gloo.AccessLogger.Image.Tag = gc.OsGlooVersion
	config.ApiServer.Deployment.Server.Image.Tag = glooEVersion
	config.ApiServer.Deployment.Envoy.Image.Tag = glooEVersion
	config.ApiServer.Deployment.Ui.Image.Tag = glooEVersion

	pullPolicy := gc.PullPolicyForVersion
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

	if gc.Arguments.RepoPrefixOverride != "" {
		config.Global.Image.Registry = gc.Arguments.RepoPrefixOverride
	}

	return writeYaml(&config, gc.GenerationFiles.ValuesOutput)
}
