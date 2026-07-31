package discovery_watchlabels

import (
	"context"
	"net/http"
	"time"

	"github.com/onsi/gomega"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	glooMatchers "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	gloocore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/stretchr/testify/suite"
)

var _ e2e.NewSuiteFunc = NewDiscoveryWatchlabelsSuite

// discoveryWatchlabelsSuite is the Suite of tests for validating Upstream discovery behavior when watchLabels are enabled
// This suite replaces the "upstream discovery" Context block from kube2e gateway tests
type discoveryWatchlabelsSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation
}

func NewDiscoveryWatchlabelsSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &discoveryWatchlabelsSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *discoveryWatchlabelsSuite) TestDiscoverUpstreamMatchingWatchLabels() {
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, serviceWithLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete service")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, serviceWithoutLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete service")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, serviceWithNoMatchingLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete service")
	})

	// add one service with labels matching our watchLabels
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, serviceWithLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can apply service")

	// add one service without labels matching our watchLabels
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, serviceWithoutLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can apply service")

	// add one service with a label matching our watchLabels but with an unwatched value
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, serviceWithNoMatchingLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can apply service")

	// eventually an Upstream should be created for the Service with matching labels
	// Upstreams no longer report status if they have not been translated at all to avoid conflicting with
	// other syncers that have translated them, so we can only detect that the objects exist here
	labeledUsName := kubernetes.UpstreamName(s.testInstallation.Metadata.InstallNamespace, "example-svc", 8000)
	s.testInstallation.Assertions.EventuallyResourceExists(
		func() (resources.Resource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, labeledUsName, clients.ReadOpts{Ctx: s.ctx})
		},
	)

	// the Upstream should have DiscoveryMetadata labels matching the parent Service
	us, err := s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, labeledUsName, clients.ReadOpts{Ctx: s.ctx})
	s.Assert().NoError(err, "can read upstream")

	s.Assert().Equal(map[string]string{
		"watchedKey": "watchedValue",
		"bonusKey":   "bonusValue",
	}, us.GetDiscoveryMetadata().GetLabels())

	// no Upstream should be created for the Service that does not have the watchLabels
	noLabelsUsName := kubernetes.UpstreamName(s.testInstallation.Metadata.InstallNamespace, "example-svc-no-labels", 8000)
	s.testInstallation.Assertions.ConsistentlyObjectsNotExist(
		s.ctx, &v1.Upstream{
			ObjectMeta: metav1.ObjectMeta{
				Name:      noLabelsUsName,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			},
		},
	)

	// no Upstream should be created for the Service that has a watched label without a watched value
	noMatchingLabelsUsName := kubernetes.UpstreamName(s.testInstallation.Metadata.InstallNamespace, "example-svc-no-matching-labels", 8000)
	s.testInstallation.Assertions.ConsistentlyObjectsNotExist(
		s.ctx, &v1.Upstream{
			ObjectMeta: metav1.ObjectMeta{
				Name:      noMatchingLabelsUsName,
				Namespace: s.testInstallation.Metadata.InstallNamespace,
			},
		},
	)

	// modify the non-watched label on the labeled service
	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, serviceWithModifiedLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can re-apply service")

	// expect the Upstream's DiscoveryMeta to eventually match the modified labels from the parent Service
	s.testInstallation.Assertions.Gomega.Eventually(func() (map[string]string, error) {
		us, err = s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, labeledUsName, clients.ReadOpts{Ctx: s.ctx})
		return us.GetDiscoveryMetadata().GetLabels(), err
	}).Should(gomega.Equal(map[string]string{
		"watchedKey": "watchedValue",
		"bonusKey":   "bonusValue-modified",
	}))
}

