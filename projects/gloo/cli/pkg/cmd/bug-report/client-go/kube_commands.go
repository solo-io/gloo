package client_go

// GetClusterContext returns the cluster context.
func GetClusterContext(kubeConfig string) (string, error) {
	return RunCmd("config current-context", "", kubeConfig, "", false)
}
