package kube

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path/filepath"

	clientset "github.com/solo-io/glue/internal/configwatcher/kube/crd/client/clientset/versioned"
	"github.com/solo-io/glue/internal/pkg/kube/upstream"
	"github.com/solo-io/glue/pkg/api/types/v1"
	. "github.com/solo-io/glue/test/helpers"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("KubeSecretWatcher", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
	)
	BeforeSuite(func() {
		namespace = RandString(8)
		mkb = NewMinikube(false, namespace)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
	})
	AfterSuite(func() {
		mkb.Teardown()
	})
	Describe("controller", func() {
		var (
			ingressCvtr *ingressConverter
			kubeClient  kubernetes.Interface
			glueClient  clientset.Interface
		)
		BeforeEach(func() {
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Must(err)

			ingressCvtr, err = NewIngressConverter(cfg, time.Second, make(chan struct{}))
			Must(err)

			kubeClient, err = kubernetes.NewForConfig(cfg)
			Must(err)

			glueClient, err = clientset.NewForConfig(cfg)
			Must(err)

			err = RegisterCrds(cfg)
			Must(err)
		})
		Context("an ingress is created without our ingress class", func() {
			var (
				createdIngress *v1beta1.Ingress
				err            error
			)
			BeforeEach(func() {
				serviceName := "somethingsomethingsomething"
				servicePort := intstr.FromInt(8080)

				// add an ingress
				ingress := &v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "ingress-",
						Namespace:    namespace,
					},
					Spec: v1beta1.IngressSpec{
						Backend: &v1beta1.IngressBackend{
							ServiceName: serviceName,
							ServicePort: servicePort,
						},
					},
				}
				createdIngress, err = kubeClient.ExtensionsV1beta1().Ingresses(namespace).Create(ingress)
				Must(err)
			})
			AfterEach(func() {
				err = kubeClient.ExtensionsV1beta1().Ingresses(namespace).Delete(createdIngress.Name, nil)
				Must(err)
			})
			It("does not return an error", func() {
				select {
				case <-time.After(time.Second):
					// passed without error
				case err := <-ingressCvtr.Error():
					Expect(err).NotTo(HaveOccurred())
					Fail("err passed, but was nil")
				}
			})
			It("ignores the ingress", func() {
				upstreams, err := glueClient.GlueV1().Upstreams(namespace).List(metav1.ListOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(upstreams.Items).To(HaveLen(0))
				routes, err := glueClient.GlueV1().Routes(namespace).List(metav1.ListOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(routes.Items).To(HaveLen(0))
			})
		})
		Context("an ingress is created with a default backend", func() {
			var (
				createdIngress *v1beta1.Ingress
				err            error
			)
			BeforeEach(func() {
				serviceName := "somethingsomethingsomething"
				servicePort := intstr.FromInt(8080)

				// add an ingress
				ingress := &v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "ingress-",
						Namespace:    namespace,
						Annotations:  map[string]string{"kubernetes.io/ingress.class": GlueIngressClass},
					},
					Spec: v1beta1.IngressSpec{
						Backend: &v1beta1.IngressBackend{
							ServiceName: serviceName,
							ServicePort: servicePort,
						},
					},
				}
				createdIngress, err = kubeClient.ExtensionsV1beta1().Ingresses(namespace).Create(ingress)
				Must(err)
				time.Sleep(time.Second)
			})
			AfterEach(func() {
				err = kubeClient.ExtensionsV1beta1().Ingresses(namespace).Delete(createdIngress.Name, nil)
				Must(err)
			})
			It("does not return an error", func() {
				select {
				case <-time.After(time.Second * 2):
					// passed without error
				case err := <-ingressCvtr.Error():
					Expect(err).NotTo(HaveOccurred())
					Fail("err passed, but was nil")
				}
			})
			It("creates the expected upstream for the ingress", func() {
				upstreams, err := glueClient.GlueV1().Upstreams(namespace).List(metav1.ListOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(upstreams.Items).To(HaveLen(1))
				us := upstreams.Items[0]
				Expect(us.Name).To(Equal(upstreamName(createdIngress.Name, *createdIngress.Spec.Backend)))
				Expect(us.Spec.Type).To(Equal(upstream.Kubernetes))
			})
			It("creates the expected route for the ingress", func() {
				glueRoute := v1.Route{
					Matcher: v1.Matcher{
						Path: v1.Path{
							Prefix: "/",
						},
					},
					Destination: v1.Destination{
						UpstreamDestination: &v1.UpstreamDestination{
							UpstreamName: upstreamName(createdIngress.Name, *createdIngress.Spec.Backend),
						},
					},
					Weight: defaultRouteWeight,
				}
				routes, err := glueClient.GlueV1().Routes(namespace).List(metav1.ListOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(routes.Items).To(HaveLen(1))
				route := routes.Items[0]
				Expect(route.Name).To(Equal(routeName(glueRoute)))
				Expect(v1.Route(route.Spec)).To(Equal(glueRoute))
			})
		})
	})
})
