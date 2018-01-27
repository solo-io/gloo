package crd

import (
	"fmt"
	"reflect"

	"github.com/solo-io/delta-configurator/test/configurator/crd/delta.io/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createCRD(clientSet apiexts.Interface) error {
	crd := &v1beta1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{Name: v1.CRDFullName},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   v1.CRDGroup,
			Version: v1.CRDVersion,
			Scope:   v1beta1.NamespaceScoped,
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: v1.CRDPlural,
				Kind:   reflect.TypeOf(v1.ApiDefinition{}).Name(),
			},
		},
	}
	if _, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd); err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create crd: %v", err)
	}
	return nil
}
