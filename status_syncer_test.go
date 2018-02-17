package ingress_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/gloo/test/helpers"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("StatusSyncer", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
		ingressService            = "test-ingress-svc"
	)
	BeforeEach(func() {
		namespace = RandString(8)
		mkb = NewMinikube(false, namespace)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
	})
	AfterEach(func() {
		mkb.Teardown()
	})
	Describe("controller", func() {
		var (
			ingressSync *ingressSyncer
			kubeClient  kubernetes.Interface
		)
		Context("an ingress is created with our ingress class", func() {
			var (
				createdIngress *v1beta1.Ingress
				err            error
			)
			BeforeEach(func() {
				cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
				Must(err)

				kubeClient, err = kubernetes.NewForConfig(cfg)
				Must(err)

				// create the "ingress pod" - really just a nginx nothingness
				// but we want the service to point at it
				labels := map[string]string{"app": ingressService}
				pod := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "pod-for-" + ingressService,
						Namespace:    namespace,
						Labels:       labels,
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "nginx",
								Image: "nginx:latest",
							},
						},
					},
				}
				_, err = kubeClient.CoreV1().Pods(namespace).Create(pod)
				Expect(err).NotTo(HaveOccurred())

				// create the ingress service, this guy is gonna be the source of our status
				service := &v1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      ingressService,
						Namespace: namespace,
					},
					Spec: v1.ServiceSpec{
						Selector: labels,
						Ports: []v1.ServicePort{
							{
								Name: "foo",
								Port: 8080,
							},
						},
						Type: v1.ServiceTypeLoadBalancer,
					},
				}
				_, err = kubeClient.CoreV1().Services(namespace).Create(service)
				Expect(err).NotTo(HaveOccurred())

				// add an ingress
				// content doesn't matter, we just want to update it with the status
				ingress := &v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "ingress-",
						Namespace:    namespace,
					},
					Spec: v1beta1.IngressSpec{
						Backend: &v1beta1.IngressBackend{
							ServiceName: "nonexistent-service",
							ServicePort: intstr.FromInt(8080),
						},
					},
				}
				createdIngress, err = kubeClient.ExtensionsV1beta1().Ingresses(namespace).Create(ingress)
				Must(err)

				time.Sleep(time.Second)

				ingressSync, err = NewIngressSyncer(cfg,
					time.Second, make(chan struct{}),
					true, ingressService)
				Must(err)

			})
			AfterEach(func() {
				err = kubeClient.ExtensionsV1beta1().Ingresses(namespace).Delete(createdIngress.Name, nil)
				Must(err)
				time.Sleep(time.Second)
			})
			It("does not return an error", func() {
				select {
				case <-time.After(time.Second):
					// passed without error
				case err := <-ingressSync.Error():
					Expect(err).NotTo(HaveOccurred())
					Fail("err passed, but was nil")
				}
			})
			It("copies the loadbalancer status from the service to the ingress", func() {
				ingress, err := kubeClient.ExtensionsV1beta1().Ingresses(namespace).Get(createdIngress.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				service, err := kubeClient.CoreV1().Services(namespace).Get(ingressService, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(ingress.Status.LoadBalancer.Ingress).To(Equal(service.Status.LoadBalancer.Ingress))
			})
		})
	})
})
