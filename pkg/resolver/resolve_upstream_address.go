package resolver

import (
	"k8s.io/client-go/kubernetes"

	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	kubeplugin "github.com/solo-io/gloo-plugins/kubernetes"
	serviceplugin "github.com/solo-io/gloo/pkg/coreplugins/service"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resolver struct {
	Kube kubernetes.Interface
}

func (r *Resolver) Resolve(us *v1.Upstream) (string, error) {
	switch us.Type {
	case kubeplugin.UpstreamTypeKube:
		return resolveKubeUpstream(r.Kube, us)
	case serviceplugin.UpstreamTypeService:
		return resolveServiceUpstream(us)
	}
	// ignore other upstream types
	return "", nil
}

func resolveKubeUpstream(kube kubernetes.Interface, us *v1.Upstream) (string, error) {
	if kube == nil {
		return "", errors.New("function discovery not enabled for kubernetes upstreams. function discovery must be run in-cluster for this feature")
	}
	spec, err := kubeplugin.DecodeUpstreamSpec(us.Spec)
	if err != nil {
		return "", errors.Wrap(err, "parsing kubernetes upstream spec")
	}
	// port specified, we're good
	if spec.ServicePort != 0 {
		return fmt.Sprintf("%v.%v.svc.cluster.local:%v", spec.ServiceName, spec.ServiceNamespace, spec.ServicePort), nil
	}
	// look up port. must be single-port service
	svc, err := kube.CoreV1().Services(spec.ServiceNamespace).Get(spec.ServiceName, metav1.GetOptions{})
	if err != nil {
		return "", errors.Wrap(err, "failed to get k8s service for upstream")
	}
	if len(svc.Spec.Ports) != 1 {
		return "", errors.New("service_port must be specified for k8s services with more than one port")
	}
	port := svc.Spec.Ports[0].Port
	return fmt.Sprintf("%v.%v.svc.cluster.local:%v", spec.ServiceName, spec.ServiceNamespace, port), nil
}

func resolveServiceUpstream(us *v1.Upstream) (string, error) {
	spec, err := serviceplugin.DecodeUpstreamSpec(us.Spec)
	if err != nil {
		return "", errors.Wrap(err, "parsing service upstream spec")
	}
	if len(spec.Hosts) < 1 {
		return "", errors.New("at least one host required to resolve service upstream")
	}
	return fmt.Sprintf("%v:%v", spec.Hosts[0].Addr, spec.Hosts[0].Port), nil
}
