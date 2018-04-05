package openfaas

import (
	"errors"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/internal/detector"
	"github.com/solo-io/gloo-plugins/rest"

	updatefaas "github.com/solo-io/gloo-function-discovery/internal/updater/openfaas"
)

type faasDetector struct {
}

func NewFaasDetector() detector.Interface {
	return &faasDetector{}
}

func (d *faasDetector) DetectFunctionalService(us *v1.Upstream, addr string) (*v1.ServiceInfo, map[string]string, error) {
	if updatefaas.IsOpenFaas(us) {
		svcInfo := &v1.ServiceInfo{
			Type: rest.ServiceTypeREST,
		}
		return svcInfo, nil, nil
	}

	return nil, nil, errors.New("not a faas upstream")
}
