package kubernetes_test

import (
	"time"

	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path/filepath"

	"fmt"

	. "github.com/solo-io/glue/internal/plugins/kubernetes"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	. "github.com/solo-io/glue/test/helpers"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("KubeSecretWatcher", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
		upstreams                 []*v1.Upstream
	)
	BeforeSuite(func() {
		namespace = RandString(8)
		mkb = NewMinikube(false, namespace)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
		singlePortServiceSpec, err := EncodeUpstreamSpec(UpstreamSpec{
			ServiceName:      "test-service",
			ServiceNamespace: namespace,
		})
		Must(err)
		serviceWithSpecificPortSpec, err := EncodeUpstreamSpec(UpstreamSpec{
			ServiceName:      "test-service-with-port",
			ServiceNamespace: namespace,
			ServicePort:      "portname",
		})
		Must(err)
		upstreams = []*v1.Upstream{
			{
				Name: "my_upstream_with_specific_port",
				Type: UpstreamTypeKube,
				Spec: serviceWithSpecificPortSpec,
			},
			{
				Name: "my_upstream_with_one_port",
				Type: UpstreamTypeKube,
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

			discovery, err := NewEndpointDiscovery(masterUrl, kubeconfigPath, time.Second)
			Expect(err).NotTo(HaveOccurred())
			go discovery.Run(make(chan struct{}))

			// add a pod and service pointing to it
			kubeClient, err := kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			go discovery.TrackUpstreams(upstreams)

			// create a service and a pod backing it for each upstream
			for _, us := range upstreams {
				serviceName := us.Spec.Fields["service_name"].Kind.(*types.Value_StringValue).StringValue
				labels := map[string]string{"app": serviceName}
				pod := &kubev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-for-" + serviceName,
						Namespace: namespace,
						Labels:    labels,
					},
					Spec: kubev1.PodSpec{
						Containers: []kubev1.Container{
							{
								Name:  "nginx",
								Image: "nginx:latest",
							},
						},
					},
				}
				_, err := kubeClient.CoreV1().Pods(namespace).Create(pod)
				Expect(err).NotTo(HaveOccurred())
				service := &kubev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      serviceName,
						Namespace: namespace,
					},
					Spec: kubev1.ServiceSpec{
						Selector: labels,
						Ports: []kubev1.ServicePort{
							{
								Name: "foo",
								Port: 8080,
							},
						},
					},
				}
				portName := us.Spec.Fields["service_port"].Kind.(*types.Value_StringValue).StringValue
				if portName != "" {
					service.Spec.Ports = append(service.Spec.Ports, kubev1.ServicePort{
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
				serviceName := us.Spec.Fields["service_name"].Kind.(*types.Value_StringValue).StringValue
				select {
				case <-time.After(time.Second * 5):
					Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
				case endpoints := <-discovery.Endpoints():
					Expect(len(endpoints)).To(Equal(len(upstreams)))
					Expect(endpoints).To(HaveKey(us.Name))
					serviceEndpoints := endpoints[us.Name]
					Expect(serviceEndpoints).NotTo(BeNil())
					portName := us.Spec.Fields["service_port"].Kind.(*types.Value_StringValue).StringValue
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
