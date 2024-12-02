package proxy_syncer

import (
	"encoding/json"
	"fmt"

	udpaannontations "github.com/cncf/xds/go/udpa/annotations"
	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"

	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"istio.io/istio/pkg/kube/krt"

	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/gloo/projects/gateway2/translator/translatorutils"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"

	oldproto "github.com/golang/protobuf/proto"
)

var (
	UseDetailedUnmarshalling = envutils.IsEnvTruthy("DETAILED_SNAP_UNMARSHALLING")
)

type XdsSnapWrapper struct {
	snap            *xds.EnvoySnapshot
	proxyKey        string
	proxyWithReport translatorutils.ProxyWithReports
	pluginRegistry  registry.PluginRegistry
	fullReports     reporter.ResourceReports
}

func (p XdsSnapWrapper) WithSnapshot(snap *xds.EnvoySnapshot) XdsSnapWrapper {
	p.snap = snap
	return p
}

var _ krt.ResourceNamer = XdsSnapWrapper{}

func (p XdsSnapWrapper) Equals(in XdsSnapWrapper) bool {
	return p.snap.Equal(in.snap)
}

func (p XdsSnapWrapper) ResourceName() string {
	return p.proxyKey
}

// note: this is feature gated, as i'm not confident the new logic can't panic, in all envoy configs
// once 1.18 is out, we can remove the feature gate.
func (p XdsSnapWrapper) MarshalJSON() (out []byte, err error) {
	if !UseDetailedUnmarshalling {
		// use a new struct to prevent infinite recursion
		return json.Marshal(struct {
			snap            *xds.EnvoySnapshot
			proxyKey        string
			proxyWithReport translatorutils.ProxyWithReports
			pluginRegistry  registry.PluginRegistry
			fullReports     reporter.ResourceReports
		}{
			snap:            p.snap,
			proxyKey:        p.proxyKey,
			proxyWithReport: p.proxyWithReport,
			pluginRegistry:  p.pluginRegistry,
			fullReports:     p.fullReports,
		})
	}

	snap := p.snap.Clone().(*xds.EnvoySnapshot)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic handling snapshot: %v", r)
		}
	}()

	// redact things
	redact(snap)
	snapJson := map[string][]any{}
	for _, r := range snap.Listeners.Items {
		rJson, _ := protojson.Marshal(r.ResourceProto().(proto.Message))
		var rAny any
		json.Unmarshal(rJson, &rAny)
		snapJson["Listeners"] = append(snapJson["Listeners"], rAny)
	}
	for _, r := range snap.Clusters.Items {
		rJson, _ := protojson.Marshal(r.ResourceProto().(proto.Message))
		var rAny any
		json.Unmarshal(rJson, &rAny)
		snapJson["Clusters"] = append(snapJson["Clusters"], rAny)
	}

	for _, r := range snap.Routes.Items {
		rJson, _ := protojson.Marshal(r.ResourceProto().(proto.Message))
		var rAny any
		json.Unmarshal(rJson, &rAny)
		snapJson["Routes"] = append(snapJson["Routes"], rAny)
	}

	for _, r := range snap.Endpoints.Items {
		rJson, _ := protojson.Marshal(r.ResourceProto().(proto.Message))
		var rAny any
		json.Unmarshal(rJson, &rAny)
		snapJson["Endpoints"] = append(snapJson["Endpoints"], rAny)
	}

	return json.Marshal(struct {
		Snap     any
		ProxyKey string
	}{
		Snap:     snapJson,
		ProxyKey: p.proxyKey,
	})
}

func redact(snap *xds.EnvoySnapshot) {
	// clusters and listener might have secrets
	for _, l := range snap.Listeners.Items {
		redactProto(l.ResourceProto())
	}
	for _, l := range snap.Clusters.Items {
		redactProto(l.ResourceProto())
	}
}

func redactProto(m oldproto.Message) {
	var msg proto.Message = m.(proto.Message)
	visitFields(msg.ProtoReflect(), false)
}

func isSensitive(fd protoreflect.FieldDescriptor) bool {
	opts := fd.Options().(*descriptorpb.FieldOptions)
	if !proto.HasExtension(opts, udpaannontations.E_Sensitive) {
		return false
	}

	maybeExt := proto.GetExtension(opts, udpaannontations.E_Sensitive)
	return maybeExt.(bool)
}

func visitFields(msg protoreflect.Message, ancestor_sensitive bool) {
	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		sensitive := ancestor_sensitive || isSensitive(fd)

		if fd.IsList() {
			list := v.List()
			for i := 0; i < list.Len(); i++ {
				elem := list.Get(i)
				if fd.Message() != nil {
					visitMessage(msg, fd, elem, sensitive)
				} else {
					// Redact scalar fields if needed
					if sensitive {
						list.Set(i, redactValue(fd, elem))
					}
				}
			}
		} else if fd.IsMap() {
			m := v.Map()
			m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
				if fd.MapValue().Message() != nil {
					visitMessage(msg, fd.MapValue(), v, sensitive)
				} else {
					// Redact scalar fields if needed
					if sensitive {
						m.Set(k, redactValue(fd.MapValue(), v))
					}
				}
				return true
			})
		} else {
			if fd.Message() != nil {
				visitMessage(msg, fd, v, sensitive)
			} else {
				// Redact scalar fields if needed
				if sensitive {
					msg.Set(fd, redactValue(fd, v))
				}
			}
		}
		return true
	})
}

func visitMessage(msg protoreflect.Message, fd protoreflect.FieldDescriptor, v protoreflect.Value, sensitive bool) {
	visitMsg := v.Message()
	var anyMsg proto.Message
	m := visitMsg.Interface()
	if anymsg, ok := m.(*anypb.Any); ok {
		anyMsg, _ = anypb.UnmarshalNew(anymsg, proto.UnmarshalOptions{})
		visitMsg = anyMsg.ProtoReflect()

	}
	visitFields(visitMsg, sensitive)
	if anyMsg != nil {
		anymsg, _ := anypb.New(anyMsg)
		msg.Set(fd, protoreflect.ValueOf(anymsg.ProtoReflect()))
	}
}

func redactValue(fd protoreflect.FieldDescriptor, v protoreflect.Value) protoreflect.Value {
	switch fd.Kind() {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString("[REDACTED]")
	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes([]byte("[REDACTED]"))
	}
	return v
}
