package crd

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Crd struct {
	GroupName string
	Plural    string
	Group     string
	Version   string
	KindName  string
	ShortName string
	Type      runtime.Object
}

func NewCrd(GroupName string,
	Plural string,
	Group string,
	Version string,
	KindName string,
	ShortName string,
	Type runtime.Object) Crd {
	return Crd{
		GroupName: GroupName,
		Plural:    Plural,
		Group:     Group,
		Version:   Version,
		KindName:  KindName,
		ShortName: ShortName,
		Type:      Type,
	}
}

func (d Crd) FullName() string {
	return d.Plural + "." + d.Group
}

func (d Crd) TypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       d.KindName,
		APIVersion: d.Group + "/" + d.Version,
	}
}

// SchemeGroupVersion is group version used to register these objects
func (d Crd) SchemeGroupVersion() schema.GroupVersion {
	return schema.GroupVersion{Group: d.GroupName, Version: d.Version}
}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func (d Crd) Kind(kind string) schema.GroupKind {
	return d.SchemeGroupVersion().WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func (d Crd) Resource(resource string) schema.GroupResource {
	return d.SchemeGroupVersion().WithResource(resource).GroupResource()
}

func (d Crd) SchemeBuilder() runtime.SchemeBuilder {
	return runtime.NewSchemeBuilder(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(d.SchemeGroupVersion(), d.Type)
		metav1.AddToGroupVersion(scheme, d.SchemeGroupVersion())
		return nil
	})
}

func (d Crd) AddToScheme(s *runtime.Scheme) error {
	return d.SchemeBuilder().AddToScheme(s)
}
