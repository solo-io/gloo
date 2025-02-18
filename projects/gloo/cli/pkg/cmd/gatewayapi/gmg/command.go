package gmg

import (
	"errors"
	"fmt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	admin_v3 "github.com/solo-io/solo-apis/client-go/admin.gloo.solo.io/v2"
	apimanagement_v2 "github.com/solo-io/solo-apis/client-go/apimanagement.gloo.solo.io/v2"
	infrastructure_v3 "github.com/solo-io/solo-apis/client-go/infrastructure.gloo.solo.io/v2"
	networking_v2 "github.com/solo-io/solo-apis/client-go/networking.gloo.solo.io/v2"
	observability_v2 "github.com/solo-io/solo-apis/client-go/observability.policy.gloo.solo.io/v2"
	resilience_v2 "github.com/solo-io/solo-apis/client-go/resilience.policy.gloo.solo.io/v2"
	security_v2 "github.com/solo-io/solo-apis/client-go/security.policy.gloo.solo.io/v2"
	traffic_v2 "github.com/solo-io/solo-apis/client-go/trafficcontrol.policy.gloo.solo.io/v2"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/utils/ptr"
	"log"
	"os"
	"path/filepath"
	"regexp"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/yaml"
	"strings"
)

const (
	RandomSuffix = 4
	RandomSeed   = 1
)

var runtimeScheme *runtime.Scheme
var codecs serializer.CodecFactory
var decoder runtime.Decoder

func RootCmd() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:   "gmg",
		Short: "Convert Gloo Mesh APIs to Gateway API",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}
	opts.addToFlags(cmd.PersistentFlags())
	cmd.SilenceUsage = true
	return cmd
}

func run(opts *Options) error {

	foundFiles, err := findFiles(opts)
	if err != nil {
		return err
	}

	filesMetrics.Add(float64(len(foundFiles)))
	var inputs []*GlooMeshInput

	for _, file := range foundFiles {
		input, err := translateFileToMeshInput(file)
		if err != nil {
			return err
		}
		inputs = append(inputs, input)
	}

	// preprocessing
	for _, input := range inputs {
		NewPreprocessor().Preprocess(input)
	}

	// now we need to convert the easy stuff like route tables
	var outputs []*GatewayAPIOutput
	for _, input := range inputs {
		output, err := translateEdgeAPIToGatewayAPI(input, opts)
		if err != nil {
			return err
		}
		outputs = append(outputs, output)
	}

	// we need to parse through the known delegated routes and translate them
	for _, output := range outputs {
		doDelegation(output)
	}

	// write all the outputs to their files
	for _, output := range outputs {
		//only write or
		txt, err := output.ToString()
		if err != nil {
			return err
		}

		if opts.Overwrite {
			fileNameSplit := strings.Split(output.FileName, ".")
			// assuming anything before the first . is the file name
			filename := fmt.Sprintf("%s.%s", fileNameSplit[0], strings.Join(fileNameSplit[1:], "."))
			if opts.OverwriteSuffix != "" {
				fmt.Sprintf("%s-%s.%s", fileNameSplit[0], opts.OverwriteSuffix, strings.Join(fileNameSplit[1:], "."))
			}
			if output.HasItems() {
				_, _ = fmt.Fprintf(os.Stdout, "Writing File: %s\n", filename)
				file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				fmt.Fprintf(file, "%s", txt)
			} else {
				_, _ = fmt.Fprintf(os.Stdout, "Skipping File because no Edge APIs Detected: %s\n", output.FileName)
			}
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "\n\n---\n# --------------------------------\n# %s\n# --------------------------------", output.FileName)
			_, _ = fmt.Fprintf(os.Stdout, "%s\n", txt)
		}
		if opts.Stats {
			totalLines.WithLabelValues("Gateway API").Add(float64(len(strings.Split(txt, "\n"))))
		}
	}
	if opts.Stats {
		//count total lines of generated yaml (probably expensive)
		//for _, input := range inputs {
		//	txt, _ := input.ToString()
		//	totalLines.WithLabelValues("Gloo").Add(float64(len(strings.Split(txt, "\n"))))
		//}
		printMetrics(outputs)
	}

	return nil
}

