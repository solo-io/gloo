package convert

import (
	"fmt"
	"os"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/convert/domain"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/pflag"
)

type ErrorType string

const (
	ERROR_TYPE_UPDATE_OBJECT             ErrorType = "UPDATE_OBJECT"
	ERROR_TYPE_NOT_SUPPORTED                       = "NOT_SUPPORTED"
	ERROR_TYPE_IGNORED                             = "IGNORED"
	ERROR_TYPE_UNKNOWN_REFERENCE                   = "UNKNOWN_REFERENCE"
	ERROR_TYPE_NO_REFERENCES                       = "NO_REFERENCES"
	ERROR_TYPE_CEL_VALIDATION_CORRECTION           = "CEL_VALIDATION_CORRECTION"
)

type Options struct {
	*options.Options
	InputFile               string
	InputDir                string
	GlooSnapshotFile        string
	OutputDir               string
	Stats                   bool
	CombineRouteOptions     bool
	RetainFolderStructure   bool
	IncludeUnknownResources bool
	DeleteOutputDir         bool
	CreateNamespaces        bool
	ControlPlaneName        string
	ControlPlaneNamespace   string
}

func (opts *Options) validate() error {

	count := 0
	if opts.InputDir != "" {
		count++
	}
	if opts.InputFile != "" {
		count++
	}
	if opts.GlooSnapshotFile != "" {
		count++
	}
	if opts.ControlPlaneName != "" {
		count++
		if opts.ControlPlaneNamespace == "" {
			return fmt.Errorf("pod namespace must be specified")
		}
	}

	if count > 1 {
		return fmt.Errorf("only one of 'input-file' or 'directory' or 'input-snapshot' or `gloo-pod-name` can be specified")
	}
	if !opts.DeleteOutputDir && folderExists(opts.OutputDir) {
		return fmt.Errorf("output-dir already %s exists. It can be deleted with --delete-output-dir", opts.OutputDir)
	}
	return nil
}
func folderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
func (o *Options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.ControlPlaneName, "gloo-control-plane", "", "Name of the Gloo control plane pod")
	flags.StringVarP(&o.ControlPlaneNamespace, "gloo-control-plane-namespace", "n", "gloo-system", "Namespace of the Gloo control plane pod")
	flags.StringVar(&o.InputFile, "input-file", "", "Convert a single YAML file to the Gateway API")
	flags.StringVar(&o.InputDir, "input-dir", "", "InputDir to read yaml/yml files recursively")
	flags.StringVar(&o.GlooSnapshotFile, "input-snapshot", "", "Gloo input snapshot file location")
	flags.BoolVar(&o.Stats, "print-stats", false, "Print stats about the conversion")
	flags.BoolVar(&o.CombineRouteOptions, "combine-route-options", false, "Combine RouteOptions that are exactly the same and share them among the HTTPRoutes")
	flags.StringVar(&o.OutputDir, "output-dir", "./_output", "Output directory to write Gateway API configurations to. The directory must not exist before starting the migration. To delete and recreate the output directory, use the --recreate-output-dir option")
	flags.BoolVar(&o.RetainFolderStructure, "retain-input-folder-structure", false, "Arrange the generated Gateway API files in the same folder structure they were read from (input-dir only).")
	flags.BoolVar(&o.IncludeUnknownResources, "include-unknown", false, "Copy non-Gloo Gateway resources to the output directory without changing them. ")
	flags.BoolVar(&o.DeleteOutputDir, "delete-output-dir", false, "Delete the output directory if it already exists.")
	flags.BoolVar(&o.CreateNamespaces, "create-namespaces", false, "Create namespaces for the objects in a file.")
}

type GatewayAPIOutput struct {
	gatewayAPICache *domain.GatewayAPICache
	edgeCache       *domain.GlooEdgeCache
	errors          map[ErrorType][]GlooError
}

type GlooError struct {
	err       error
	errorType ErrorType
	name      string
	namespace string
	crdType   string
}

func (o *GatewayAPIOutput) AddError(errType ErrorType, msg string, args ...interface{}) {
	if o.errors == nil {
		o.errors = make(map[ErrorType][]GlooError)
	}
	if o.errors[errType] == nil {
		o.errors[errType] = make([]GlooError, 0)
	}

	o.errors[errType] = append(o.errors[errType], GlooError{
		err:       fmt.Errorf(msg, args...),
		errorType: errType,
		name:      "none",
		namespace: "none",
		crdType:   "none",
	})
}
func (o *GatewayAPIOutput) AddErrorFromWrapper(errType ErrorType, wrapper domain.Wrapper, msg string, args ...interface{}) {
	if o.errors == nil {
		o.errors = make(map[ErrorType][]GlooError)
	}
	if o.errors[errType] == nil {
		o.errors[errType] = make([]GlooError, 0)
	}
	o.errors[errType] = append(o.errors[errType], GlooError{
		err:       fmt.Errorf(msg, args...),
		errorType: errType,
		name:      wrapper.GetName(),
		namespace: wrapper.GetNamespace(),
		crdType:   fmt.Sprintf("%s/%s", wrapper.GetObjectKind().GroupVersionKind().Group, wrapper.GetObjectKind().GroupVersionKind().Kind),
	})
}

func NewGatewayAPIOutput() *GatewayAPIOutput {
	return &GatewayAPIOutput{
		gatewayAPICache: &domain.GatewayAPICache{},
		edgeCache:       &domain.GlooEdgeCache{},
	}
}
