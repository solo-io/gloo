package generate

import (
	"os"

	"github.com/solo-io/gloo/install/helm/gloo/generate"

	flag "github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/k8s-utils/installutils/helmchart"
)

const (
	devPullPolicy          = string(v1.PullAlways)
	distributionPullPolicy = string(v1.PullIfNotPresent)
)

// We produce two helm artifacts: GlooE and Gloo with a read-only version of the GlooE UI

// these arguments are provided on the command line during make or CI. They are shared by all artifacts
type GenerationArguments struct {
	Version            string
	RepoPrefixOverride string
	// Allows for overriding the gloo-fed chart repo; used in local builds to specify a
	// local directory instead of the official gloo-fed-helm release repository.
	GlooFedRepoOverride string
	GenerateHelmDocs    bool
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
	DocsOutput           string
	ChartTemplate        string
	ChartOutput          string
	RequirementsTemplate string
	RequirementsOutput   string
}

type Artifact int

const (
	GlooE Artifact = iota
)

func ArtifactName(artifact Artifact) string {
	switch artifact {
	case GlooE:
		return "GlooE"
	default:
		return "unknown artifact"
	}
}

// Run generates the helm artifacts for the corresponding file sets
func Run(args *GenerationArguments, fileSets ...*GenerationFiles) error {
	osGlooVersion, err := GetGlooOsVersion(fileSets...)
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
	}

	// Parse optional arguments
	var repoPrefixOverride = flag.String(
		"repo-prefix-override",
		"",
		"(Optional) repository prefix override.")
	var glooFedRepoOverride = flag.String(
		"gloo-fed-repo-override",
		"",
		"(Optional) repository override for gloo-fed chart.")
	var generateHelmDocs = flag.Bool(
		"generate-helm-docs",
		false,
		"(Optional) if set, will generate docs for the helm values")
	flag.Parse()

	if *repoPrefixOverride != "" {
		args.RepoPrefixOverride = *repoPrefixOverride
	}
	if *glooFedRepoOverride != "" {
		args.GlooFedRepoOverride = *glooFedRepoOverride
	}
	if *generateHelmDocs == true {
		args.GenerateHelmDocs = *generateHelmDocs
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
	if err := generateRequirementsYaml(
		gc.GenerationFiles.RequirementsTemplate,
		gc.GenerationFiles.RequirementsOutput,
		gc.OsGlooVersion,
		gc.Arguments.Version,
		gc.Arguments.GlooFedRepoOverride,
	); err != nil {
		return errors.Wrapf(err, "unable to parse requirements.yaml")
	}

	if gc.Arguments.GenerateHelmDocs {
		log.Printf("Generating helm value docs in file: %v", gc.GenerationFiles.DocsOutput)
		if err := gc.generateValueDocs(); err != nil {
			return errors.Wrapf(err, "Generating values.txt failed")
		}
	}
	return nil
}

func (gc *GenerationConfig) generateValuesYamls() error {
	switch gc.GenerationFiles.Artifact {
	case GlooE:
		return gc.generateValuesYamlForGlooE()
	default:
		return errors.New("unknown artifact specified")
	}
}

////////////////////////////////////////////////////////////////////////////////
// generate Gloo-ee values file
////////////////////////////////////////////////////////////////////////////////

