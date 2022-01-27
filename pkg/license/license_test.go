package license_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/licensing/pkg/model"
	"github.com/solo-io/solo-projects/pkg/license"
)

var (
	licenseState = &model.License{
		IssuedAt:      time.Now(),
		ExpiresAt:     time.Now(),
		RandomPayload: "",
		LicenseType:   "",
		Product:       "",
		AddOns:        model.AddOns{},
	}
	nilLicense = &license.ValidatedLicense{
		License: nil,
		Err:     nil,
		Warn:    nil,
	}
	invalidLicense = &license.ValidatedLicense{
		License: licenseState,
		Err:     eris.New("License is invalid"),
		Warn:    nil,
	}
	expiredLicense = &license.ValidatedLicense{
		License: licenseState,
		Err:     nil,
		Warn:    eris.New("License is expired"),
	}
	validLicense = &license.ValidatedLicense{
		License: licenseState,
		Err:     nil,
		Warn:    nil,
	}
)

var _ = Describe("LicensedFeatureProvider", func() {

	var (
		licensedFeatureProvider *license.LicensedFeatureProvider
	)

	BeforeEach(func() {
		licensedFeatureProvider = license.NewLicensedFeatureProvider()
	})

	table.DescribeTable("Enterprise",
		func(validatedLicense *license.ValidatedLicense, expectedEnabled bool) {
			licensedFeatureProvider.SetValidatedLicense(validatedLicense)

			featureState := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
			Expect(featureState.Enabled).To(Equal(expectedEnabled))
		},
		table.Entry("nil license", nilLicense, true),
		table.Entry("invalid license", invalidLicense, true),
		table.Entry("expired license", expiredLicense, true),
		table.Entry("valid license", validLicense, true),
	)

	table.DescribeTable("GraphQL",
		func(validatedLicense *license.ValidatedLicense, expectedEnabled bool) {
			licensedFeatureProvider.SetValidatedLicense(validatedLicense)

			featureState := licensedFeatureProvider.GetStateForLicensedFeature(license.GraphQL)
			Expect(featureState.Enabled).To(Equal(expectedEnabled))
		},
		table.Entry("nil license", nilLicense, false),
		table.Entry("invalid license", invalidLicense, true),
		table.Entry("valid license without add-on", validLicense, false),
		table.Entry("valid license with add-on", getGraphQLLicense(false, false), true),
		table.Entry("expired license with add-on", getGraphQLLicense(true, false), true),
		table.Entry("invalid license with add-on", getGraphQLLicense(true, true), true),
	)
})

func getGraphQLLicense(isExpired, isInvalid bool) *license.ValidatedLicense {
	l := &license.ValidatedLicense{
		License: &model.License{
			IssuedAt:      time.Now(),
			ExpiresAt:     time.Now(),
			RandomPayload: "",
			LicenseType:   "",
			Product:       "",
			AddOns: model.AddOns{
				GraphQL: true,
			},
		},
		Err:  nil,
		Warn: nil,
	}
	if isExpired {
		l.Warn = eris.New("License is expired")
	}
	if isInvalid {
		l.Err = eris.New("License is invalid")
	}
	return l
}
