package kubernetes_test

import (
	"time"

	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/log"
	"k8s.io/apimachinery/pkg/labels"

	"os"
	"path/filepath"

	"fmt"

	. "github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/endpointdiscovery"
	. "github.com/solo-io/gloo/test/helpers"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var podsWithLabels = make(map[string]bool)

var _ = Describe("KubeSecretWatcher", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
		upstreams                 []*v1.Upstream
		specLabels                = map[string]string{"version": "v1"}
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
		serviceWithSpecificLabels, err := EncodeUpstreamSpec(UpstreamSpec{
			ServiceName:      "test-service-with-labels",
			ServiceNamespace: namespace,
			ServicePort:      "portname",
			Labels:           specLabels,
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
			{
				Name: "my_upstream_with_specific_port_v1",
				Type: UpstreamTypeKube,
				Spec: serviceWithSpecificLabels,
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

			// add a pod and service pointing to it
			kubeClient, err := kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			go discovery.TrackUpstreams(upstreams)

			// create a service and a pod backing it for each upstream
			for _, us := range upstreams {
				serviceName := us.Spec.Fields["service_name"].Kind.(*types.Value_StringValue).StringValue
				podLabels := map[string]string{"app": serviceName}
				withoutLabels := &kubev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-for-" + serviceName,
						Namespace: namespace,
						Labels:    podLabels,
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
				_, err := kubeClient.CoreV1().Pods(namespace).Create(withoutLabels)
				Expect(err).NotTo(HaveOccurred())
				withLabels := &kubev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-withlabels-for-" + serviceName,
						Namespace: namespace,
						Labels:    labels.Merge(podLabels, specLabels),
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
				_, err = kubeClient.CoreV1().Pods(namespace).Create(withLabels)
				Expect(err).NotTo(HaveOccurred())
				service := &kubev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      serviceName,
						Namespace: namespace,
					},
					Spec: kubev1.ServiceSpec{
						Selector: podLabels,
						Ports: []kubev1.ServicePort{
							{
								Name: "foo",
								Port: 8080,
							},
						},
					},
				}

				podsWithLabels[withLabels.Name] = true

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

			go discovery.Run(make(chan struct{}))
			time.Sleep(time.Second)

			Eventually(func() endpointdiscovery.EndpointGroups {
				select {
				case endpoints := <-discovery.Endpoints():
					return endpoints
				case err := <-discovery.Error():
					Expect(err).NotTo(HaveOccurred())
				}
				return nil
			}, time.Second*10).Should(HaveLen(len(upstreams)))

			for _, us := range upstreams {
				decodedSpec, err := DecodeUpstreamSpec(us.Spec)
				Expect(err).NotTo(HaveOccurred())
				serviceName := decodedSpec.ServiceName

				select {
				case <-time.After(time.Second * 5):
					Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
				case endpoints := <-discovery.Endpoints():
					Expect(len(endpoints)).To(Equal(len(upstreams)))
					Expect(endpoints).To(HaveKey(us.Name))
					serviceEndpoints := endpoints[us.Name]
					Expect(serviceEndpoints).NotTo(BeNil())
					careAboutLabels := len(decodedSpec.Labels) > 0
					log.Printf("%v got eps %v", decodedSpec.Labels, endpoints)
					if careAboutLabels {
						Expect(serviceEndpoints).To(HaveLen(1))
					} else {
						Expect(serviceEndpoints).To(HaveLen(2))
					}

					if decodedSpec.ServicePort != "" {
						Expect(serviceEndpoints[0].Port).To(Equal(int32(8081)))
					} else {
						Expect(serviceEndpoints[0].Port).To(Equal(int32(8080)))
					}

					if careAboutLabels {
						podWithIp, err := kubeClient.CoreV1().Pods(namespace).Get(podName(serviceName, true), metav1.GetOptions{})
						Expect(err).NotTo(HaveOccurred())
						log.Printf("pod %v", podWithIp)
						Expect(serviceEndpoints[0].Address).To(Equal(podWithIp.Status.PodIP))
					} else {
						podNoLabels, err := kubeClient.CoreV1().Pods(namespace).Get(podName(serviceName, false), metav1.GetOptions{})
						Expect(err).NotTo(HaveOccurred())
						podLabels, err := kubeClient.CoreV1().Pods(namespace).Get(podName(serviceName, true), metav1.GetOptions{})
						Expect(err).NotTo(HaveOccurred())
						var endpointFoundForPod bool
					loop:
						for _, ep := range serviceEndpoints {
							for _, pod := range []*kubev1.Pod{podNoLabels, podLabels} {
								if pod.Status.PodIP == ep.Address {
									endpointFoundForPod = true
									break loop
								}
							}
						}
						Expect(endpointFoundForPod).To(BeTrue())
					}

				case err := <-discovery.Error():
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})
	})
})

func podName(serviceName string, haslabels bool) string {
	if haslabels {
		return "pod-withlabels-for-" + serviceName
	}
	return "pod-for-" + serviceName
}
