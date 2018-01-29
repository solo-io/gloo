package controller

import (
	"fmt"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/solo-io/glue/config/watcher/crd/solo.io/v1"
	"github.com/solo-io/glue/pkg/log"
)

// register crds
func RegisterCrds(restConfig *rest.Config) error {
	clientset, err := apiexts.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to create api extension clientset: %v", err)
	}
	for _, crd := range v1.KnownCRDs {
		toRegister := &v1beta1.CustomResourceDefinition{
			ObjectMeta: meta_v1.ObjectMeta{Name: crd.FullName()},
			Spec: v1beta1.CustomResourceDefinitionSpec{
				Group:   crd.Group,
				Version: crd.Version,
				Scope:   v1beta1.NamespaceScoped,
				Names: v1beta1.CustomResourceDefinitionNames{
					Plural: crd.Plural,
					Kind:   crd.Kind,
				},
			},
		}
		log.Printf("registering crd %v", crd)
		if _, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(toRegister); err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create crd: %v", err)
		}
	}
	return nil
}
