package license

import (
	"context"
	"os"
	"time"

	"github.com/rotisserie/eris"
	"github.com/solo-io/licensing/pkg/client"
	"github.com/solo-io/licensing/pkg/model"
)

const (
	EnvName = "GLOO_LICENSE_KEY"

	ParseError     = "License is invalid: Could not parse license"
	InvalidLicense = "License is invalid"
	ExpiredLicense = "License expired, please contact support to renew."
	ExpiredGraphql = "License does not support GraphQL"
	MissingGloo    = "No license for Gloo Edge found"
	MissingGraphql = "License does not support GraphQL"
)

// LicensedFeature represents the set of features that Gloo Edge Enterprise supports
type LicensedFeature int

const (
	Enterprise LicensedFeature = iota
	GraphQL
)

// LicensedFeatureProvider decides whether a provided license supports a set of Edge features
// The purpose of this LicensedFeatureProvider is to:
//
//	(A) codify how the state of that license affects the behavior of our application.
//	(B) decouple application logic from license logic. Now our application can be aware of
//		features, and unaware of the state of the license
type LicensedFeatureProvider struct {
	license *ValidatedLicense

	// memoize the state of enabled features to avoid re-processing
	stateByLicensedFeature map[LicensedFeature]*FeatureState
}

// A ValidatedLicense aggregates the state of the license and any errors or warnings on that license
type ValidatedLicense struct {
	*model.License
	Warn, Err error
}

type FeatureState struct {
	Enabled bool
	Reason  string
}

func NewLicensedFeatureProvider() *LicensedFeatureProvider {
	return &LicensedFeatureProvider{
		stateByLicensedFeature: map[LicensedFeature]*FeatureState{},
	}
}

func (l *LicensedFeatureProvider) validateLicense(ctx context.Context) (*model.License, error, error) {
	// If a license is provided, validate the key and configure features based on the license state
	if _, err := client.ParseLicenseInfo(ctx, os.Getenv(EnvName)); err != nil {
		return nil, nil, eris.New(ParseError)
	}
	licenseClient, err := client.NewLicensingClient(ctx, "", EnvName, nil)
	if err != nil {
		return nil, nil, eris.Wrapf(err, "could not initialize licensing client")
	}
	foundLicenseForState := func(state client.LicenseState, products ...model.Product) *model.License {
		for _, product := range products {
			licenseState, license := licenseClient.GetLicense(product)
			if licenseState == state {
				return license
			}
		}
		return nil
	}
	mainGlooProducts := []model.Product{model.Product_Gloo, model.Product_GlooTrial}
	license := foundLicenseForState(client.LicenseStateOk, mainGlooProducts...)
	if license != nil {
		return license, nil, nil
	}
	license = foundLicenseForState(client.LicenseStateExpired, mainGlooProducts...)
	if license != nil {
		return license, eris.New(ExpiredLicense), nil
	}
	license = foundLicenseForState(client.LicenseStateInvalid, mainGlooProducts...)
	if license != nil {
		return license, nil, eris.New(InvalidLicense)
	}
	return nil, nil, eris.New(MissingGloo)
}

// ValidateAndSetLicense sets l.license based off of the license key provided in the GLOO_LICENSE_KEY env var.
// If the license key is unparsable/ empty, an error is emitted. Otherwise, we look for the most valid license state for
// each of these products in the following order: gloo, gloo-trial. If the license key is corrupted (i.e. validation
// errors exist, it will be treated as invalid. If the license key is expired, a warning will be emitted with the date
// of expiry. If a license is not found or the provided license isnt for gloo or gloo-trial (i.e. for gloo-mesh)), an
// error will be emitted.
func (l *LicensedFeatureProvider) ValidateAndSetLicense(ctx context.Context) {
	if os.Getenv(EnvName) == "" {
		// If a license is not provided, mark it as nil, which will be treated as disabling enterprise features
		l.SetValidatedLicense(nil)
		return
	}
	license, warn, err := l.validateLicense(ctx)
	l.SetValidatedLicense(&ValidatedLicense{
		License: license,
		Warn:    warn,
		Err:     err,
	})
}

func (l *LicensedFeatureProvider) SetValidatedLicense(license *ValidatedLicense) {
	l.license = license

	l.setFeatureStateForEnterprise(license)
	l.setFeatureStateForGraphql(license)
}

func (l *LicensedFeatureProvider) setFeatureStateForEnterprise(license *ValidatedLicense) {
	featureState := &FeatureState{
		Enabled: false,
		Reason:  "License not found, Enterprise features not included",
	}
	l.stateByLicensedFeature[Enterprise] = featureState

	if license == nil {
		return
	}

	if license.Err != nil {
		// Error on the license means that the license is invalid
		// There are some ongoing decisions around how to properly validate Enterprise keys and expose features
		// To avoid crashing the proxy or control-plane.
		// Until the decision is made, we will keep the original behavior and enable enterprise features
		// for both expired and invalid licenses.
		// https://github.com/solo-io/solo-projects/issues/2918
		featureState.Enabled = true
		featureState.Reason = license.Err.Error()
		return
	}

	// Enterprise features are enabled, even if there is a warning on the license
	featureState.Enabled = true
	featureState.Reason = ""
	if license.Warn != nil {
		featureState.Reason = license.Warn.Error()
	}
}

func (l *LicensedFeatureProvider) setFeatureStateForGraphql(license *ValidatedLicense) {
	featureState := &FeatureState{
		Enabled: false,
		Reason:  "License not found, GraphQL features not included",
	}
	l.stateByLicensedFeature[GraphQL] = featureState

	if license == nil {
		return
	}

	if license.Err != nil {
		// Error on the license means that the license is invalid
		// There are some ongoing decisions around how to properly validate Enterprise keys and expose features
		// To avoid crashing the proxy or control-plane.
		// Until the decision is made, we will keep the original behavior and enable enterprise features
		// for both expired and invalid licenses.
		// https://github.com/solo-io/solo-projects/issues/2918
		featureState.Enabled = true
		featureState.Reason = license.Err.Error()
		return
	}
	foundGraphql := false
	graphqlExpiration := time.Time{}
	for _, addOn := range license.License.AddOns {
		if addOn.AddOn == model.GraphQL {
			graphqlExpiration = addOn.ExpiresAt
			foundGraphql = true
		}
	}

	if license.License == nil || !foundGraphql {
		// If GraphQL is not explicitly enabled, disable the feature
		featureState.Reason = MissingGraphql
		return
	}

	featureState.Enabled = true
	featureState.Reason = ""
	if graphqlExpiration.Equal(time.Time{}) {
		featureState.Reason = "GraphQL plugin license doesnt include expiration date, please contact support for an up to date license key."
	} else if foundGraphql && time.Now().After(graphqlExpiration) {
		featureState.Reason = ExpiredGraphql
	}

	if license.Warn != nil {
		featureState.Reason = license.Warn.Error()
	}
}

func (l *LicensedFeatureProvider) GetStateForLicensedFeature(feature LicensedFeature) *FeatureState {
	featureState, ok := l.stateByLicensedFeature[feature]
	if !ok {
		// This should never happen
		// We would only encounter this if we forgot to load the license
		return &FeatureState{
			Enabled: false,
			Reason:  "License information not available, treating as invalid",
		}
	}
	return featureState
}
