package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	CRDPlural   = "apis"
	CRDGroup    = "public.delta.io"
	CRDVersion  = "v1"
	CRDFullName = CRDPlural + "." + CRDGroup
	Kind        = "ApiDefinition"
)

var ApiDefinitionsResource = schema.GroupVersionResource{Group: CRDGroup, Version: CRDVersion, Resource: CRDPlural}

var ApiDefinitionsKind = schema.GroupVersionKind{Group: CRDGroup, Version: CRDVersion, Kind: Kind}

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: CRDGroup, Version: "v1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ApiDefinition{},
		&ApiDefinitionList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
