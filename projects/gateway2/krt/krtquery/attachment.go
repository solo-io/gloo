package krtquery

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// Attachment indexes a resource by its attachments.
// No section name means it is attached directly to the resource.
type Attachment[T Namespaced] struct {
	// of the attached item
	Namespace, Name, Section string
	Resource                 T
}

type Namespaced interface {
	GetNamespace() string
}

type TargetRef interface {
	GetNamespace() *wrappers.StringValue
	GetName() string
	GetSectionName() *wrappers.StringValue
}

// ResourceName is used by krt for keying elements of a Collection.
// We query for the attachments themselves using this key.
func (h Attachment[T]) ResourceName() string {
	return h.Namespace + "/" + h.Name + "/" + h.Section
}

// TargetKey returns the key of the attached Kubernetes resource (no section name).
func (h Attachment[T]) TargetKey() string {
	return h.Namespace + "/" + h.Name
}

func (a Attachment[T]) AttachedToGateway(ctx krt.HandlerContext, Gateways krt.Collection[*gwv1.Gateway]) bool {
	// TODO check ref Group/Kind
	gw := ptr.Flatten(krt.FetchOne(ctx, Gateways, krt.FilterKey(a.TargetKey())))
	if gw == nil {
		return false
	}
	if a.Section == "" {
		return true
	}
	found := false
	for _, l := range gw.Spec.Listeners {
		if string(l.Name) == a.Section {
			found = true
			break
		}
	}
	return found
}

func attachementFromParentRef[T Namespaced](resource T, ref gwv1.ParentReference) Attachment[T] {
	return Attachment[T]{
		Namespace: string(ptr.OrDefault(ref.Namespace, gwv1.Namespace(resource.GetNamespace()))),
		Name:      string(ref.Name),
		Section:   string(ptr.OrEmpty(ref.SectionName)),
		Resource:  resource,
	}
}

func attachementFromTargetRef[T Namespaced](resource T, ref TargetRef) Attachment[T] {
	return Attachment[T]{
		Namespace: ptr.NonEmptyOrDefault(ref.GetNamespace().GetValue(), resource.GetNamespace()),
		Name:      ref.GetName(),
		Section:   ref.GetSectionName().GetValue(),
		Resource:  resource,
	}
}
