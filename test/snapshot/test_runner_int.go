package snapshot

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/onsi/ginkgo/v2"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	. "github.com/solo-io/gloo/projects/gateway2/translator"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/test/kube2e"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
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
	Name             string
	ResultsByGateway map[types.NamespacedName]ExpectedTestResult

	Client    client.Client
	ClientSet *kube2e.KubeResourceClientSet

	ToCleanup []client.Object // all objects written for an individual test run should be cleaned up at the end
}

type ActualTestResult struct {
	Proxy *v1.Proxy
	// Reports     map[types.NamespacedName]*reports.GatewayReport
}

type ExpectedTestResult struct {
	Proxy string
	// Reports     map[types.NamespacedName]*reports.GatewayReport
}

func (r ExpectedTestResult) Equals(actual ActualTestResult) (bool, error) {
	proxy, err := testutils.ReadProxyFromFile(r.Proxy)
	if err != nil {
		return false, err
	}
	return proxy.Equal(actual.Proxy), nil
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

	}

	return nil
}

// map of gwv1.GW namespace/name to translation result
func (tr TestRunner) RunInMemory(ctx context.Context, inputs []client.Object) (map[types.NamespacedName]bool, error) {
	var gateways []*gwv1.Gateway
	for _, obj := range inputs {
		switch obj := obj.(type) {
		case *gwv1.Gateway:
			gateways = append(gateways, obj)
		}
	}

	queries := testutils.BuildGatewayQueries(inputs)

	results := make(map[types.NamespacedName]bool)
	for _, gw := range gateways {

		ref := types.NamespacedName{
			Namespace: gw.Namespace,
			Name:      gw.Name,
		}
		reportsMap := reports.NewReportMap()
		reporter := reports.NewReporter(&reportsMap)
		pluginRegistry := registry.NewPluginRegistry(queries)
		// translate gateway
		proxy := NewTranslator(*pluginRegistry).TranslateProxy(
			ctx,
			gw,
			queries,
			reporter,
		)

		act, _ := testutils.MarshalYaml(proxy)
		fmt.Fprintf(ginkgo.GinkgoWriter, "actual result:\n %s \n", act)

		actReport, _ := testutils.MarshalAnyYaml(reportsMap)
		fmt.Fprintf(ginkgo.GinkgoWriter, "actual reports:\n %s \n", actReport)

		actual := ActualTestResult{
			Proxy: proxy,
			// Reports:     reportsMap.Gateways,
		}

		expected, ok := tr.ResultsByGateway[ref]
		if !ok {
			return nil, errors.Errorf("no expected result found for gateway %v", ref)
		}

		equal, err := expected.Equals(actual)
		if err != nil {
			return nil, err
		}

		results[ref] = equal
	}

	return results, nil
}

func (tr TestRunner) RunFromFile(ctx context.Context, inputFiles []string) (map[types.NamespacedName]bool, error) {
	// load inputs
	var (
		gateways []*gwv1.Gateway
		inputs   []client.Object
	)
	for _, file := range inputFiles {
		objs, err := testutils.LoadFromFiles(ctx, file)
		if err != nil {
			return nil, err
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

	return tr.RunInMemory(ctx, inputs)
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
