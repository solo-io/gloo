package kubegateway

import (
	"context"
	"fmt"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"

	"github.com/solo-io/gloo/projects/gateway2/crds"
	"github.com/solo-io/gloo/projects/gateway2/deployer"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func uninstall(opts *options.Options, installOpts *Options) error {
	ctx := context.Background()

	cfg, err := kubeutils.GetRestConfigWithKubeContext(opts.Top.KubeContext)
	if err != nil {
		return err
	}

	cli, err := client.New(cfg, client.Options{})
	if err != nil {
		return err
	}

	dep, err := deployer.NewDeployer(cli, &deployer.Inputs{
		ControllerName: "glooctl",
	}, nil)
	if err != nil {
		return err
	}

	objs, err := dep.Render("default", installOpts.Namespace, map[string]any{})
	if err != nil {
		return err
	}

	for _, obj := range objs {
		obj.SetNamespace(installOpts.Namespace)
	}

	fmt.Printf("Deleting Manifest... ")
	if err := deleteObjs(ctx, objs, cli); err != nil {
		fmt.Printf("Failed, but continuing\n")
	} else {
		fmt.Printf("Done\n")
	}

	fmt.Printf("Deleting Gateway CRDs... ")
	crds, err := deployer.ConvertYAMLToObjects(cli.Scheme(), crds.GatewayCrds)
	if err != nil {
		fmt.Printf("Failed\n")
	} else {
		if err := dep.DeployObjs(ctx, crds); err != nil {
			fmt.Printf("Failed\n")
		}
		fmt.Printf("Done\n")
	}

	deleteNamespace(ctx, cli, installOpts.Namespace)

	return nil
}

func deleteObjs(ctx context.Context, objs []client.Object, cli client.Client) error {
	for _, obj := range objs {
		if err := cli.Delete(ctx, obj); err != nil {
			return fmt.Errorf("failed to delete object %s %s: %w", obj.GetObjectKind().GroupVersionKind().String(), obj.GetName(), err)
		}
	}
	return nil
}

func deleteNamespace(ctx context.Context, cli client.Client, namespace string) {
	ns := corev1.Namespace{}
	err := cli.Get(ctx, client.ObjectKey{Name: namespace}, &ns)
	if err != nil {
		if apierrors.IsNotFound(err) {
			fmt.Printf("Nmespace %s has already been deleted... ", namespace)
		} else {
			fmt.Printf("\nUnable to check if namespace %s exists. Continuing...\n", namespace)
		}
	}
	fmt.Printf("Deleting namespace %s... ", namespace)
	if err := cli.Delete(ctx, &ns); err != nil {
		fmt.Printf("Failed\n")
	}
	fmt.Printf("Done\n")
}
