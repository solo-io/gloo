package snapshot

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	"github.com/solo-io/gloo/test/kube2e"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TestEnv struct {
	GatewayName      string
	GatewayNamespace string
	GatewayPort      int

	ClusterName    string
	ClusterContext string
}

type TestRunner struct {
	Name string

	Scheme    *runtime.Scheme
	Client    client.Client
	ClientSet *kube2e.KubeResourceClientSet

	Inputs    []client.Object // all objects written for an individual test run should be cleaned up at the end
	InputFile string
}

func (tr TestRunner) Run(ctx context.Context) error {
	if tr.Inputs != nil {
		return tr.run(ctx, tr.Inputs)
	} else if tr.InputFile != "" {
		return tr.runFromFile(ctx, []string{tr.InputFile})
	} else {
		return fmt.Errorf("no inputs provided")
	}
}

func (tr TestRunner) run(ctx context.Context, inputs []client.Object) error {
	for _, obj := range inputs {
		err := tr.Client.Create(ctx, obj, &client.CreateOptions{})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				// ignore already exists from previous test runs
				fmt.Fprintf(ginkgo.GinkgoWriter, "Object %s.%s already exists: %v\n", obj.GetName(), obj.GetNamespace(), err)
				continue
			}
			return err
		}
	}
	return nil
}

func (tr TestRunner) runFromFile(ctx context.Context, inputFiles []string) error {
	// load inputs
	var inputs []client.Object

	for _, file := range inputFiles {
		objs, err := testutils.LoadFromFiles(ctx, file, tr.Scheme)
		if err != nil {
			return err
		}
		for _, obj := range objs {
			inputs = append(inputs, obj)
		}
	}

	// set inputs to clean up after test run
	tr.Inputs = inputs

	return tr.run(ctx, inputs)
}

func (tr *TestRunner) Cleanup(ctx context.Context) error {
	var errs error
	for _, obj := range tr.Inputs {
		if obj == nil {
			continue
		}
		if err := tr.Client.Delete(ctx, obj); err != nil {
			if apierrors.IsNotFound(err) {
				fmt.Printf("warning to devs! resource deleted multiple times; this is likely a bug %s.%s", obj.GetName(), obj.GetNamespace())
				continue
			}
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}
