package printers

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/grpc_json"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/xdsinspection"
	plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws/ec2"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/olekukonko/tablewriter"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

// PrintUpstreams
func PrintUpstreams(upstreams v1.UpstreamList, outputType OutputType, xdsDump *xdsinspection.XdsDump, shouldListGrpcMethods bool) error {
	if outputType == KUBE_YAML {
		return PrintKubeCrdList(upstreams.AsInputResources(), v1.UpstreamCrd)
	}
	if outputType == JSON && shouldListGrpcMethods {
		return printGrpcServiceAsJson(upstreams, os.Stdout)
	}
	return cliutils.PrintList(outputType.String(), "", upstreams,
		func(data interface{}, w io.Writer) error {
			UpstreamTable(xdsDump, data.(v1.UpstreamList), w)
			return nil
		}, os.Stdout)
}

// PrintTable prints upstreams using tables to io.Writer
func UpstreamTable(xdsDump *xdsinspection.XdsDump, upstreams []*v1.Upstream, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Upstream", "type", "status", "details"})

	for _, us := range upstreams {
		name := us.GetMetadata().GetName()
		s := upstreamStatus(us)

		u := upstreamType(us)
		details := upstreamDetails(us, xdsDump)

		if len(details) == 0 {
			details = []string{""}
		}
		for i, line := range details {
			if i == 0 {
				table.Append([]string{name, u, s, line})
			} else {
				table.Append([]string{"", "", "", line})
			}
		}
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}

func upstreamStatus(us *v1.Upstream) string {
	return AggregateNamespacedStatuses(us.GetNamespacedStatuses(), func(status *core.Status) string {
		return status.GetState().String()
	})
}

func upstreamType(up *v1.Upstream) string {
	if up == nil {
		return "Invalid"
	}

	switch up.GetUpstreamType().(type) {
	case *v1.Upstream_Aws:
		return "AWS Lambda"
	case *v1.Upstream_Azure:
		return "Azure"
	case *v1.Upstream_Consul:
		return "Consul"
	case *v1.Upstream_AwsEc2:
		return "AWS EC2"
	case *v1.Upstream_Kube:
		return "Kubernetes"
	case *v1.Upstream_Static:
		return "Static"
	case *v1.Upstream_Gcp:
		return "GCP"
	default:
		return "Unknown"
	}
}

func upstreamDetails(up *v1.Upstream, xdsDump *xdsinspection.XdsDump) []string {
	if up == nil {
		return []string{"invalid: upstream was nil"}
	}

	var details []string
	add := func(s ...string) {
		details = append(details, s...)
	}
	switch usType := up.GetUpstreamType().(type) {
	case *v1.Upstream_Aws:
		var functions []string
		for _, fn := range usType.Aws.GetLambdaFunctions() {
			functions = append(functions, fn.GetLambdaFunctionName())
		}

		add(
			fmt.Sprintf("region: %v", usType.Aws.GetRegion()),
			fmt.Sprintf("secret: %s", stringifyKey(usType.Aws.GetSecretRef())),
		)
		for i := range functions {
			if i == 0 {
				add("functions:")
			}
			add(fmt.Sprintf("- %v", functions[i]))
		}
	case *v1.Upstream_AwsEc2:
		add(
			fmt.Sprintf("role:           %v", usType.AwsEc2.GetRoleArn()),
			fmt.Sprintf("uses public ip: %v", usType.AwsEc2.GetPublicIp()),
			fmt.Sprintf("port:           %v", usType.AwsEc2.GetPort()),
		)
		add(getEc2TagFiltersString(usType.AwsEc2.GetFilters())...)
		instances := xdsDump.GetEc2InstancesForUpstream(up.GetMetadata().Ref())
		add(
			"EC2 Instance Ids:",
		)
		add(
			instances...,
		)
	case *v1.Upstream_Azure:
		var functions []string
		for _, fn := range usType.Azure.GetFunctions() {
			functions = append(functions, fn.GetFunctionName())
		}
		add(
			fmt.Sprintf("function app name: %v", usType.Azure.GetFunctionAppName()),
			fmt.Sprintf("secret: %s", stringifyKey(usType.Azure.GetSecretRef())),
		)

		for i := range functions {
			if i == 0 {
				add("functions:")
			}
			add(fmt.Sprintf("- %v", functions[i]))
		}
	case *v1.Upstream_Consul:
		add(
			fmt.Sprintf("svc name: %v", usType.Consul.GetServiceName()),
			fmt.Sprintf("svc tags: %v", usType.Consul.GetServiceTags()),
		)
		if usType.Consul.GetServiceSpec() != nil {
			add(linesForServiceSpec(usType.Consul.GetServiceSpec())...)
		}
	case *v1.Upstream_Kube:
		add(
			fmt.Sprintf("svc name:      %v", usType.Kube.GetServiceName()),
			fmt.Sprintf("svc namespace: %v", usType.Kube.GetServiceNamespace()),
			fmt.Sprintf("port:          %v", usType.Kube.GetServicePort()),
		)
		if usType.Kube.GetServiceSpec() != nil {
			add(linesForServiceSpec(usType.Kube.GetServiceSpec())...)
		}
	case *v1.Upstream_Static:
		for i := range usType.Static.GetHosts() {
			if i == 0 {
				add("hosts:")
			}
			add(fmt.Sprintf("- %v:%v", usType.Static.GetHosts()[i].GetAddr(), usType.Static.GetHosts()[i].GetPort()))
		}
		if usType.Static.GetServiceSpec() != nil {
			add(linesForServiceSpec(usType.Static.GetServiceSpec())...)
		}
	case *v1.Upstream_Gcp:
		add(fmt.Sprintf("host: %v", usType.Gcp.GetHost()))
		if usType.Gcp.GetAudience() != "" {
			add(fmt.Sprintf("host: %v", usType.Gcp.GetAudience()))
		}

	}
	add("")
	return details
}

func linesForServiceSpec(serviceSpec *plugins.ServiceSpec) []string {
	var spec []string
	add := func(s ...string) {
		spec = append(spec, s...)
	}
	switch plug := serviceSpec.GetPluginType().(type) {
	case *plugins.ServiceSpec_Rest:
		add("REST service:")
		var functions []string
		for restFunc, transform := range plug.Rest.GetTransformations() {
			var method, path string
			methodP, ok := transform.GetHeaders()[":method"]
			if ok {
				method = methodP.GetText()
			}
			pathP, ok := transform.GetHeaders()[":path"]
			if ok {
				path = pathP.GetText()
			}
			if false {
				//TODO ilackarms: save it for -o wide
				functions = append(functions, fmt.Sprintf("- %v (%v %v)", restFunc, method, path))
			}
			functions = append(functions, fmt.Sprintf("- %v", restFunc))
		}
		// needed because map
		sort.Strings(functions)

		for i := range functions {
			if i == 0 {
				add("functions:")
			}
			add(functions[i])
		}
	case *plugins.ServiceSpec_Grpc:
		add("gRPC service:")
		for _, grpcService := range plug.Grpc.GetGrpcServices() {
			add(fmt.Sprintf("  %v", grpcService.GetServiceName()))
			for _, fn := range grpcService.GetFunctionNames() {
				add(fmt.Sprintf("  - %v", fn))
			}
		}
	case *plugins.ServiceSpec_GrpcJsonTranscoder:
		add("gRPC service:")
		descriptorSet := plug.GrpcJsonTranscoder.GetProtoDescriptorBin()
		for _, grpcService := range plug.GrpcJsonTranscoder.GetServices() {
			add(fmt.Sprintf("  %v", grpcService))
			md := getAllFuncNames(grpcService, descriptorSet)
			for i := 0; i < md.Len(); i++ {
				add(fmt.Sprintf("  - %v", md.Get(i).Name()))
			}
		}
	}

	return spec
}

func getAllFuncNames(service string, descriptorSet []byte) protoreflect.MethodDescriptors {
	fds := &descriptor.FileDescriptorSet{}
	err := proto.Unmarshal(descriptorSet, fds)
	if err != nil {
		panic("unable to unmarshal descriptor")
	}
	files, _ := protodesc.NewFiles(fds)
	d, _ := files.FindDescriptorByName(protoreflect.FullName(service))
	s, ok := d.(protoreflect.ServiceDescriptor)
	if !ok {
		panic("unable to decode service descriptor")
	}
	return s.Methods()
}

// stringifyKey for a resource likely could be done more nicely with spew
// or a better accessor but minimal this avoids panicing on nested references to nils
func stringifyKey(plausiblyNilRef *core.ResourceRef) string {

	if plausiblyNilRef == nil {
		return "<None>"
	}
	return plausiblyNilRef.Key()

}

func getEc2TagFiltersString(filters []*ec2.TagFilter) []string {
	var out []string
	add := func(s ...string) {
		out = append(out, s...)
	}

	var kFilters []*ec2.TagFilter_Key
	var kvFilters []*ec2.TagFilter_KvPair
	for _, f := range filters {
		switch x := f.GetSpec().(type) {
		case *ec2.TagFilter_Key:
			kFilters = append(kFilters, x)
		case *ec2.TagFilter_KvPair_:
			kvFilters = append(kvFilters, x.KvPair)
		}
	}
	if len(kFilters) == 0 {
		add(fmt.Sprintf("key filters: (none)"))
	} else {
		add(fmt.Sprintf("key filters:"))
		for _, f := range kFilters {
			add(fmt.Sprintf("- %v", f.Key))
		}
	}
	if len(kvFilters) == 0 {
		add(fmt.Sprintf("key-value filters: (none)"))
	} else {
		add(fmt.Sprintf("key-value filters:"))
		for _, f := range kvFilters {
			add(fmt.Sprintf("- %v: %v", f.GetKey(), f.GetValue()))
		}
	}
	return out
}

type UpstreamGrpc struct {
	Name      string      `json:"upstreamName,omitempty"`
	Namespace string      `json:"upstreamNamespace,omitempty"`
	GrpcAttrs *[]GrpcAttr `json:"grpc,omitempty"`
}

type GrpcAttr struct {
	PackageName   string   `json:"packageName,omitempty"`
	ServiceName   string   `json:"serviceName,omitempty"`
	FunctionNames []string `json:"functionNames,omitempty"`
}

func BuildUpstreamGrpcAsJson(up *v1.Upstream) string {
	switch usType := up.GetUpstreamType().(type) {
	case *v1.Upstream_Static:
		return buildUpstreamGrpc(up.GetMetadata(), usType.GetServiceSpec())
	case *v1.Upstream_Kube:
		return buildUpstreamGrpc(up.GetMetadata(), usType.GetServiceSpec())
	}
	return ""
}

func buildUpstreamGrpc(metadata *core.Metadata, serviceSpec *plugins.ServiceSpec) string {
	ug := &UpstreamGrpc{}
	switch plug := serviceSpec.GetPluginType().(type) {
	case *plugins.ServiceSpec_GrpcJsonTranscoder:
		ug.Name = metadata.GetName()
		ug.Namespace = metadata.GetNamespace()
		ug.GrpcAttrs = buildGrpcAttr(plug.GrpcJsonTranscoder)
		return convertToJson(ug)
	}
	return ""
}

func convertToJson(ug *UpstreamGrpc) string {
	j, err := json.Marshal(ug)
	if err != nil {
		return ""
	}
	return string(j)
}

func buildGrpcAttr(gjt *grpc_json.GrpcJsonTranscoder) *[]GrpcAttr {
	var gaList []GrpcAttr
	descriptorSet := gjt.GetProtoDescriptorBin()
	for _, grpcService := range gjt.GetServices() {
		ga := GrpcAttr{}
		parts := strings.Split(grpcService, ".")
		ga.ServiceName = parts[len(parts)-1]
		ga.PackageName = strings.Join(parts[:len(parts)-1], ".")
		md := getAllFuncNames(grpcService, descriptorSet)
		var funcNames []string
		for i := 0; i < md.Len(); i++ {
			funcNames = append(funcNames, string(md.Get(i).Name()))
		}
		ga.FunctionNames = funcNames
		gaList = append(gaList, ga)
	}
	return &gaList
}

func printGrpcServiceAsJson(data interface{}, w io.Writer) error {
	list := reflect.ValueOf(data)
	_, err := fmt.Fprintln(w, "{")
	_, err = fmt.Fprintln(w, "\"upstreams\": [")
	if err != nil {
		return errors.Wrap(err, "unable to print JSON list")
	}
	for i := 0; i < list.Len(); i++ {
		v, ok := list.Index(i).Interface().(proto.Message)
		if !ok {
			return eris.New("unable to convert to proto message")
		}
		if i != 0 {
			_, err = fmt.Fprintln(w, ",")
			if err != nil {
				return errors.Wrap(err, "unable to print JSON list")
			}
		}
		err = cliutils.PrintJSON(v, w)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintln(w, "]")
	_, err = fmt.Fprintln(w, ", \"grpcServices\": [")
	if err != nil {
		return errors.Wrap(err, "unable to print JSON list")
	}
	for i := 0; i < list.Len(); i++ {
		v, ok := list.Index(i).Interface().(proto.Message)
		if !ok {
			return eris.New("unable to convert to proto message")
		}
		j := BuildUpstreamGrpcAsJson(v.(*v1.Upstream))
		if j != "" {
			if i != 0 {
				_, err = fmt.Fprintln(w, ",")
				if err != nil {
					return errors.Wrap(err, "unable to print JSON list")
				}
			}
			_, err = fmt.Fprintln(w, j)
			if err != nil {
				return errors.Wrap(err, "unable to print JSON list")
			}
		}
	}
	_, err = fmt.Fprintln(w, "]")
	_, err = fmt.Fprintln(w, "}")
	return err
}
