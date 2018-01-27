package v1

import meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +resource:path=apidefinition
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiDefinition is the generic Kubernetes API object wrapper
type ApiDefinition struct {
	meta_v1.TypeMeta `json:",inline"`
	// +optional
	meta_v1.ObjectMeta `json:"metadata,omitempty"`
	Spec               map[string]string `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiDefinitionList is the generic Kubernetes API object wrapper
type ApiDefinitionList struct {
	meta_v1.TypeMeta `json:",inline"`
	// +optional
	meta_v1.ObjectMeta `json:"metadata,omitempty"`
	Items              []ApiDefinition `json:"items"`
}
