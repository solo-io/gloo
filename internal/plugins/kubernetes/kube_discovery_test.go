package kubernetes

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path/filepath"

	"fmt"

	. "github.com/solo-io/glue/internal/endpointdiscovery/kube"
	"github.com/solo-io/glue/internal/pkg/kube/upstream"
	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	. "github.com/solo-io/glue/test/helpers"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("KubeSecretWatcher", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
		upstreams                 []gluev1.Upstream
	)
	BeforeSuite(func() {
		namespace = RandString(8)
		mkb = NewMinikube(false, namespace)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
		singlePortServiceSpec, err := ToMap(upstream.Spec{
			ServiceName:      "test-service",
			ServiceNamespace: namespace,
		})
		Must(err)
		serviceWithSpecificPortSpec, err := ToMap(upstream.Spec{
			ServiceName:      "test-service-with-port",
			ServiceNamespace: namespace,
			ServicePortName:  "portname",
		})
		Must(err)
		upstreams = []gluev1.Upstream{
			{
				Name: "my_upstream_with_specific_port",
				Type: upstream.Kubernetes,
				Spec: serviceWithSpecificPortSpec,
			},
			{
				Name: "my_upstream_with_one_port",
				Type: upstream.Kubernetes,
				Spec: singlePortServiceSpec,
			},
		}
	})
	AfterSuite(func() {
		mkb.Teardown()
	})
	Describe("controller", func() {
		It("watches kube endpoints", func() {
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			discovery, err := NewEndpointDiscovery(masterUrl, kubeconfigPath, time.Second, make(chan struct{}))
			Expect(err).NotTo(HaveOccurred())

			// add a pod and service pointing to it
			kubeClient, err := kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			go discovery.TrackUpstreams(upstreams)

			// create a service and a pod backing it for each upstream
			for _, us := range upstreams {
				serviceName := us.Spec["service_name"].(string)
				labels := map[string]string{"app": serviceName}
				pod := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-for-" + serviceName,
						Namespace: namespace,
						Labels:    labels,
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
				_, err := kubeClient.CoreV1().Pods(namespace).Create(pod)
				Expect(err).NotTo(HaveOccurred())
				service := &v1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      serviceName,
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
					},
				}
				portName := us.Spec["service_port_name"].(string)
				if portName != "" {
					service.Spec.Ports = append(service.Spec.Ports, v1.ServicePort{
						Name: portName,
						Port: 8081,
					})
				}
				_, err = kubeClient.CoreV1().Services(namespace).Create(service)
				Expect(err).NotTo(HaveOccurred())
				//created pods and services
			}

			Eventually(func() endpointdiscovery.EndpointGroups {
				select {
				case endpoints := <-discovery.Endpoints():
					return endpoints
				case err := <-discovery.Error():
					Expect(err).NotTo(HaveOccurred())
				}
				return nil
			}, time.Second*5).Should(HaveLen(len(upstreams)))

			for _, us := range upstreams {
				serviceName := us.Spec["service_name"].(string)
				select {
				case <-time.After(time.Second * 5):
					Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
				case endpoints := <-discovery.Endpoints():
					Expect(len(endpoints)).To(Equal(len(upstreams)))
					Expect(endpoints).To(HaveKey(us.Name))
					serviceEndpoints := endpoints[us.Name]
					Expect(serviceEndpoints).NotTo(BeNil())
					portName := us.Spec["service_port_name"].(string)
					Expect(serviceEndpoints).To(HaveLen(1))
					if portName != "" {
						Expect(serviceEndpoints[0].Port).To(Equal(int32(8081)))
					} else {
						Expect(serviceEndpoints[0].Port).To(Equal(int32(8080)))
					}
					podWithIp, err := kubeClient.CoreV1().Pods(namespace).Get("pod-for-"+serviceName, metav1.GetOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(podWithIp.Status.PodIP).To(Equal(serviceEndpoints[0].Address))
				case err := <-discovery.Error():
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})
	})
})
