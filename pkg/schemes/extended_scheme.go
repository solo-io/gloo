package schemes

import (
	"fmt"

	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"

	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// AddGatewayV1A2Scheme adds the gwv1a2 scheme to the provided scheme if the TCPRoute CRD exists.
func AddGatewayV1A2Scheme(restConfig *rest.Config, scheme *runtime.Scheme) error {
	exists, err := CRDExists(restConfig, gwv1a2.GroupVersion.Group, gwv1a2.GroupVersion.Version, wellknown.TCPRouteKind)
	if err != nil {
		return fmt.Errorf("error checking if %s CRD exists: %w", wellknown.TCPRouteKind, err)
	}

	if exists {
		if err := gwv1a2.Install(scheme); err != nil {
			return fmt.Errorf("error adding Gateway API v1alpha2 to scheme: %w", err)
		}
	}

	return nil
}

// Helper function to check if a CRD exists
func CRDExists(restConfig *rest.Config, group, version, kind string) (bool, error) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return false, err
	}

	groupVersion := fmt.Sprintf("%s/%s", group, version)
	apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		if errors.IsNotFound(err) || discovery.IsGroupDiscoveryFailedError(err) || meta.IsNoMatchError(err) {
			return false, nil
		}
		return false, err
	}

	for _, apiResource := range apiResourceList.APIResources {
		if apiResource.Kind == kind {
			return true, nil
		}
	}

	return false, nil
}
