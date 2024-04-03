package models

// ClusterContext contains the metadata about a Kubernetes cluster
// It also includes useful utilities for interacting with that cluster
type ClusterContext struct {
	// The name of the Kubernetes cluster
	Name string

	KubeContext string
}

/**

resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

clientScheme, err = utils.BuildClientScheme()
Expect(err).NotTo(HaveOccurred(), "can build client scheme")

kubeClient, err = utils.GetClient(kubeCtx, clientScheme)
Expect(err).NotTo(HaveOccurred(), "can create client")
*/
