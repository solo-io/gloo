package apiserverutils

import "os"

const (
	InstallNamespaceKey     = "POD_NAMESPACE"
	DefaultInstallNamespace = "gloo-system"
	// currently the deployment name is not customizable
	DefaultDeploymentName = "gloo-fed"
)

// Returns the namespace that Gloo Fed apiserver was installed in
func GetInstallNamespace() string {
	if installNamespace := os.Getenv(InstallNamespaceKey); installNamespace != "" {
		return installNamespace
	}
	return DefaultInstallNamespace
}