func (gc *GenerationConfig) generateValuesConfig(versionOverride string) (*HelmConfig, error) {
	config, err := readConfig(gc.GenerationFiles.ValuesTemplate)
	if err != nil {
		return nil, err
	}

	version := &gc.Arguments.Version
	if versionOverride != "" {
		version = &versionOverride
	}
	tag := &gc.OsGlooVersion
	if tag == nil {
		tag = version
	}
	config.Gloo.Gloo.Deployment.Image.Tag = version
	for _, v := range config.Gloo.GatewayProxies {
		v.PodTemplate.Image.Tag = version
	}
	if config.Gloo.IngressProxy != nil {
		config.Gloo.IngressProxy.Deployment.Image.Tag = version
	}
	config.Gloo.Settings.Integrations.Knative.Proxy.Image.Tag = version
	// Use open source gloo version for discovery and gateway

	// This code used to assume that all relavant structs were already instantiated.
	// But since we no longer duplicate certain most values between the OS and enterprise
	// values-template.yaml file, we need to nil check and create several values that
	// are no longer present in the default enterprise values-template.
	if config.Gloo.Discovery == nil {
		config.Gloo.Discovery = &generate.Discovery{}
	}
	if config.Gloo.Discovery.Deployment == nil {
		config.Gloo.Discovery.Deployment = &generate.DiscoveryDeployment{}
	}
	config.Gloo.Discovery.Deployment.Image.Tag = tag

	if config.Gloo.Gateway == nil {
		config.Gloo.Gateway = &generate.Gateway{}
	}
	if config.Gloo.Gateway.Deployment == nil {
		config.Gloo.Gateway.Deployment = &generate.GatewayDeployment{}
	}
	if config.Gloo.Gateway.Deployment.Image == nil {
		config.Gloo.Gateway.Deployment.Image = &generate.Image{}
	}
	config.Gloo.Gateway.Deployment.Image.Tag = tag

	if config.Gloo.Gateway.CertGenJob == nil {
		config.Gloo.Gateway.CertGenJob = &generate.CertGenJob{}
	}
	if config.Gloo.Gateway.CertGenJob.Image == nil {
		config.Gloo.Gateway.CertGenJob.Image = &generate.Image{}
	}
	config.Gloo.Gateway.CertGenJob.Image.Tag = tag

	config.Observability.Deployment.Image.Tag = version

	if config.Global.GlooMtls.Sds.Image == nil {
		config.Global.GlooMtls.Sds.Image = &generate.Image{}
	}
	config.Global.GlooMtls.Sds.Image.Tag = tag
	config.Global.GlooMtls.EnvoySidecar.Image.Tag = version

	pullPolicy := gc.PullPolicyForVersion
	config.Gloo.Gloo.Deployment.Image.PullPolicy = &pullPolicy
	for _, v := range config.Gloo.GatewayProxies {
		v.PodTemplate.Image.PullPolicy = &pullPolicy
	}
	if config.Gloo.IngressProxy != nil {
		config.Gloo.IngressProxy.Deployment.Image.PullPolicy = &pullPolicy
	}

	config.Gloo.Settings.Integrations.Knative.Proxy.Image.PullPolicy = &pullPolicy
	config.Gloo.Discovery.Deployment.Image.PullPolicy = &pullPolicy
	config.Gloo.Gateway.Deployment.Image.PullPolicy = &pullPolicy
	config.Gloo.Gateway.CertGenJob.Image.PullPolicy = &pullPolicy
	config.Observability.Deployment.Image.PullPolicy = &pullPolicy
	config.Redis.Deployment.Image.PullPolicy = &pullPolicy

	if err = updateExtensionsImageVersionAndPullPolicy(config, pullPolicy, version); err != nil {
		return nil, err
	}

	if gc.Arguments.RepoPrefixOverride != "" {
		config.Global.Image.Registry = &gc.Arguments.RepoPrefixOverride
	}
	return &config, nil
}

func (gc *GenerationConfig) generateValuesYamlForGlooE() error {
	config, err := gc.generateValuesConfig("")
	if err != nil {
		return errors.Wrapf(err, "Unable to generate values config")
	}

	if err := writeYaml(config, gc.GenerationFiles.ValuesOutput); err != nil {
		return errors.Wrapf(err, "unable to generate GlooE")
	}
	return nil
}

func (gc *GenerationConfig) generateValueDocs() error {
	// Overwrite the literal version with a description of the field value
	config, err := gc.generateValuesConfig("Version number, ex. 1.8.0")
	if err != nil {
		return errors.Wrapf(err, "Unable to generate values config")
	}
	return writeDocs(helmchart.Doc(config), gc.GenerationFiles.DocsOutput)
}

func updateExtensionsImageVersionAndPullPolicy(config HelmConfig, pullPolicy string, version *string) (err error) {
	bytes, err := yaml.Marshal(config.Global.Extensions)
	if err != nil {
		return err
	}
	var glooEeExtensions GlooEeExtensions
	err = yaml.Unmarshal(bytes, &glooEeExtensions)
	if err != nil {
		return err
	}
	// Extauth and rate-limit are both referenced in Values.gloo.settings, and thus need to be retro-typed
	// to avoid type-leakage into gloo-OS. Because helm like re-typing values defined in imported charts,
	// we must also place these in the shared `.Values.global.` struct.
	// The following code simply applies the version/pull policy cohesion that generateValuesYamlForGlooE() does
	// for everything else.

	glooEeExtensions.ExtAuth.Deployment.Image.Tag = version
	glooEeExtensions.ExtAuth.Deployment.Image.PullPolicy = &pullPolicy

	glooEeExtensions.RateLimit.Deployment.Image.Tag = version
	glooEeExtensions.RateLimit.Deployment.Image.PullPolicy = &pullPolicy

	config.Global.Extensions = glooEeExtensions
	return nil
}
