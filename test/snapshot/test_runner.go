package snapshot

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	"github.com/solo-io/gloo/test/kube2e"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	InsecureIngressClusterPort = 31080
	IngressPortClusterSSL      = 31443
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

	Client    client.Client
	ClientSet *kube2e.KubeResourceClientSet

	ToCleanup []client.Object // all objects written for an individual test run should be cleaned up at the end
}

func (tr TestRunner) Run(ctx context.Context, inputs []client.Object) error {
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
		// add to cleanup list
		tr.ToCleanup = append(tr.ToCleanup, obj)
	}
	return nil
}

func (tr TestRunner) RunFromFile(ctx context.Context, inputFiles []string) error {
	// load inputs
	var (
		gateways []*gwv1.Gateway
		inputs   []client.Object
	)
	for _, file := range inputFiles {
		objs, err := testutils.LoadFromFiles(ctx, file)
		if err != nil {
			return err
		}
		for _, obj := range objs {
			switch obj := obj.(type) {
			case *gwv1.Gateway:
				gateways = append(gateways, obj)
			default:
				inputs = append(inputs, obj)
			}
		}
	}

	return tr.Run(ctx, inputs)
}

func (tr *TestRunner) Cleanup(ctx context.Context) error {
	var errs error
	for _, obj := range tr.ToCleanup {
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
