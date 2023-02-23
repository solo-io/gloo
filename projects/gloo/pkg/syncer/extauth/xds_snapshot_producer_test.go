package extauth_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	extauthsyncer "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth/test_fixtures"
)

const (
	validAuthConfigName         = "valid-ac"
	invalidAuthConfigName       = "invalid-ac"
	validCustomAuthServerName   = "custom-auth-server"
	invalidCustomAuthServerName = "invalid-custom-auth-server"
)

var _ = Describe("xDS Snapshot Producer", func() {

	Context("proxySourcedXdsSnapshotProducer", func() {
		// We intend to remove this instance of the XdsSnapshotProducer
		// All existing tests live in ../extauth_translator_syncer_test.go
	})

	Context("snapshotSourcedXdsSnapshotProducer", func() {

		var (
			ctx    context.Context
			cancel context.CancelFunc

			snapshotProducer extauthsyncer.XdsSnapshotProducer

			settings *v1.Settings
			snapshot *gloosnapshot.ApiSnapshot

			invalidAuthConfigError = extauthsyncer.NewInvalidAuthConfigError("passthrough grpc", &core.ResourceRef{
				Name:      invalidAuthConfigName,
				Namespace: writeNamespace,
			})
		)

		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())

			snapshotProducer = extauthsyncer.NewSnapshotSourcedXdsSnapshotProducer()

			snapshot = &gloosnapshot.ApiSnapshot{
				AuthConfigs: extauth.AuthConfigList{
					test_fixtures.BasicAuthConfig(validAuthConfigName, writeNamespace),
					test_fixtures.InvalidAuthConfig(invalidAuthConfigName, writeNamespace),
				},
			}
			settings = &v1.Settings{
				Extauth: &extauth.Settings{},
				NamedExtauth: map[string]*extauth.Settings{
					validCustomAuthServerName: {},
				},
			}
		})

		AfterEach(func() {
			cancel()
		})

		// We validate the various configuration that can be defined on just the Route level
		// This same configuration can be applied on the VirtualHost and WeightedDestination levels
		// but the logic is identical
		DescribeTable("ExtAuthExtension defined on Route",
			func(extAuthExtension *extauth.ExtAuthExtension, reportsAssertion types.GomegaMatcher) {
				By("Define a simple Proxy, with a single Route containing an ExtAuthExtension")
				snapshot.Proxies = v1.ProxyList{
					test_fixtures.ProxyWithExtAuthExtensionOnRoute("proxy", writeNamespace, extAuthExtension),
				}

				reports := make(reporter.ResourceReports)
				reports.Accept(snapshot.AuthConfigs.AsInputResources()...)
				reports.Accept(snapshot.Proxies.AsInputResources()...)

				By("Execute ProduceXdsSnapshot and perform assertions on results")
				xDSConfig := snapshotProducer.ProduceXdsSnapshot(ctx, settings, snapshot, reports)
				Expect(xDSConfig).To(HaveLen(1), "There is 1 valid AuthConfig defined in the API Snapshot")

				Expect(reports.ValidateStrict()).To(reportsAssertion)
			},
			Entry(
				"nil extension",
				nil,
				MatchError(ContainSubstring(invalidAuthConfigError.Error())),
			),
			Entry(
				"ref to valid AuthConfig",
				&extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_ConfigRef{
						ConfigRef: &core.ResourceRef{
							Name:      validAuthConfigName,
							Namespace: writeNamespace,
						},
					},
				},
				MatchError(ContainSubstring(invalidAuthConfigError.Error())),
			),
			Entry(
				"ref to invalid AuthConfig",
				&extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_ConfigRef{
						ConfigRef: &core.ResourceRef{
							Name:      invalidAuthConfigName,
							Namespace: writeNamespace,
						},
					},
				},
				And(
					MatchError(ContainSubstring(invalidAuthConfigError.Error())),
					// https://github.com/solo-io/gloo/issues/7272
					// If we add to the Proxy report, we will start validating VirtualServices (and other resources) which reference
					// Invalid AuthConfigs. I tried to introduce this, but a number of tests will need to be updated, so I backed
					// out the change.
					//MatchError(ContainSubstring("proxy references an authConfig gloo-system.invalid-ac which is invalid")),
				),
			),
			Entry(
				"ref to missing AuthConfig",
				&extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_ConfigRef{
						ConfigRef: &core.ResourceRef{
							Name:      "auth-config-not-in-snapshot",
							Namespace: writeNamespace,
						},
					},
				},
				And(
					MatchError(ContainSubstring(invalidAuthConfigError.Error())),
					MatchError(ContainSubstring("list did not find authConfig gloo-system.auth-config-not-in-snapshot")),
				),
			),
			Entry(
				"custom auth extension to valid auth server",
				&extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_CustomAuth{
						CustomAuth: &extauth.CustomAuth{
							Name: validCustomAuthServerName,
						},
					},
				},
				MatchError(ContainSubstring(invalidAuthConfigError.Error())),
			),
			Entry(
				"custom auth extension to invalid auth server",
				&extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_CustomAuth{
						CustomAuth: &extauth.CustomAuth{
							Name: invalidCustomAuthServerName,
						},
					},
				},
				And(
					MatchError(ContainSubstring(invalidAuthConfigError.Error())),
					MatchError(ContainSubstring("Unable to find custom auth server [invalid-custom-auth-server] in namedExtauth in Settings")),
				),
			),
		)
	})
})
