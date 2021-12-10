package utils

const (

	// KubeSystemNamespace is the system namespace where we place kubernetes system components.
	KubeSystemNamespace string = "kube-system"

	// KubePublicNamespace is the namespace where we place kubernetes public info (ConfigMaps).
	KubePublicNamespace string = "kube-public"

	// KubeNodeLeaseNamespace is the namespace for the lease objects associated with each kubernetes node.
	KubeNodeLeaseNamespace string = "kube-node-lease"

	// LocalPathStorageNamespace is the namespace for dynamically provisioning persistent local storage with
	// Kubernetes. Typically used with the Kind cluster: https://github.com/rancher/local-path-provisioner
	LocalPathStorageNamespace string = "local-path-storage"
)

var (
	systemNamespaces = map[string]bool{
		KubeNodeLeaseNamespace:    true,
		LocalPathStorageNamespace: true,
		KubePublicNamespace:       true,
		KubeSystemNamespace:       true,
	}
)

func IsSystemNamespace(ns string) bool {
	return systemNamespaces[ns]
}
