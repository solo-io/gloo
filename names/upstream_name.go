package names

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func UpstreamName(serviceNamespace, serviceName string, servicePort intstr.IntOrString) string {
	return fmt.Sprintf("%s-%s-%v", serviceNamespace, serviceName, servicePort.String())
}