// find all the HTTPRoutes that match and add a parent ref to the correct parent
func doDelegation(output *GatewayAPIOutput) {
	for _, delegateRoute := range output.DelegationReferences {
		addParentRefsToHTTPRoutes(output, delegateRoute)
	}
}

func addParentRefsToHTTPRoutes(output *GatewayAPIOutput, delegateRoute *DelegateParentReference) {

	// for each delegate route we need to go find the RouteTable that matches and assign a parent
	for _, httpRoute := range output.HTTPRoutes {
		if len(httpRoute.Labels) > 0 && doHttpRouteLabelsMatch(delegateRoute.Labels, httpRoute.Labels) {

			parentRef := gwv1.ParentReference{
				Name:      gwv1.ObjectName(delegateRoute.ParentName),
				Namespace: (*gwv1.Namespace)(&delegateRoute.ParentNamespace),
				Kind:      (*gwv1.Kind)(ptr.To("HTTPRoute")),
				Group:     (*gwv1.Group)(ptr.To("gateway.networking.k8s.io")),
			}

			httpRoute.Spec.ParentRefs = append(httpRoute.Spec.ParentRefs, parentRef)
		}
	}
}

func doHttpRouteLabelsMatch(matches map[string]string, labels map[string]string) bool {
	for k, v := range matches {
		if labels[k] != v {
			return false
		}
	}
	return true
}

func translateEdgeAPIToGatewayAPI(input *GlooMeshInput, opts *Options) (*GatewayAPIOutput, error) {

	output := &GatewayAPIOutput{
		FileName:    input.FileName,
		YamlObjects: input.YamlObjects,
	}

	return output, nil
}

