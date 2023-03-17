package license_test

import (
	"context"
	"os"
	"time"

	"github.com/solo-io/licensing/pkg/defaults"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/licensing/pkg/model"
	"github.com/solo-io/solo-projects/pkg/license"
)

var (
	ctx          context.Context
	cancel       context.CancelFunc
	licenseState = &model.License{
		IssuedAt:      time.Now(),
		ExpiresAt:     time.Now(),
		RandomPayload: "",
		LicenseType:   "",
		Product:       "",
		AddOns:        nil,
	}
	nilLicense     *license.ValidatedLicense
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
	km = defaults.KeyManager()
	kg = km.KeyGenerator()
)

var _ = Describe("LicensedFeatureProvider", func() {

	var (
		licensedFeatureProvider *license.LicensedFeatureProvider
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		licensedFeatureProvider = license.NewLicensedFeatureProvider()
	})

	AfterEach(func() {
		cancel()
	})

	DescribeTable("Enterprise",
		func(validatedLicense *license.ValidatedLicense, expectedEnabled bool) {
			licensedFeatureProvider.SetValidatedLicense(validatedLicense)

			featureState := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
			Expect(featureState.Enabled).To(Equal(expectedEnabled))
		},
		Entry("nil license", nilLicense, false),
		Entry("invalid license", invalidLicense, true),
		Entry("expired license", expiredLicense, true),
		Entry("valid license", validLicense, true),
	)

	DescribeTable("GraphQL",
		func(validatedLicense *license.ValidatedLicense, expectedEnabled bool) {
			licensedFeatureProvider.SetValidatedLicense(validatedLicense)

			featureState := licensedFeatureProvider.GetStateForLicensedFeature(license.GraphQL)
			Expect(featureState.Enabled).To(Equal(expectedEnabled))
		},
		Entry("nil license", nilLicense, false),
		Entry("invalid license", invalidLicense, true),
		Entry("valid license without add-on", validLicense, false),
		Entry("valid license with add-on", getGraphQLLicense(false, false), true),
		Entry("expired license with add-on", getGraphQLLicense(true, false), true),
		Entry("invalid license with add-on", getGraphQLLicense(true, true), true),
	)

	// getLicenseKey generates a license key string from the following:
	// The license type, either enterprise or trial (licenseType)
	// The product (product) and its days until expiration (days)
	// The products addons (addons) and each of their days until expiration (daysToAddonExpiration)
	getLicenseKey := func(licenseType model.LicenseType, product model.Product, days int, addons string, daysToAddonExpiration []int) string {
		Expect(model.AreValidAddons(addons)).ToNot(HaveOccurred(), "invalid addons")
		addOnLicensesWithoutExp, err := model.HandleLegacyAddons(addons)

		Expect(err).ToNot(HaveOccurred())
		Expect(len(daysToAddonExpiration)).To(Equal(len(addOnLicensesWithoutExp)))

		var addOnLicenses []model.AddOnLicense
		for i, lic := range addOnLicensesWithoutExp {
			lic.ExpiresAt = time.Now().Add(time.Duration(daysToAddonExpiration[i]*24) * time.Hour)
			lic.LicenseType = licenseType
			addOnLicenses = append(addOnLicenses, lic)
		}

		key, err := kg.GenerateKey(ctx, time.Now(), time.Now().Add(time.Duration(days*24)*time.Hour), licenseType, product, addOnLicenses)
		Expect(err).ToNot(HaveOccurred())
		return key
	}

	testLicense := func(key string, enterpriseLicense license.FeatureState, graphqlLicense license.FeatureState) {
		os.Setenv(license.EnvName, key)
		licensedFeatureProvider.ValidateAndSetLicense(ctx)
		featureState := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
		Expect(featureState).To(BeEquivalentTo(&enterpriseLicense))
		featureState = licensedFeatureProvider.GetStateForLicensedFeature(license.GraphQL)
		Expect(featureState).To(BeEquivalentTo(&graphqlLicense))
		os.Unsetenv(license.EnvName)
	}
	DescribeTable("New Licensing Client", func(key string, enterpriseLicense license.FeatureState, graphqlLicense license.FeatureState) {
		testLicense(key, enterpriseLicense, graphqlLicense)
	},
		Entry("no license", "",
			license.FeatureState{Enabled: false, Reason: "License not found, Enterprise features not included"},
			license.FeatureState{Enabled: false, Reason: "License not found, GraphQL features not included"}),
		Entry("invalid license", "some invalid license key",
			license.FeatureState{Enabled: true, Reason: license.ParseError},
			license.FeatureState{Enabled: true, Reason: license.ParseError}),
		Entry("license for wrong product", getLicenseKey(model.LicenseType_Enterprise, model.Product_GlooMesh, 1, "", nil),
			license.FeatureState{Enabled: true, Reason: license.MissingGloo},
			license.FeatureState{Enabled: true, Reason: license.MissingGloo}),
		Entry("expired license", getLicenseKey(model.LicenseType_Enterprise, model.Product_Gloo, -1, "", nil),
			license.FeatureState{Enabled: true, Reason: license.ExpiredLicense},
			license.FeatureState{Enabled: false, Reason: license.MissingGraphql}),
		Entry("expired trial license (treated as invalid)", getLicenseKey(model.LicenseType_Trial, model.Product_GlooTrial, -1, "", nil),
			license.FeatureState{Enabled: true, Reason: license.InvalidLicense},
			license.FeatureState{Enabled: true, Reason: license.InvalidLicense}),
		Entry("valid license without add-on", getLicenseKey(model.LicenseType_Enterprise, model.Product_Gloo, 1, "", nil),
			license.FeatureState{Enabled: true, Reason: ""},
			license.FeatureState{Enabled: false, Reason: license.MissingGraphql}),
		Entry("valid license with add-on", getLicenseKey(model.LicenseType_Enterprise, model.Product_Gloo, 1, model.GraphQL.String(), []int{1}),
			license.FeatureState{Enabled: true, Reason: ""},
			license.FeatureState{Enabled: true, Reason: ""}),
		Entry("valid trial license with add-on", getLicenseKey(model.LicenseType_Trial, model.Product_GlooTrial, 1, model.GraphQL.String(), []int{1}),
			license.FeatureState{Enabled: true, Reason: ""},
			license.FeatureState{Enabled: true, Reason: ""}),
		Entry("valid license with expired add-on", getLicenseKey(model.LicenseType_Enterprise, model.Product_Gloo, 1, model.GraphQL.String(), []int{-1}),
			license.FeatureState{Enabled: true, Reason: ""},
			license.FeatureState{Enabled: true, Reason: license.MissingGraphql}),
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
			AddOns: []model.AddOnLicense{
				{
					AddOn: model.GraphQL,
				},
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
