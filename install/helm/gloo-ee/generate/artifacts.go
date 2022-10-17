package generate

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

	"github.com/solo-io/gloo/install/helm/gloo/generate"

	flag "github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/k8s-utils/installutils/helmchart"
)

const (
	devPullPolicy          = string(v1.PullAlways)
	distributionPullPolicy = string(v1.PullIfNotPresent)
	// this is the default arm image repo for the arm work around
	defaultArmImageRegistry = "localhost:5000"
	// this is the default image registry
	defaultImageRegistry = "quay.io/solo-io"
)

// We produce two helm artifacts: GlooE and Gloo with a read-only version of the GlooE UI

// these arguments are provided on the command line during make or CI. They are shared by all artifacts
type GenerationArguments struct {
	Version string
	// Allows for overriding the global image registry
	RepoPrefixOverride string
	// Allows for overriding the gloo-fed chart repo; used in local builds to specify a
	// local directory instead of the official gloo-fed-helm release repository.
	GlooFedRepoOverride string
	GenerateHelmDocs    bool

	// specify using image digests, rather than image tags
	UseDigests bool
}

// GenerationConfig represents all the artifact-specific config
type GenerationConfig struct {
	Arguments       *GenerationArguments
	OsGlooVersion   string
	GenerationFiles *GenerationFiles
	// override the global image Pull Policy of the helm chart
	PullPolicyForVersion string
	// if running arm, and not using Generation Arguments RepoPrefixOverride,
	// this will set the registry for all images if the architecture is arm
	ArmImageRegistry string
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
	// when running in arm64, use work around settings
	imageRepo := args.RepoPrefixOverride
	if runtime.GOARCH == "arm64" && imageRepo == "" && os.Getenv("RUNNING_REGRESSION_TESTS") == "true" {
		imageRepo = os.Getenv("IMAGE_REPO")
		if imageRepo == "" {
			imageRepo = defaultArmImageRegistry
		}
		pullPolicyForVersion = devPullPolicy
	}
	return &GenerationConfig{
		Arguments:            args,
		OsGlooVersion:        osGlooVersion,
		PullPolicyForVersion: pullPolicyForVersion,
		GenerationFiles:      generationFiles,
		ArmImageRegistry:     imageRepo,
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
	var useDigests = flag.Bool(
		"use-digests",
		false,
		"(Optional) specify using image digests, rather than image tags",
	)
	flag.Parse()

	if *repoPrefixOverride != "" {
		args.RepoPrefixOverride = *repoPrefixOverride
	}
	if *glooFedRepoOverride != "" {
		args.GlooFedRepoOverride = *glooFedRepoOverride
	}
	if *generateHelmDocs {
		args.GenerateHelmDocs = *generateHelmDocs
	}
	if *useDigests {
		args.UseDigests = *useDigests
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
	config.Gloo.Gloo.Deployment.OssImageTag = tag
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
	config.Gloo.Discovery.Deployment.Image.Tag = version

	if config.Gloo.Gateway == nil {
		config.Gloo.Gateway = &generate.Gateway{}
	}

	if config.Gloo.Gateway.CertGenJob == nil {
		config.Gloo.Gateway.CertGenJob = &generate.CertGenJob{}
	}
	if config.Gloo.Gateway.CertGenJob.Image == nil {
		config.Gloo.Gateway.CertGenJob.Image = &generate.Image{}
	}
	config.Gloo.Gateway.CertGenJob.Image.Tag = tag

	if config.Gloo.Gateway.RolloutJob == nil {
		config.Gloo.Gateway.RolloutJob = &generate.RolloutJob{}
	}
	if config.Gloo.Gateway.RolloutJob.Image == nil {
		config.Gloo.Gateway.RolloutJob.Image = &generate.Image{}
	}
	config.Gloo.Gateway.RolloutJob.Image.Tag = tag

	if config.Gloo.Gateway.CleanupJob == nil {
		config.Gloo.Gateway.CleanupJob = &generate.CleanupJob{}
	}
	if config.Gloo.Gateway.CleanupJob.Image == nil {
		config.Gloo.Gateway.CleanupJob.Image = &generate.Image{}
	}
	config.Gloo.Gateway.CleanupJob.Image.Tag = tag

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
	config.Gloo.Gateway.CertGenJob.Image.PullPolicy = &pullPolicy
	config.Gloo.Gateway.RolloutJob.Image.PullPolicy = &pullPolicy
	config.Gloo.Gateway.CleanupJob.Image.PullPolicy = &pullPolicy
	config.Observability.Deployment.Image.PullPolicy = &pullPolicy
	config.Redis.Deployment.Image.PullPolicy = &pullPolicy

	updateExtensionsImageVersionAndPullPolicy(config, pullPolicy, version)

	if gc.Arguments.RepoPrefixOverride != "" {
		config.Global.Image.Registry = &gc.Arguments.RepoPrefixOverride
	} else if gc.ArmImageRegistry != "" {
		// assumes that if running arm, client wants registry to be default registry AKA (quay.io)
		// for the below images
		config.Global.Image.Registry = &gc.ArmImageRegistry
		gtw := config.Gloo.Gateway
		dir := defaultImageRegistry
		gtw.CertGenJob.Image.Registry = &dir
		gtw.RolloutJob.Image.Registry = &dir
		gtw.CleanupJob.Image.Registry = &dir
		config.Global.GlooMtls.Sds.Image.Registry = &dir
	}

	if gc.Arguments.UseDigests {
		return gc.convertAllTagsToDigests(config), nil
	}
	return &config, nil
}

func (gc *GenerationConfig) convertAllTagsToDigests(config HelmConfig) *HelmConfig {
	// "intelligently" process all nested structs of structs
	for _, imagePath := range gc.findImagePaths(config) {
		// extract image object coresponding to imagePath
		if !canUseFieldByIndex(reflect.ValueOf(config), imagePath) {
			continue // fields _above_ "Image" are nil
		}

		img := reflect.ValueOf(config).FieldByIndex(imagePath).Interface().(*generate.Image)
		if img == nil || img.Tag == nil || img.Repository == nil {
			continue // "Image" is nil
		}

		setDigest(img, config)
	}

	// handle special cases
	for _, v := range config.Gloo.GatewayProxies { // <--- "list of structs" != "struct of struct" assumptions
		setDigest(v.PodTemplate.Image, config)
	}

	return &config
}

func (gc *GenerationConfig) findImagePaths(config HelmConfig) [][]int {
	// to start, we have 2 top-level fields to search: Config and GlobalConfig
	fieldPaths := [][]int{{0}, {1}}   // collector object; in-progress field paths stored here
	imagePaths := [][]int{}           // collector object; matched "Image" field paths stored here
	var curPath []int                 // currently-considered field path.  Ex:  {0,0,0} --> config.Config.Settings.WatchNamespaces
	tConfig := reflect.TypeOf(config) // top level config type
	var cConfig reflect.Type          // current field type

	// DFS for all paths with an "Image" field
	for len(fieldPaths) > 0 {
		// pop oldest-added subpath
		curPath, fieldPaths = fieldPaths[0], fieldPaths[1:]
		cConfig = tConfig.FieldByIndex(curPath).Type
		if cConfig.Kind() == reflect.Ptr {
			cConfig = cConfig.Elem()
		}

		if cConfig.Name() == "Image" {
			// found an image to overwrite
			imagePaths = append(imagePaths, curPath)
			continue
		} else if cConfig.Kind() != reflect.Struct {
			// We have delved too greedily and too deep.  Turn back
			continue
		}

		// add newly-discovered fields to collector
		for i := 0; i < cConfig.NumField(); i++ {
			tmp := make([]int, len(curPath))
			copy(tmp, curPath) // copy op needed to prevent mutation of earlier added items to fieldPaths
			fieldPaths = append(fieldPaths, append(tmp, i))
		}
	}
	return imagePaths
}

func setDigest(img *generate.Image, config HelmConfig) {
	registry := config.Global.Image.Registry
	if img.Registry != nil {
		registry = img.Registry
	}
	imageUrl := *registry + "/" + *img.Repository + ":" + *img.Tag

	digest, _, _ := shellout("docker manifest inspect " + imageUrl + " -v | jq -r \".Descriptor.digest\"")
	digest = strings.TrimSpace(digest)
	if digest != "" {
		// notably, non-solo-produced images are ignored, here.  Though not exhaustively true, many of them
		// don't have a _single_ digest.  Rather, they have several, on a per-platform basis.
		img.Digest = &digest
	}

	// FIPS exceptions pulled from "gloo.image": https://github.com/solo-io/gloo/blob/master/install/helm/gloo/templates/_helpers.tpl#L22-L24
	if *img.Repository == "gloo-ee" || *img.Repository == "extauth-ee" ||
		*img.Repository == "gloo-ee-envoy-wrapper" || *img.Repository == "rate-limit-ee" || *img.Repository == "ext-auth-plugins" {

		fipsUrl := *registry + "/" + *img.Repository + "-fips:" + *img.Tag
		digest, _, _ := shellout("docker manifest inspect " + fipsUrl + " -v | jq -r \".Descriptor.digest\"")
		digest = strings.TrimSpace(digest)
		if digest != "" {
			// notably, non-solo-produced images are ignored, here.  Though not exhaustively true, many of them
			// don't have a _single_ digest.  Rather, they have several, on a per-platform basis.
			img.FipsDigest = &digest
		}
	}
}

func shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func canUseFieldByIndex(v reflect.Value, index []int) bool {
	// stolen/lightly modified from reflect.FieldByIndex.  For our use case, we don't want to panic
	// on a v.IsNil().  Rather, we want to abort our computation
	if len(index) == 1 {
		return true
	}
	for i, x := range index {
		if i > 0 {
			if v.Kind() == reflect.Pointer {
				if v.IsNil() {
					return false
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return true
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

func updateExtensionsImageVersionAndPullPolicy(config HelmConfig, pullPolicy string, version *string) {
	// Extauth and rate-limit are both referenced in Values.gloo.settings, and thus need to be retro-typed
	// to avoid type-leakage into gloo-OS. Because helm like re-typing values defined in imported charts,
	// we must also place these in the shared `.Values.global.` struct.
	// The following code simply applies the version/pull policy cohesion that generateValuesYamlForGlooE() does
	// for everything else.

	config.Global.Extensions.ExtAuth.Deployment.Image.Tag = version
	config.Global.Extensions.ExtAuth.Deployment.Image.PullPolicy = &pullPolicy

	config.Global.Extensions.RateLimit.Deployment.Image.Tag = version
	config.Global.Extensions.RateLimit.Deployment.Image.PullPolicy = &pullPolicy

	config.Global.Extensions.Caching.Deployment.Image.Tag = version
	config.Global.Extensions.Caching.Deployment.Image.PullPolicy = &pullPolicy
}