func translateFileToMeshInput(fileName string) (*GlooMeshInput, error) {

	gei := &GlooMeshInput{
		FileName: fileName,
	}

	// Read the file
	data, err := os.ReadFile(fileName)
	if err != nil {
		return gei, err
	}
	for _, resourceYAML := range strings.Split(string(data), "---") {
		if len(resourceYAML) == 0 {
			continue
		}
		// yaml to object
		obj, k, err := decoder.Decode([]byte(resourceYAML), nil, nil)
		if err != nil {
			if runtime.IsNotRegisteredError(err) {
				// we just want to add the yaml and move on
				gei.YamlObjects = append(gei.YamlObjects, resourceYAML)
				continue
			}

			// TODO if we cant decode it, don't do anything and continue
			//log.Printf("# Skipping object due to error file parsing error %s", err)
			continue
		}

		// a lot of times lists are missing the group so this object doesnt match
		if k.Kind == "List" {
			var list unstructured.UnstructuredList
			if err := yaml.Unmarshal([]byte(resourceYAML), &list); err != nil {
				return nil, err
			}

			for _, item := range list.Items {

				tmpGei, err := parseObjects(item, k)
				if err != nil {
					return nil, err
				}
				gei.YamlObjects = append(gei.YamlObjects, tmpGei.YamlObjects...)
				gei.RouteTables = append(gei.RouteTables, tmpGei.RouteTables...)
			}
			continue
		}

		switch o := obj.(type) {
		case *networking_v2.RouteTable:
			glooConfigMetric.WithLabelValues("RotueTable").Inc()
			gei.RouteTables = append(gei.RouteTables, o)
		case *networking_v2.RouteTableList:
			for _, routeTable := range o.Items {
				glooConfigMetric.WithLabelValues("RouteTable").Inc()
				gei.RouteTables = append(gei.RouteTables, &routeTable)
			}
		case *networking_v2.VirtualDestination:
			glooConfigMetric.WithLabelValues("VirtualDestination").Inc()
			gei.VirtualDestinations = append(gei.VirtualDestinations, o)
		case *networking_v2.VirtualGateway:
			glooConfigMetric.WithLabelValues("VirtualGateway").Inc()
			gei.VirtualGateways = append(gei.VirtualGateways, o)
		case *networking_v2.ExternalService:
			glooConfigMetric.WithLabelValues("ExternalService").Inc()
			gei.ExternalServices = append(gei.ExternalServices, o)
		case *networking_v2.ExternalEndpoint:
			glooConfigMetric.WithLabelValues("ExternalEndpoint").Inc()
			gei.ExternalEndpoints = append(gei.ExternalEndpoints, o)
		case *traffic_v2.RateLimitPolicy:
			glooConfigMetric.WithLabelValues("RateLimitPolicy").Inc()
			gei.RateLimitPolicies = append(gei.RateLimitPolicies, o)
		case *security_v2.ExtAuthPolicy:
			glooConfigMetric.WithLabelValues("ExtAuthPolicy").Inc()
			gei.ExtAuthPolicies = append(gei.ExtAuthPolicies, o)
		case *security_v2.CORSPolicy:
			glooConfigMetric.WithLabelValues("CORSPolicy").Inc()
			gei.CORSPolicies = append(gei.CORSPolicies, o)
		case *security_v2.JWTPolicy:
			glooConfigMetric.WithLabelValues("JWTPolicy").Inc()
			gei.JWTPolicies = append(gei.JWTPolicies, o)
		case *traffic_v2.HeaderManipulationPolicy:
			glooConfigMetric.WithLabelValues("HeaderManipulationPolicy").Inc()
			gei.HeaderManipulationPolicies = append(gei.HeaderManipulationPolicies, o)
		case *resilience_v2.ConnectionPolicy:
			glooConfigMetric.WithLabelValues("ConnectionPolicy").Inc()
			gei.ConnectionPolicies = append(gei.ConnectionPolicies, o)
		case *admin_v3.RateLimitServerConfig:
			glooConfigMetric.WithLabelValues("RateLimitServerConfig").Inc()
			gei.RateLimitServerConfigs = append(gei.RateLimitServerConfigs, o)
		case *traffic_v2.RateLimitClientConfig:
			glooConfigMetric.WithLabelValues("RateLimitClientConfig").Inc()
			gei.RateLimitClientConfigs = append(gei.RateLimitClientConfigs, o)
		case *apimanagement_v2.Portal:
			glooConfigMetric.WithLabelValues("Portal").Inc()
			gei.Portals = append(gei.Portals, o)
		case *apimanagement_v2.PortalGroup:
			glooConfigMetric.WithLabelValues("PortalGroup").Inc()
			gei.PortalGroups = append(gei.PortalGroups, o)
		case *infrastructure_v3.CloudProvider:
			glooConfigMetric.WithLabelValues("CloudProvider").Inc()
			gei.CloudProviders = append(gei.CloudProviders, o)
		case *infrastructure_v3.CloudResources:
			glooConfigMetric.WithLabelValues("CloudResources").Inc()
			gei.CloudResources = append(gei.CloudResources, o)
		case *observability_v2.AccessLogPolicy:
			glooConfigMetric.WithLabelValues("AccessLogPolicy").Inc()
			gei.AccessLogPolicies = append(gei.AccessLogPolicies, o)
		case *resilience_v2.ActiveHealthCheckPolicy:
			glooConfigMetric.WithLabelValues("ActiveHealthCheckPolicy").Inc()
			gei.ActiveHealthCheckPolicies = append(gei.ActiveHealthCheckPolicies, o)
		case *resilience_v2.FailoverPolicy:
			glooConfigMetric.WithLabelValues("FailoverPolicy").Inc()
			gei.FailoverPolicies = append(gei.FailoverPolicies, o)
		case *resilience_v2.FaultInjectionPolicy:
			glooConfigMetric.WithLabelValues("FaultInjectionPolicy").Inc()
			gei.FaultInjectionPolicies = append(gei.FaultInjectionPolicies, o)
		case *resilience_v2.ListenerConnectionPolicy:
			glooConfigMetric.WithLabelValues("ListenerConnectionPolicy").Inc()
			gei.ListenerConnectionPolicies = append(gei.ListenerConnectionPolicies, o)
		case *resilience_v2.OutlierDetectionPolicy:
			glooConfigMetric.WithLabelValues("OutlierDetectionPolicy").Inc()
			gei.OutlierDetectionPolicies = append(gei.OutlierDetectionPolicies, o)
		case *resilience_v2.RetryTimeoutPolicy:
			glooConfigMetric.WithLabelValues("RetryTimeoutPolicy").Inc()
			gei.RetryTimeoutPolicies = append(gei.RetryTimeoutPolicies, o)
		case *security_v2.CSRFPolicy:
			glooConfigMetric.WithLabelValues("CSRFPolicy").Inc()
			gei.CSRFPolicies = append(gei.CSRFPolicies, o)
		case *security_v2.DLPPolicy:
			glooConfigMetric.WithLabelValues("DLPPolicy").Inc()
			gei.DLPPolicies = append(gei.DLPPolicies, o)
		case *security_v2.WAFPolicy:
			glooConfigMetric.WithLabelValues("WAFPolicy").Inc()
			gei.WAFPolicies = append(gei.WAFPolicies, o)
		case *traffic_v2.HTTPBufferPolicy:
			glooConfigMetric.WithLabelValues("HTTPBufferPolicy").Inc()
			gei.HTTPBufferPolicies = append(gei.HTTPBufferPolicies, o)
		case *traffic_v2.LoadBalancerPolicy:
			glooConfigMetric.WithLabelValues("LoadBalancerPolicy").Inc()
			gei.LoadBalancerPolicies = append(gei.LoadBalancerPolicies, o)
		case *traffic_v2.MirrorPolicy:
			glooConfigMetric.WithLabelValues("MirrorPolicy").Inc()
			gei.MirrorPolicies = append(gei.MirrorPolicies, o)
		case *traffic_v2.ProxyProtocolPolicy:
			glooConfigMetric.WithLabelValues("ProxyProtocolPolicy").Inc()
			gei.ProxyProtocolPolicies = append(gei.ProxyProtocolPolicies, o)
		case *traffic_v2.TransformationPolicy:
			glooConfigMetric.WithLabelValues("TransformationPolicy").Inc()
			gei.TransformationPolicies = append(gei.TransformationPolicies, o)
		default:
			// if we dont know what type it is we just add it back
			// no change so just add it back
			glooConfigMetric.WithLabelValues(k.Kind).Inc()
			gei.YamlObjects = append(gei.YamlObjects, resourceYAML)
		}
	}
	return gei, nil
}

