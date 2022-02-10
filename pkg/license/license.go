package license

import (
	"context"

	"github.com/solo-io/licensing/pkg/model"
	"github.com/solo-io/licensing/pkg/validate"
)

const (
	EnvName = "GLOO_LICENSE_KEY"
)

// LicensedFeature represents the set of features that Gloo Edge Enterprise supports
type LicensedFeature int

const (
	Enterprise LicensedFeature = iota
	GraphQL
)

// LicensedFeatureProvider decides whether a provided license supports a set of Edge features
// The purpose of this LicensedFeatureProvider is to:
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

func (l *LicensedFeatureProvider) ValidateAndSetLicense(licenseString string) {
	if licenseString == "" {
		// If a license is not provided, mark it as nil, which will be treated as disabling enterprise features
		l.SetValidatedLicense(nil)
		return
	}

	// If a license is provided, validate the key and configure features based on the license state
	license, warn, err := validate.ValidateLicenseKey(context.TODO(), licenseString, model.Product_Gloo, model.AddOns{})
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

	if license.License == nil || license.License.AddOns.GraphQL == false {
		// If GraphQL is not explicitly enabled, disable the feature
		featureState.Reason = "License does not support GraphQL"
		return
	}

	featureState.Enabled = true
	featureState.Reason = ""
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
