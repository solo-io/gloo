package bootstrap

import "os"

const (
	InstallNamespaceKey     = "POD_NAMESPACE"
	DefaultInstallNamespace = "gloo-system"
)

// Returns the namespace that Gloo Fed was installed in
func GetInstallNamespace() string {
	if installNamespace := os.Getenv(InstallNamespaceKey); installNamespace != "" {
		return installNamespace
	}
	return DefaultInstallNamespace
}
