package create

import (
	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/argsutils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/constants"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/printers"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/surveyutils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	optionsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	flagutilsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/flagutils"
	surveyutilsExt "github.com/solo-io/solo-projects/projects/gloo/cli/pkg/surveyutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	extauth2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	ratelimit2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	"github.com/spf13/cobra"
)

var defaultDomains = []string{"*"}

var (
	ProvideNamespaceAndNameError = func(namespace, secretName string) error {
		return errors.Errorf("provide both a secret namespace [%v] and secret name [%v]", namespace, secretName)
	}
	UnableToMarshalApiKeyConfig = func(err error) error {
		return errors.Wrapf(err, "Error marshalling apikey config")
	}
)

func VSCreate(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {

	optsExt := &optionsExt.ExtraOptions{}
	optsExt.OIDCAuth.ClientSecretRef = new(core.ResourceRef)

	cmd := &cobra.Command{
		// Use command constants to aid with replacement.
		Use:     constants.VIRTUAL_SERVICE_COMMAND.Use,
		Aliases: constants.VIRTUAL_SERVICE_COMMAND.Aliases,
		Short:   "Create a Virtual Service",
		Long: "A virtual service describes the set of routes to match for a set of domains. \n" +
			"Virtual services are containers for routes assigned to a domain or set of domains. \n" +
			"Virtual services must not have overlapping domains, as the virtual service to match a request " +
			"is selected by the Host header (in HTTP1) or :authority header (in HTTP2). " +
			"When using Gloo Enterprise, virtual services can be configured with rate limiting, oauth, and apikey auth.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Top.Interactive {
				if err := surveyutils.AddVirtualServiceFlagsInteractive(&opts.Create.VirtualService); err != nil {
					return err
				}
				if err := surveyutilsExt.AddVirtualServiceFlagsInteractive(optsExt); err != nil {
					return err
				}
			}
			err := argsutils.MetadataArgsParse(opts, args)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return createVirtualService(opts, optsExt, args)
		},
	}

	pflags := cmd.PersistentFlags()
	flagutils.AddMetadataFlags(pflags, &opts.Metadata)
	flagutils.AddVirtualServiceFlags(pflags, &opts.Create.VirtualService)
	flagutilsExt.AddVirtualServiceFlags(pflags, optsExt)
	cliutils.ApplyOptions(cmd, optionsFunc)

	return cmd
}

func createVirtualService(opts *options.Options, optsExt *optionsExt.ExtraOptions, args []string) error {
	vs, err := virtualServiceFromOpts(opts.Metadata, opts.Create.VirtualService, *optsExt)
	if err != nil {
		return err
	}

	if !opts.Create.DryRun {
		virtualServiceClient := helpers.MustVirtualServiceClient()
		vs, err = virtualServiceClient.Write(vs, clients.WriteOpts{})
		if err != nil {
			return err
		}
	}

	printers.PrintVirtualServices(v1.VirtualServiceList{vs}, opts.Top.Output)

	return nil
}

