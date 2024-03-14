package translator_test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/onsi/ginkgo/v2"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	. "github.com/solo-io/gloo/projects/gateway2/translator"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"google.golang.org/protobuf/testing/protocmp"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type TestCase struct {
	Name             string
	InputFiles       []string
	ResultsByGateway map[types.NamespacedName]ExpectedTestResult
}

type ActualTestResult struct {
	Proxy *v1.Proxy
	// Reports     map[types.NamespacedName]*reports.GatewayReport
	// TODO(Law): figure out how RouteReports fit in
}

type ExpectedTestResult struct {
	ProxyFile string
	// Reports     map[types.NamespacedName]*reports.GatewayReport
}

var cmpOpts = []cmp.Option{protocmp.Transform(), cmpopts.EquateEmpty()}

func (r ExpectedTestResult) Cmp(actual ActualTestResult) (string, error) {
	proxy, err := r.Proxy()
	if err != nil {
		return "", err
	}
	return cmp.Diff(proxy, actual.Proxy, cmpOpts...), nil
}

func (r ExpectedTestResult) Proxy() (*v1.Proxy, error) {
	proxy, err := testutils.ReadProxyFromFile(r.ProxyFile)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}

// map of gwv1.GW namespace/name to translation result
func (tc TestCase) Run(ctx context.Context) (map[types.NamespacedName]string, error) {
	// load inputs

	var (
		gateways     []*gwv1.Gateway
		dependencies []client.Object
	)
	for _, file := range tc.InputFiles {
		objs, err := testutils.LoadFromFiles(ctx, file)
		if err != nil {
			return nil, err
		}
		for _, obj := range objs {
			switch obj := obj.(type) {
			case *gwv1.Gateway:
				gateways = append(gateways, obj)
			default:
				dependencies = append(dependencies, obj)
			}
		}
	}

	queries := testutils.BuildGatewayQueries(dependencies)
	pluginRegistry := registry.NewPluginRegistry(registry.BuildPlugins(queries))

	diffs := make(map[types.NamespacedName]string)
	for _, gw := range gateways {

		ref := types.NamespacedName{
			Namespace: gw.Namespace,
			Name:      gw.Name,
		}
		reportsMap := reports.NewReportMap()
		reporter := reports.NewReporter(&reportsMap)

		// translate gateway
		proxy := NewTranslator(queries, pluginRegistry).TranslateProxy(
			ctx,
			gw,
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

		expected, ok := tc.ResultsByGateway[ref]
		if !ok {
			return nil, errors.Errorf("no expected result found for gateway %v", ref)
		}

		diff, err := expected.Cmp(actual)
		if err != nil {
			return nil, err
		}
		diffs[ref] = diff
	}

	return diffs, nil
}
