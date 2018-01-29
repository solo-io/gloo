package v1

import (
	"encoding/json"

	"github.com/solo-io/glue/pkg/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Route is the generic Kubernetes API object wrapper for Glue Routes
type Route struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            Status        `json:"status,omitempty"`
	Spec              DeepCopyRoute `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RouteList is the generic Kubernetes API object wrapper
type RouteList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []Route `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Upstream is the generic Kubernetes API object wrapper for Glue Upstreams
type Upstream struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            Status `json:"status,omitempty"`

	Spec DeepCopyUpstream `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UpstreamList is the generic Kubernetes API object wrapper
type UpstreamList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []Upstream `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualHost is the generic Kubernetes API object wrapper for Glue VirtualHosts
type VirtualHost struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            Status `json:"status,omitempty"`

	Spec DeepCopyVirtualHost `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualHostList is the generic Kubernetes API object wrapper
type VirtualHostList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	metav1.Status     `json:"status,omitempty"`
	Items             []VirtualHost `json:"items"`
}

type Status string

type DeepCopyRoute v1.Route

func (in *DeepCopyRoute) DeepCopyInto(out *DeepCopyRoute) {
	data, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}

type DeepCopyUpstream v1.Upstream

func (in *DeepCopyUpstream) DeepCopyInto(out *DeepCopyUpstream) {
	data, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}

type DeepCopyVirtualHost v1.VirtualHost

func (in *DeepCopyVirtualHost) DeepCopyInto(out *DeepCopyVirtualHost) {
	data, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}
