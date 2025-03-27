package envoy

import (
	"fmt"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/external/envoy/extensions/waf"
	"github.com/spf13/cobra"
	"os"
	"time"

	_ "github.com/cncf/xds/go/udpa/type/v1"
	udpav1 "github.com/cncf/xds/go/udpa/type/v1"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	glooenvoywaf "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/waf"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	_ "istio.io/api/envoy/config/filter/http/alpn/v2alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func RootCmd() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:   "envoy",
		Short: "Convert Envoy Config to Gateway API",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}
	opts.addToFlags(cmd.PersistentFlags())

	cmd.SilenceUsage = true
	return cmd
}

func run(opts *Options) error {

	msgs, err := load(opts.InputFile)
	if err != nil {
		return err
	}
	outputs := Outputs{
		OutputDir:          opts.OutputDir,
		FolderPerNamespace: opts.FolderPerNamespace,
		Settings:           v1.Settings{},
		Errors:             make([]error, 0),
	}
	if err := outputs.Convert(msgs); err != nil {
		return err
	}
	if err := outputs.PostProcess(opts.RouteTableFile); err != nil {
		return err
	}
	if err := outputs.Write(); err != nil {
		return err
	}
	return nil
}

func load(inputFile string) ([]proto.Message, error) {
	fmt.Printf("Loading config file %s\n", inputFile)
	// Read the Envoy configuration file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, err
	}

	var cd adminv3.ConfigDump

	m := protojson.UnmarshalOptions{}
	err = m.Unmarshal(data, &cd)
	if err != nil {
		return nil, err
	}

	var msgs []proto.Message

	for _, r := range cd.Configs {
		msg, err := r.UnmarshalNew()
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	fmt.Println("File loading complete")
	return msgs, nil
}

func (o *Outputs) Convert(msgs []proto.Message) error {
	//var bootstrap *adminv3.BootstrapConfigDump
	var listeners *adminv3.ListenersConfigDump
	var routes *adminv3.RoutesConfigDump
	var clusters *adminv3.ClustersConfigDump
	for _, msg := range msgs {
		switch m := msg.(type) {
		//case *adminv3.BootstrapConfigDump:
		//	bootstrap = m
		case *adminv3.ListenersConfigDump:
			listeners = m
		case *adminv3.RoutesConfigDump:
			routes = m
		case *adminv3.ClustersConfigDump:
			clusters = m
			err := o.AddClusters(clusters)
			if err != nil {
				return err
			}
		}
	}

	//gwName := bootstrap.GetBootstrap().GetNode().GetId()
	gwName := "ingress"

	o.Gateway = &gwv1.Gateway{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Gateway",
			APIVersion: "gateway.networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: gwName,
		},
		Spec: gwv1.GatewaySpec{
			GatewayClassName: "gateway.solo.io",
			Listeners:        []gwv1.Listener{},
		},
	}
	fmt.Printf("Parsing %d listeners\n", len(listeners.DynamicListeners))
	for _, dl := range listeners.DynamicListeners {
		var l listenerv3.Listener
		err := dl.GetActiveState().GetListener().UnmarshalTo(&l)
		if err != nil {
			return err
		}
		// translate each filter chain
		l2, rds := o.convertListener(&l)

		parentRef := gwv1.ParentReference{
			Name:      gwv1.ObjectName(gwName),
			Group:     ptr.To(gwv1.Group(gwv1.GroupName)),
			Kind:      ptr.To(gwv1.Kind("Gateway")),
			Namespace: ptr.To(gwv1.Namespace(o.Gateway.Namespace)),
		}
		if len(l2) == 1 {
			parentRef.SectionName = ptr.To(l2[0].Name)
		}
		fmt.Printf("\tParsing %d httpRoutes for listener\n", len(rds))
		for _, r := range rds {
			o.doRoutes(routes, parentRef, r)
		}

		o.Gateway.Spec.Listeners = append(o.Gateway.Spec.Listeners, l2...)
	}

	return nil
}

func convertDuration(d *durationpb.Duration) gwv1.Duration {
	dur := d.AsDuration()
	durMs := dur.Milliseconds()
	// we need this 4 digits MSB with wither h, m, s, ms.
	if dur < 9999*time.Millisecond {
		return gwv1.Duration(fmt.Sprintf("%dms", durMs))
	}
	if dur < 9999*time.Second {
		return gwv1.Duration(fmt.Sprintf("%ds", durMs/1e3))
	}
	if dur < 9999*time.Minute {
		return gwv1.Duration(fmt.Sprintf("%dm", durMs/1e3))
	}
	panic("duration too big")
}

func convertAny(anym *anypb.Any) (proto.Message, error) {
	v, err := anym.UnmarshalNew()
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling per filter config %w", err)

	}
	if ts, ok := v.(*udpav1.TypedStruct); ok {
		if ts.GetTypeUrl() == "type.googleapis.com/io.istio.http.peer_metadata.Config" || ts.GetTypeUrl() == "type.googleapis.com/envoy.tcp.metadataexchange.config.MetadataExchange" {
			// no need to convert this, and there's no pb.go for it..
			return nil, nil
		}
		msgt, err := protoregistry.GlobalTypes.FindMessageByURL(ts.GetTypeUrl())
		if err != nil {
			return nil, fmt.Errorf("error finding message by url %w", err)

		}
		msg := msgt.New().Interface()
		// marshal to json
		b, err := protojson.Marshal(ts.Value)
		if err != nil {
			return nil, fmt.Errorf("error marshalling to json %w", err)

		}
		// unmarshal to gateway route option
		err = protojson.Unmarshal(b, msg)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling to gateway route option %w", err)
		}
		return msg, nil
	}
	return v, nil
}

func getRef(cluster string) *core.ResourceRef {
	// TODO: this needs to be the name of a gloo upstream
	// need to map this concept to ggv2/kgateway
	return &core.ResourceRef{
		Name: cluster,
	}
}

func protoRoundTrip[P interface {
	proto.Message
	*T
}, T any](m proto.Message) P {
	if m == nil {
		return nil
	}
	b, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	var t T
	var p P = &t
	if err := proto.Unmarshal(b, p); err != nil {
		panic(err)
	}
	return p
}

//func convertTrue(b *wrapperspb.BoolValue) *gwv1.LowercaseTrue {
//	if !b.GetValue() {
//		return nil
//	}
//	tr := gwv1.LowercaseTrue("true")
//	return &tr
//}

func convertRuleSets(rules []*waf.RuleSet) []*glooenvoywaf.RuleSet {
	var ret []*glooenvoywaf.RuleSet
	for _, rs := range rules {
		ret = append(ret, &glooenvoywaf.RuleSet{
			RuleStr:   rs.RuleStr,
			Files:     rs.Files,
			Directory: rs.Directory,
		})
	}

	return ret
}

func convertSlice[T ~string](s []string) []T {
	var ret []T
	for _, v := range s {
		ret = append(ret, T(v))
	}
	return ret
}

func convertOrigins(oms []*envoy_matcher_v3.StringMatcher) []gwv1.AbsoluteURI {

	var ret []gwv1.AbsoluteURI

	for _, om := range oms {
		switch om.GetMatchPattern().(type) {
		case *envoy_matcher_v3.StringMatcher_Exact:
			ret = append(ret, gwv1.AbsoluteURI(om.GetExact()))
		case *envoy_matcher_v3.StringMatcher_SafeRegex:

		}
	}

	return ret
}

func isEmpty(m proto.Message) bool {
	b, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	return len(b) == 0
}

type RouteInfo struct {
	Rds            string
	FiltersOnChain map[string][]proto.Message
}
