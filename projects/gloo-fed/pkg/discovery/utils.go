package discovery

import (
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"k8s.io/apimachinery/pkg/api/errors"
)

func GetAdminProxyForInstance(instance *fedv1.GlooInstance) *fed_types.GlooInstanceSpec_Proxy {
	for _, proxy := range instance.Spec.GetProxies() {
		if proxy.GetName() == instance.Spec.GetAdmin().GetProxyId().GetName() &&
			proxy.GetNamespace() == instance.Spec.GetAdmin().GetProxyId().GetNamespace() {
			return proxy
		}
	}
	// This should never happen
	return nil
}

func IgnoreAlreadyExists(err error) error {
	if err != nil && errors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
