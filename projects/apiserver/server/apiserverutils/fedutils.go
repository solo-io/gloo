package apiserverutils

import (
	"context"

	appsv1client "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Determines whether Gloo Fed is enabled in the current Gloo installation by checking for the existence
// of the Gloo Fed deployment.
func IsGlooFedEnabled(ctx context.Context, config *rest.Config) (bool, error) {
	kubeClient, err := client.New(config, client.Options{})
	if err != nil {
		return false, err
	}

	// look for the Gloo Fed deployment in the install namespace
	deploymentClient := appsv1client.NewDeploymentClient(kubeClient)
	_, err = deploymentClient.GetDeployment(ctx, client.ObjectKey{
		Namespace: GetInstallNamespace(),
		Name:      DefaultDeploymentName,
	})
	if apierrors.IsNotFound(err) {
		// could not find the deployment, so Gloo Fed is not enabled
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
