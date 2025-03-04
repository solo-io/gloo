package setup_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"
	corev1 "k8s.io/api/core/v1"
)

const (
	otherNS = "other-ns"
)

var _ = Describe("ControlPlane", func() {

	Context("xds service", func() {
		var (
			ctx       context.Context
			svcClient skkube.ServiceClient
			err       error

			svc1 *skkube.Service
			svc2 *skkube.Service
		)

		BeforeEach(func() {
			ctx = context.Background()
			inMemoryFactory := &factory.MemoryResourceClientFactory{
				Cache: memory.NewInMemoryResourceCache(),
			}
			svcClient, err = skkube.NewServiceClient(ctx, inMemoryFactory)
			Expect(err).NotTo(HaveOccurred())

			svc1 = skkube.NewService(defaults.GlooSystem, kubeutils.GlooServiceName)
			svc1.Labels = kubeutils.GlooServiceLabels
			svc1.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: kubeutils.GlooXdsPortName,
						Port: 1111,
					},
				},
			}

			svc2 = skkube.NewService(otherNS, kubeutils.GlooServiceName)
			svc2.Labels = kubeutils.GlooServiceLabels
			svc2.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: kubeutils.GlooXdsPortName,
						Port: 2222,
					},
				},
			}
		})

		AfterEach(func() {
			err := os.Unsetenv(statusutils.PodNamespaceEnvName)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns xds service from gloo service in default namespace", func() {
			// write both services
			_, err = svcClient.Write(svc1, clients.WriteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())
			_, err = svcClient.Write(svc2, clients.WriteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())

			// should return the port from gloo service in gloo-system
			service, err := setup.GetControlPlaneService(ctx, defaults.GlooSystem, svcClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(service).NotTo(BeNil())
			Expect(service.Name).To(Equal(svc1.Name))
			Expect(service.Namespace).To(Equal(svc1.Namespace))
		})

		It("returns error when no gloo service exists in the namespace", func() {
			// only write svc1
			_, err = svcClient.Write(svc1, clients.WriteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())

			_, err = setup.GetControlPlaneService(ctx, "another-ns", svcClient)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(setup.NoGlooSvcFoundError))
		})

		It("returns error when multiple gloo services exist in the namespace", func() {
			dupeSvc := skkube.NewService("other-ns", "duplicate-gloo-service")
			dupeSvc.Labels = kubeutils.GlooServiceLabels
			dupeSvc.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: kubeutils.GlooXdsPortName,
						Port: 3333,
					},
				},
			}

			// write both services
			_, err = svcClient.Write(svc2, clients.WriteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())
			_, err = svcClient.Write(dupeSvc, clients.WriteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())

			_, err = setup.GetControlPlaneService(ctx, otherNS, svcClient)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(setup.MultipleGlooSvcFoundError))
		})
	})

	Context("xds host", func() {
		It("get FQDN for service", func() {
			svc := skkube.NewService(defaults.GlooSystem, kubeutils.GlooServiceName)
			xdsHost := setup.GetControlPlaneXdsHost(svc)
			Expect(xdsHost).To(Equal("gloo.gloo-system.svc.cluster.local"))
		})
	})

	Context("xds port", func() {
		It("returns xds port", func() {
			svc := skkube.NewService(defaults.GlooSystem, kubeutils.GlooServiceName)
			svc.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: kubeutils.GlooXdsPortName,
						Port: 1111,
					},
				},
			}

			port, err := setup.GetControlPlaneXdsPort(svc)
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(int32(1111)))
		})

		It("returns error when the expected port name is not found", func() {
			svc := skkube.NewService(defaults.GlooSystem, kubeutils.GlooServiceName)
			svc.Spec = corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name: "not-the-right-name",
						Port: 1111,
					},
				},
			}

			_, err := setup.GetControlPlaneXdsPort(svc)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(setup.NoXdsPortFoundError))
		})
	})

})
