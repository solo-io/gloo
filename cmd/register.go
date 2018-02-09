package cmd

import (
	"log"

	"github.com/pkg/errors"
	kubev1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
	"github.com/spf13/cobra"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

// registerCmd is CLI command to register CRDs. This should only
// be necessary during development as Glue would do this in
// production environment
func registerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register CRDs (only needed for developing)",
		RunE: func(c *cobra.Command, args []string) error {
			cfg, err := getClientConfig()
			if err != nil {
				return errors.Wrap(err, "unable to get client configuration")
			}
			register(cfg)
			return nil
		},
	}
	return cmd
}

func register(cfg *rest.Config) {
	client, err := apiexts.NewForConfig(cfg)
	if err != nil {
		log.Printf("Error getting client to register CRDs %q\n", err)
		return
	}
	upstream := &v1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: "upstreams.glue.solo.io"},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   kubev1.GroupName,
			Version: kubev1.Version,
			Scope:   v1beta1.NamespaceScoped,
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural: "upstreams",
				Kind:   "Upstream",
			},
		},
	}
	if _, err = client.ApiextensionsV1beta1().
		CustomResourceDefinitions().
		Create(upstream); err != nil && !apierrors.IsAlreadyExists(err) {
		log.Printf("Unable to register CRDs %q\n", err)
	}
}