// TODO: dedupe with Gloo
func virtualServiceFromOpts(meta core.Metadata, input options.InputVirtualService, extopts optionsExt.ExtraOptions) (*v1.VirtualService, error) {
	if len(input.Domains) == 0 {
		input.Domains = defaultDomains
	}
	displayName := meta.Name
	if input.DisplayName != "" {
		displayName = input.DisplayName
	}
	vs := &v1.VirtualService{
		Metadata: meta,
		VirtualHost: &gloov1.VirtualHost{
			Domains: input.Domains,
		},
		DisplayName: displayName,
	}
	rl := extopts.RateLimit
	if rl.Enable {
		if vs.VirtualHost.VirtualHostPlugins == nil {
			vs.VirtualHost.VirtualHostPlugins = &gloov1.VirtualHostPlugins{}
		}
		timeUnit, ok := ratelimit.RateLimit_Unit_value[rl.TimeUnit]
		if !ok {
			return nil, errors.Errorf("invalid time unit specified: %v", rl.TimeUnit)
		}
		ingressRateLimit := &ratelimit.IngressRateLimit{
			AnonymousLimits: &ratelimit.RateLimit{
				Unit:            ratelimit.RateLimit_Unit(timeUnit),
				RequestsPerUnit: rl.RequestsPerTimeUnit,
			},
		}
		ingressRateLimitStruct, err := envoyutil.MessageToStruct(ingressRateLimit)
		if err != nil {
			return nil, errors.Wrapf(err, "Error marshalling ingress rate limit")
		}
		if vs.VirtualHost.VirtualHostPlugins.Extensions == nil {
			vs.VirtualHost.VirtualHostPlugins.Extensions = new(gloov1.Extensions)
		}
		if vs.VirtualHost.VirtualHostPlugins.Extensions.Configs == nil {
			vs.VirtualHost.VirtualHostPlugins.Extensions.Configs = make(map[string]*types.Struct)
		}
		vs.VirtualHost.VirtualHostPlugins.Extensions.Configs[ratelimit2.ExtensionName] = ingressRateLimitStruct
	}

	oidc := extopts.OIDCAuth
	if oidc.Enable {
		if vs.VirtualHost.VirtualHostPlugins == nil {
			vs.VirtualHost.VirtualHostPlugins = &gloov1.VirtualHostPlugins{}
		}
		if oidc.AppUrl == "" {
			return nil, errors.Errorf("invalid app url specified: %v", oidc.AppUrl)
		}
		if oidc.IssuerUrl == "" {
			return nil, errors.Errorf("invalid issuer url specified: %v", oidc.IssuerUrl)
		}
		if oidc.ClientId == "" {
			return nil, errors.Errorf("invalid client id specified: %v", oidc.ClientId)
		}
		if oidc.CallbackPath == "" {
			return nil, errors.Errorf("invalid callback path specified: %v", oidc.CallbackPath)
		}
		if oidc.ClientSecretRef.Name == "" || oidc.ClientSecretRef.Namespace == "" {
			return nil, errors.Errorf("invalid client secret ref specified: %v.%v", oidc.ClientSecretRef.Namespace, oidc.ClientSecretRef.Name)
		}
		vhostAuth := &extauth.VhostExtension{
			AuthConfig: &extauth.VhostExtension_Oauth{
				Oauth: &extauth.OAuth{
					AppUrl:          oidc.AppUrl,
					CallbackPath:    oidc.CallbackPath,
					ClientId:        oidc.ClientId,
					ClientSecretRef: oidc.ClientSecretRef,
					IssuerUrl:       oidc.IssuerUrl,
				},
			},
		}
		vhostAuthStruct, err := envoyutil.MessageToStruct(vhostAuth)
		if err != nil {
			return nil, errors.Wrapf(err, "Error marshalling oauth config")
		}
		if vs.VirtualHost.VirtualHostPlugins.Extensions == nil {
			vs.VirtualHost.VirtualHostPlugins.Extensions = new(gloov1.Extensions)
		}
		if vs.VirtualHost.VirtualHostPlugins.Extensions.Configs == nil {
			vs.VirtualHost.VirtualHostPlugins.Extensions.Configs = make(map[string]*types.Struct)
		}
		vs.VirtualHost.VirtualHostPlugins.Extensions.Configs[extauth2.ExtensionName] = vhostAuthStruct

	}

	apiKey := extopts.ApiKeyAuth
	if apiKey.Enable {
		if vs.VirtualHost.VirtualHostPlugins == nil {
			vs.VirtualHost.VirtualHostPlugins = &gloov1.VirtualHostPlugins{}
		}

		var secretRefs []*core.ResourceRef
		if apiKey.SecretNamespace != "" && apiKey.SecretName != "" {
			secretRefs = []*core.ResourceRef{
				{
					Namespace: apiKey.SecretNamespace,
					Name:      apiKey.SecretName,
				},
			}
		} else if apiKey.SecretNamespace != "" || apiKey.SecretName != "" {
			return nil, ProvideNamespaceAndNameError(apiKey.SecretNamespace, apiKey.SecretName)
		}

		var labels options.InputMapStringString
		labels.Entries = apiKey.Labels
		var labelSelector map[string]string
		if len(labels.MustMap()) > 0 {
			labelSelector = labels.MustMap()
		}

		vhostAuth := &extauth.VhostExtension{
			AuthConfig: &extauth.VhostExtension_ApiKeyAuth{
				ApiKeyAuth: &extauth.ApiKeyAuth{
					LabelSelector:    labelSelector,
					ApiKeySecretRefs: secretRefs,
				},
			},
		}
		vhostAuthStruct, err := envoyutil.MessageToStruct(vhostAuth)
		if err != nil {
			return nil, UnableToMarshalApiKeyConfig(err)
		}
		if vs.VirtualHost.VirtualHostPlugins.Extensions == nil {
			vs.VirtualHost.VirtualHostPlugins.Extensions = new(gloov1.Extensions)
		}
		if vs.VirtualHost.VirtualHostPlugins.Extensions.Configs == nil {
			vs.VirtualHost.VirtualHostPlugins.Extensions.Configs = make(map[string]*types.Struct)
		}
		vs.VirtualHost.VirtualHostPlugins.Extensions.Configs[extauth2.ExtensionName] = vhostAuthStruct
	}

	return vs, nil
}
