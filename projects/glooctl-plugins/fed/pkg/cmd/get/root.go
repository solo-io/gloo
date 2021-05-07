package get

import (
	enterprise_gloo_cli_client "github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/api/enterprise.gloo.solo.io/v1/cli"
	gateway_cli_client "github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/api/gateway.solo.io/v1/cli"
	gloo_cli_client "github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/api/gloo.solo.io/v1/cli"
	ratelimit_cli_client "github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/api/ratelimit.api.solo.io/v1alpha1/cli"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/cmd/options"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/constants"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/flagutils"
	"github.com/solo-io/solo-projects/projects/glooctl-plugins/fed/pkg/helpers"
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var UnsetNamespaceError = eris.New("Gloo Fed installation namespace does not exist. Did you install it in another namespace and forgot to add the '-n NAMESPACE' flag?")

func RootCmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or a list of Gloo Fed resources",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			client := helpers.MustKubeClient()
			_, err := client.CoreV1().Namespaces().Get(opts.Ctx, opts.Namespace, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					return UnsetNamespaceError
				}
				// we would still like to attempt the command even if we don't have rbac to list namespaces, so just log a warning
				contextutils.LoggerFrom(opts.Ctx).Warnf("unable to locate gloo fed installation namespace", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return constants.SubcommandError
		},
	}

	pflags := cmd.PersistentFlags()
	flagutils.AddApiserverFlags(pflags, opts)

	cmd.AddCommand(gloo_cli_client.Upstream(opts))
	cmd.AddCommand(gloo_cli_client.UpstreamGroup(opts))
	cmd.AddCommand(gloo_cli_client.Settings(opts))
	cmd.AddCommand(gloo_cli_client.Proxy(opts))
	cmd.AddCommand(gateway_cli_client.VirtualService(opts))
	cmd.AddCommand(gateway_cli_client.RouteTable(opts))
	cmd.AddCommand(gateway_cli_client.Gateway(opts))
	cmd.AddCommand(enterprise_gloo_cli_client.AuthConfig(opts))
	cmd.AddCommand(ratelimit_cli_client.RateLimitConfig(opts))

	return cmd

}
