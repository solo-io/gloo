package surveyutils

import (
	"fmt"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
)

func AddVirtualServiceFlagsInteractive(opts *options.ExtraOptions) error {

	if err := rateLimitingSurvey(&opts.RateLimit); err != nil {
		return err
	}
	if err := oidcSurvey(&opts.OIDCAuth); err != nil {
		return err
	}
	if err := apiKeySurvey(&opts.ApiKeyAuth); err != nil {
		return err
	}

	return nil
}

// TODO: move this to input virtual host
func rateLimitingSurvey(input *options.RateLimit) error {
	yes, err := cliutil.GetYesInput("do you wish to add rate limiting to the virtual service [y/n]?")
	if err != nil {
		return err
	}

	if !yes {
		return nil
	}

	input.Enable = true

	if err := cliutil.ChooseFromList(
		"Limit requests over this unit of time for clients: ",
		&input.TimeUnit,
		options.RateLimit_TimeUnits,
	); err != nil {
		return err
	}

	if err := cliutil.GetUint32InputDefault(
		fmt.Sprintf("how many requests per %v?", input.TimeUnit),
		&input.RequestsPerTimeUnit,
		100,
	); err != nil {
		return err
	}
	return nil
}

func oidcSurvey(input *options.OIDCAuth) error {
	yes, err := cliutil.GetYesInput("do you wish to add oidc auth to the virtual service [y/n]?")
	if err != nil {
		return err
	}

	if !yes {
		return nil
	}

	input.Enable = true

	err = cliutil.GetStringInput("What is your app url?", &input.AppUrl)
	if err != nil {
		return err
	}
	err = cliutil.GetStringInput("What is your issuer url?", &input.IssuerUrl)
	if err != nil {
		return err
	}
	err = cliutil.GetStringInputDefault("What path (relative to your app url) should we use as a callback from the issuer?", &input.CallbackPath, "/oidc-gloo-callback")
	if err != nil {
		return err
	}
	err = cliutil.GetStringInput("What is your client id?", &input.ClientId)
	if err != nil {
		return err
	}
	err = cliutil.GetStringInput("What is your client secret name?", &input.ClientSecretRef.Name)
	if err != nil {
		return err
	}
	err = cliutil.GetStringInput("What is your client secret namespace?", &input.ClientSecretRef.Namespace)
	if err != nil {
		return err
	}

	return nil
}

func apiKeySurvey(input *options.ApiKeyAuth) error {
	yes, err := cliutil.GetYesInput("do you wish to add apikey auth to the virtual service [y/n]?")
	if err != nil {
		return err
	}

	if !yes {
		return nil
	}

	input.Enable = true

	err = cliutil.GetStringSliceInput("provide a label (key=value) to be a part of your label selector (empty to finish)", &input.Labels)
	if err != nil {
		return err
	}

	err = cliutil.GetStringInput("apikey secret name to attach to this virtual service? (empty to skip)", &input.SecretName)
	if err != nil {
		return err
	}

	if input.SecretName != "" {
		err = cliutil.GetStringInput("provide a namespace to search for the secret in (empty to finish)", &input.SecretNamespace)
		if err != nil {
			return err
		}
	}

	return nil
}
