package v1

import (
	"github.com/name5566/leaf/util"
	"github.com/solo-io/glue/pkg/api/types"
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
	Spec              DeepCopyRoute
	Status            `json:"status"`
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
	Spec              DeepCopyUpstream
	Status            `json:"status"`
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
	Spec              DeepCopyVirtualHost
	Status            `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualHostList is the generic Kubernetes API object wrapper
type VirtualHostList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []VirtualHost `json:"items"`
}

type Status int

const (
	Status_Pending Status = iota
	Status_Completed
	Status_Failed
)

func (s Status) String() string {
	switch s {
	case Status_Pending:
		return "Pending"
	case Status_Completed:
		return "Completed"
	case Status_Failed:
		return "Failed"
	}
	return "Unknown"
}

type DeepCopyRoute types.Route

func (r DeepCopyRoute) DeepCopyInto(r2 *DeepCopyRoute) {
	util.DeepCopy(&r, r2)
}

type DeepCopyUpstream types.Upstream

func (r DeepCopyUpstream) DeepCopyInto(r2 *DeepCopyUpstream) {
	util.DeepCopy(&r, r2)
}

type DeepCopyVirtualHost types.VirtualHost

func (r DeepCopyVirtualHost) DeepCopyInto(r2 *DeepCopyVirtualHost) {
	util.DeepCopy(&r, r2)
}
