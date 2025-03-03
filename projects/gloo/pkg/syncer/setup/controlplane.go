package setup

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
)

var (
	NoXdsPortFoundError = eris.New("failed to find xds port")
	noXdsPortFoundError = func(portName string, svcNamespace string, svcName string) error {
		return eris.Wrapf(NoXdsPortFoundError, "no port with the name %s found in service %s.%s", portName, svcNamespace, svcName)
	}
	NoGlooSvcFoundError = eris.New("failed to find Gloo service")
	noGlooSvcFoundError = func(svcNamespace string) error {
		return eris.Wrapf(NoGlooSvcFoundError, "service in %s with gloo=gloo", svcNamespace)
	}
	MultipleGlooSvcFoundError = eris.New("found multiple Gloo services")
	multipleGlooSvcFoundError = func(svcNamespace string) error {
		return eris.Wrapf(MultipleGlooSvcFoundError, "found multiple services in %s with gloo=gloo label", svcNamespace)
	}
)

func GetControlPlaneService(ctx context.Context, svcNamespace string, svcClient skkube.ServiceClient) (*skkube.Service, error) {
	opts := clients.ListOpts{
		Ctx:      ctx,
		Selector: kubeutils.GlooServiceLabels,
	}
	services, err := svcClient.List(svcNamespace, opts)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, noGlooSvcFoundError(svcNamespace)
	}

	if len(services) > 1 {
		return nil, multipleGlooSvcFoundError(svcNamespace)
	}

	return services[0], nil
}

// GetControlPlaneXdsPort gets the xDS port from the Gloo Service
func GetControlPlaneXdsPort(service *skkube.Service) (int32, error) {
	// find the xds port on the Gloo Service
	for _, port := range service.Spec.Ports {
		if port.Name == kubeutils.GlooXdsPortName {
			return port.Port, nil
		}
	}
	return 0, noXdsPortFoundError(kubeutils.GlooXdsPortName, service.Namespace, service.Name)
}

// GetControlPlaneXdsHost gets the xDS address from the gloo Service.
func GetControlPlaneXdsHost(service *skkube.Service) string {
	return kubeutils.ServiceFQDN(metav1.ObjectMeta{
		Name:      service.Name,
		Namespace: service.Namespace,
	})
}