func parseObjects(item unstructured.Unstructured, k *schema.GroupVersionKind) (*GlooMeshInput, error) {
	gei := &GlooMeshInput{}

	resourceYaml, err := yaml.Marshal(item)
	if err != nil {
		return nil, err
	}
	gvk := item.GroupVersionKind()
	obj, err := runtimeScheme.New(gvk)
	if runtime.IsNotRegisteredError(err) {
		// we just want to add the yaml and move on
		gei.YamlObjects = append(gei.YamlObjects, string(resourceYaml))
		return gei, nil
	} else if err != nil {
		return nil, err
	}

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
		return nil, errors.New(fmt.Sprintf("Error converting unstructured to typed: %v", err))
	}

	switch o := obj.(type) {
	case *networking_v2.RouteTable:
		glooConfigMetric.WithLabelValues("RotueTable").Inc()
		gei.RouteTables = append(gei.RouteTables, o)
	case *networking_v2.RouteTableList:
		for _, routeTable := range o.Items {
			glooConfigMetric.WithLabelValues("RouteTable").Inc()
			gei.RouteTables = append(gei.RouteTables, &routeTable)
		}
	case *networking_v2.VirtualDestination:
		glooConfigMetric.WithLabelValues("VirtualDestination").Inc()
		gei.VirtualDestinations = append(gei.VirtualDestinations, o)
	case *networking_v2.VirtualGateway:
		glooConfigMetric.WithLabelValues("VirtualGateway").Inc()
		gei.VirtualGateways = append(gei.VirtualGateways, o)
	case *networking_v2.ExternalService:
		glooConfigMetric.WithLabelValues("ExternalService").Inc()
		gei.ExternalServices = append(gei.ExternalServices, o)
	case *networking_v2.ExternalEndpoint:
		glooConfigMetric.WithLabelValues("ExternalEndpoint").Inc()
		gei.ExternalEndpoints = append(gei.ExternalEndpoints, o)
	case *traffic_v2.RateLimitPolicy:
		glooConfigMetric.WithLabelValues("RateLimitPolicy").Inc()
		gei.RateLimitPolicies = append(gei.RateLimitPolicies, o)
	case *security_v2.ExtAuthPolicy:
		glooConfigMetric.WithLabelValues("ExtAuthPolicy").Inc()
		gei.ExtAuthPolicies = append(gei.ExtAuthPolicies, o)
	case *security_v2.CORSPolicy:
		glooConfigMetric.WithLabelValues("CORSPolicy").Inc()
		gei.CORSPolicies = append(gei.CORSPolicies, o)
	case *security_v2.JWTPolicy:
		glooConfigMetric.WithLabelValues("JWTPolicy").Inc()
		gei.JWTPolicies = append(gei.JWTPolicies, o)
	case *traffic_v2.HeaderManipulationPolicy:
		glooConfigMetric.WithLabelValues("HeaderManipulationPolicy").Inc()
		gei.HeaderManipulationPolicies = append(gei.HeaderManipulationPolicies, o)
	case *resilience_v2.ConnectionPolicy:
		glooConfigMetric.WithLabelValues("ConnectionPolicy").Inc()
		gei.ConnectionPolicies = append(gei.ConnectionPolicies, o)
	case *admin_v3.RateLimitServerConfig:
		glooConfigMetric.WithLabelValues("RateLimitServerConfig").Inc()
		gei.RateLimitServerConfigs = append(gei.RateLimitServerConfigs, o)
	case *traffic_v2.RateLimitClientConfig:
		glooConfigMetric.WithLabelValues("RateLimitClientConfig").Inc()
		gei.RateLimitClientConfigs = append(gei.RateLimitClientConfigs, o)
	case *apimanagement_v2.Portal:
		glooConfigMetric.WithLabelValues("Portal").Inc()
		gei.Portals = append(gei.Portals, o)
	case *apimanagement_v2.PortalGroup:
		glooConfigMetric.WithLabelValues("PortalGroup").Inc()
		gei.PortalGroups = append(gei.PortalGroups, o)
	case *infrastructure_v3.CloudProvider:
		glooConfigMetric.WithLabelValues("CloudProvider").Inc()
		gei.CloudProviders = append(gei.CloudProviders, o)
	case *infrastructure_v3.CloudResources:
		glooConfigMetric.WithLabelValues("CloudResources").Inc()
		gei.CloudResources = append(gei.CloudResources, o)
	case *observability_v2.AccessLogPolicy:
		glooConfigMetric.WithLabelValues("AccessLogPolicy").Inc()
		gei.AccessLogPolicies = append(gei.AccessLogPolicies, o)
	case *resilience_v2.ActiveHealthCheckPolicy:
		glooConfigMetric.WithLabelValues("ActiveHealthCheckPolicy").Inc()
		gei.ActiveHealthCheckPolicies = append(gei.ActiveHealthCheckPolicies, o)
	case *resilience_v2.FailoverPolicy:
		glooConfigMetric.WithLabelValues("FailoverPolicy").Inc()
		gei.FailoverPolicies = append(gei.FailoverPolicies, o)
	case *resilience_v2.FaultInjectionPolicy:
		glooConfigMetric.WithLabelValues("FaultInjectionPolicy").Inc()
		gei.FaultInjectionPolicies = append(gei.FaultInjectionPolicies, o)
	case *resilience_v2.ListenerConnectionPolicy:
		glooConfigMetric.WithLabelValues("ListenerConnectionPolicy").Inc()
		gei.ListenerConnectionPolicies = append(gei.ListenerConnectionPolicies, o)
	case *resilience_v2.OutlierDetectionPolicy:
		glooConfigMetric.WithLabelValues("OutlierDetectionPolicy").Inc()
		gei.OutlierDetectionPolicies = append(gei.OutlierDetectionPolicies, o)
	case *resilience_v2.RetryTimeoutPolicy:
		glooConfigMetric.WithLabelValues("RetryTimeoutPolicy").Inc()
		gei.RetryTimeoutPolicies = append(gei.RetryTimeoutPolicies, o)
	case *security_v2.CSRFPolicy:
		glooConfigMetric.WithLabelValues("CSRFPolicy").Inc()
		gei.CSRFPolicies = append(gei.CSRFPolicies, o)
	case *security_v2.DLPPolicy:
		glooConfigMetric.WithLabelValues("DLPPolicy").Inc()
		gei.DLPPolicies = append(gei.DLPPolicies, o)
	case *security_v2.WAFPolicy:
		glooConfigMetric.WithLabelValues("WAFPolicy").Inc()
		gei.WAFPolicies = append(gei.WAFPolicies, o)
	case *traffic_v2.HTTPBufferPolicy:
		glooConfigMetric.WithLabelValues("HTTPBufferPolicy").Inc()
		gei.HTTPBufferPolicies = append(gei.HTTPBufferPolicies, o)
	case *traffic_v2.LoadBalancerPolicy:
		glooConfigMetric.WithLabelValues("LoadBalancerPolicy").Inc()
		gei.LoadBalancerPolicies = append(gei.LoadBalancerPolicies, o)
	case *traffic_v2.MirrorPolicy:
		glooConfigMetric.WithLabelValues("MirrorPolicy").Inc()
		gei.MirrorPolicies = append(gei.MirrorPolicies, o)
	case *traffic_v2.ProxyProtocolPolicy:
		glooConfigMetric.WithLabelValues("ProxyProtocolPolicy").Inc()
		gei.ProxyProtocolPolicies = append(gei.ProxyProtocolPolicies, o)
	case *traffic_v2.TransformationPolicy:
		glooConfigMetric.WithLabelValues("TransformationPolicy").Inc()
		gei.TransformationPolicies = append(gei.TransformationPolicies, o)
	default:
		// if we dont know what type it is we just add it back
		// no change so just add it back
		glooConfigMetric.WithLabelValues(k.Kind).Inc()
		gei.YamlObjects = append(gei.YamlObjects, string(resourceYaml))
	}
	return gei, nil
}