func (s *discoveryWatchlabelsSuite) TestDiscoverySpecPreserved() {
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, serviceWithLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete service")
	})

	// add one service with labels matching our watchLabels
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, serviceWithLabelsManifest, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can apply service")

	// eventually an Upstream should be created for the Service with matching labels
	// Upstreams no longer report status if they have not been translated at all to avoid conflicting with
	// other syncers that have translated them, so we can only detect that the objects exist here
	labeledUsName := kubernetes.UpstreamName(s.testInstallation.Metadata.InstallNamespace, "example-svc", 8000)
	s.testInstallation.Assertions.EventuallyResourceExists(
		func() (resources.Resource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, labeledUsName, clients.ReadOpts{Ctx: s.ctx})
		},
	)

	// the Upstream should have DiscoveryMetadata labels matching the parent Service
	us, err := s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, labeledUsName, clients.ReadOpts{Ctx: s.ctx})
	s.Assert().NoError(err, "can read upstream")

	s.Assert().NotNil(us.GetKube())
	s.Assert().Nil(us.GetKube().GetServiceSpec())

	// modify the Upstream to have a ServiceSpec
	us.GetKube().ServiceSpec = &options.ServiceSpec{
		PluginType: &options.ServiceSpec_GrpcJsonTranscoder{},
	}
	updatedUs, err := s.testInstallation.ResourceClients.UpstreamClient().Write(us, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: true})
	s.Assert().NoError(err, "can update upstream")
	s.Assert().NotNil(updatedUs.GetKube().GetServiceSpec())

	// expect the Upstream to consistently have the modified Spec
	s.testInstallation.Assertions.Gomega.Consistently(func() (*options.ServiceSpec, error) {
		us, err := s.testInstallation.ResourceClients.UpstreamClient().Read(us.GetMetadata().GetNamespace(), us.GetMetadata().GetName(), clients.ReadOpts{Ctx: s.ctx})
		return us.GetKube().GetServiceSpec(), err
	}).Should(gomega.Not(gomega.BeNil()))
}

// TestDiscoveredSelectorTracksService verifies that changing a Service selector
// updates the discovered Upstream and its Envoy endpoints.
func (s *discoveryWatchlabelsSuite) TestDiscoveredSelectorTracksService() {
	installNs := s.testInstallation.Metadata.InstallNamespace
	usName := kubernetes.UpstreamName(installNs, "echo-svc", 8000)

	s.T().Cleanup(func() {
		err := s.testInstallation.ResourceClients.VirtualServiceClient().Delete(installNs, "echo-vs",
			clients.DeleteOpts{Ctx: s.ctx, IgnoreNotExist: true})
		s.Assertions.NoError(err, "can delete virtual service")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, serviceBlueGreenManifest, "-n", installNs)
		s.Assertions.NoError(err, "can delete blue/green service")

		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, testdefaults.CurlPodManifest)
		s.Assertions.NoError(err, "can delete curl pod")
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.CurlPodManifest)
	s.Assert().NoError(err, "can apply curl pod")

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, serviceBlueGreenManifest, "-n", installNs)
	s.Assert().NoError(err, "can apply blue/green service")

	s.testInstallation.Assertions.EventuallyResourceExists(
		func() (resources.Resource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(installNs, usName, clients.ReadOpts{Ctx: s.ctx})
		},
	)

	// Route directly to the discovered Upstream. A kube destination would use an
	// in-memory Upstream instead.
	_, err = s.testInstallation.ResourceClients.VirtualServiceClient().Write(&gatewayv1.VirtualService{
		Metadata: &gloocore.Metadata{
			Name:      "echo-vs",
			Namespace: installNs,
		},
		VirtualHost: &gatewayv1.VirtualHost{
			Domains: []string{"echo.example.com"},
			Routes: []*gatewayv1.Route{{
				Matchers: []*glooMatchers.Matcher{{
					PathSpecifier: &glooMatchers.Matcher_Prefix{Prefix: "/"},
				}},
				Action: &gatewayv1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: &gloocore.ResourceRef{
										Name:      usName,
										Namespace: installNs,
									},
								},
							},
						},
					},
				},
			}},
		},
	}, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: true})
	s.Assert().NoError(err, "can write virtual service")

	curlOpts := []curl.Option{
		curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{
			Name:      defaults.GatewayProxyName,
			Namespace: installNs,
		})),
		curl.WithHostHeader("echo.example.com"),
		curl.WithPort(80),
	}

	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		curlOpts,
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring("BLUE"),
		},
		time.Minute,
	)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, serviceBlueGreenFlippedManifest, "-n", installNs)
	s.Assert().NoError(err, "can flip service selector to green")

	s.testInstallation.Assertions.Gomega.Eventually(func() (map[string]string, error) {
		upstream, err := s.testInstallation.ResourceClients.UpstreamClient().Read(installNs, usName, clients.ReadOpts{Ctx: s.ctx})
		return upstream.GetKube().GetSelector(), err
	}, time.Minute, time.Second).Should(gomega.Equal(map[string]string{
		"app":     "echo",
		"version": "green",
	}))

	// Verify Envoy receives the EDS update after the selector changes.
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		curlOpts,
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring("GREEN"),
		},
		time.Minute,
	)
}
