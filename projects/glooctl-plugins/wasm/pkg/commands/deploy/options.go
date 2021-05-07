package deploy

import (
	"context"

	"github.com/solo-io/k8s-utils/kubeutils"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	"github.com/solo-io/wasm-kit/pkg/commands/opts"
	"github.com/solo-io/wasm-kit/pkg/deploy"
	"github.com/solo-io/wasm-kit/pkg/deploy/gloo"
	v1 "github.com/solo-io/wasm-kit/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasm-kit/pkg/pull"
	"github.com/solo-io/wasm-kit/pkg/resolver"
	"github.com/spf13/pflag"
)

type options struct {
	// filter to deploy
	filter v1.FilterSpec

	// configuration string for filter
	filterConfig string

	// deployment implementation
	providerOptions

	// login
	opts.AuthOptions

	// remove a deployed filter instead of deploying
	remove bool
}

func (opts *options) addToFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&opts.filterConfig, "config", "", "", "optional config that will be passed to the filter. accepts an inline string.")
	flags.StringVarP(&opts.filter.RootID, "root-id", "", "", "optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.")
	opts.addIdToFlags(flags)
}

func (opts *options) addIdToFlags(flags *pflag.FlagSet) {
	flags.StringVar(&opts.filter.Id, "id", "", "unique id for naming the deployed filter. This is used for logging as well as removing the filter.")
}

type providerOptions struct {
	providerType string

	glooOpts glooOpts
}

type glooOpts struct {
	selector gloo.Selector
}

func (opts *glooOpts) addToFlags(flags *pflag.FlagSet) {
	flags.StringSliceVarP(&opts.selector.Namespaces, "namespaces", "n", nil, "deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.")
	flags.StringToStringVarP(&opts.selector.GatewayLabels, "labels", "l", nil, "select deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.")
}

const (
	Provider_Gloo = "gloo"
)

var SupportedProviders = []string{
	Provider_Gloo,
}

const (
	WorkloadType_DaemonSet   = "daemonset"
	WorkloadType_Deployment  = "deployment"
	WorkloadType_Statefulset = "statefulset"
)

var SupportedWorkloadTypes = []string{
	WorkloadType_DaemonSet,
	WorkloadType_Deployment,
	WorkloadType_Statefulset,
}

func (opts *options) makeProvider(ctx context.Context) (deploy.Provider, error) {
	switch opts.providerType {
	case Provider_Gloo:
		var gwClient gatewayv1.GatewayClient

		cfg, err := kubeutils.GetConfig("", "")
		kubeClientSet, err := gatewayv1.NewClientsetFromConfig(cfg)
		if err != nil {
			return nil, err
		}
		gwClient = kubeClientSet.Gateways()

		return &gloo.Provider{
			Ctx:           ctx,
			GatewayClient: gwClient,
			Selector:      opts.glooOpts.selector,
		}, nil
	}

	return nil, nil
}

func makeDeployer(ctx context.Context, opts *options) (*deploy.Deployer, error) {
	resolver, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.CredentialsFiles...)
	puller := pull.NewPuller(resolver)

	provider, err := opts.makeProvider(ctx)
	if err != nil {
		return nil, err
	}
	return &deploy.Deployer{
		Ctx:      ctx,
		Puller:   puller,
		Provider: provider,
	}, nil
}