func init() {
	runtimeScheme = runtime.NewScheme()

	//if err := metav1.AddToSche; err != nil {
	//	log.Fatal(err)
	//}
	//if err := glookube.AddToScheme(runtimeScheme); err != nil {
	//	log.Fatal(err)
	//}
	//if err := gatewaykube.AddToScheme(runtimeScheme); err != nil {
	//	log.Fatal(err)
	//}
	//if err := v2.AddToScheme(runtimeScheme); err != nil {
	//	log.Fatal(err)
	//}
	//if err := gwv1.Install(runtimeScheme); err != nil {
	//	log.Fatal(err)
	//}

	if err := SchemeBuilder.AddToScheme(runtimeScheme); err != nil {
		log.Fatal(err)
	}

	codecs = serializer.NewCodecFactory(runtimeScheme)
	decoder = codecs.UniversalDeserializer()
}

func findFiles(opts *Options) ([]string, error) {
	var files []string
	if opts.Directory != "" {
		fs, err := findYamlFiles(opts.Directory)
		if err != nil {
			return nil, err
		}
		files = fs
	} else {
		files = append(files, opts.InputFile)
	}

	return files, nil
}

func findYamlFiles(directory string) ([]string, error) {
	var files []string
	libRegEx, e := regexp.Compile("^.+\\.(yaml|yml)$")
	if e != nil {
		return nil, e
	}

	e = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			// println(path)
			if !strings.Contains(info.Name(), "kustomization") {
				files = append(files, path)
			}
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return files, nil
}

// This function validates that the RouteRable matchers are the same match type prefix or exact
// The reason being is that if you are doing a rewrite you can only have one type of filter applied
func validateMatchersAreSame(matches []*matchers.Matcher) error {

	var foundExact, foundPrefix, foundRegex bool
	for _, m := range matches {
		if m.GetExact() != "" {
			if foundPrefix || foundRegex {
				return errors.New("multiple matchers found")
			}
			foundExact = true
		}
		if m.GetPrefix() != "" {
			if foundExact || foundRegex {
				return errors.New("multiple matchers found")
			}
			foundPrefix = true
		}
		if m.GetRegex() != "" {
			if foundExact || foundPrefix {
				return errors.New("multiple matchers found")
			}
			foundRegex = true
		}
	}
	return nil
}

// tests to see if all matchers are exact
func isExactMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetExact() != "" {
			return true, nil
		}
	}
	return false, nil
}

// tests to see if all matchers are exact
func isPrefixMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetPrefix() != "" {
			return true, nil
		}
	}
	return false, nil
}

// tests to see if all matchers are regex
func isRegexMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetRegex() != "" {
			return true, nil
		}
	}
	return false, nil
}
