package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionSecretFunc func(original, desired *Secret)

type SecretReconciler interface {
	Reconcile(namespace string, desiredResources []*Secret, opts clients.ListOpts) error
}

func secretsToResources(list SecretList) []resources.Resource {
	var resourceList []resources.Resource
	for _, secret := range list {
		resourceList = append(resourceList, secret)
	}
	return resourceList
}

func NewSecretReconciler(client SecretClient, transition TransitionSecretFunc) SecretReconciler {
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) {
			transition(original.(*Secret), desired.(*Secret))
		}
	}
	return &secretReconciler{
		base: reconcile.NewReconciler(client.BaseClient(), transitionResources),
	}
}

type secretReconciler struct {
	base reconcile.Reconciler
}

func (r *secretReconciler) Reconcile(namespace string, desiredResources []*Secret, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "secret_reconciler")
	return r.base.Reconcile(namespace, secretsToResources(desiredResources), opts)
}
