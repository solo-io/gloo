package controller

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	kubeplugin "github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	. "github.com/solo-io/gloo/test/helpers"
)

var _ = Describe("KubeIngressController", func() {
	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
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
			ingressCtl *ServiceController
			kubeClient kubernetes.Interface
			glooClient storage.Interface
		)
		BeforeEach(func() {
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Must(err)

			glooClient, err = crd.NewStorage(cfg, namespace, time.Second)
			Must(err)

			ingressCtl, err = NewServiceController(cfg, glooClient, time.Second, true)
			Must(err)

			go ingressCtl.Run(make(chan struct{}))

			kubeClient, err = kubernetes.NewForConfig(cfg)
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
				time.Sleep(time.Second)
			})
			It("does not return an error", func() {
				select {
				case <-time.After(time.Second):
					// passed without error
				case err := <-ingressCtl.Error():
					Expect(err).NotTo(HaveOccurred())
					Fail("err passed, but was nil")
				}
			})
			It("ignores the ingress", func() {
				upstreams, err := glooClient.V1().Upstreams().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(upstreams).To(HaveLen(0))
				virtualHostList, err := glooClient.V1().VirtualHosts().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(virtualHostList).To(HaveLen(0))
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
						Annotations:  map[string]string{"kubernetes.io/ingress.class": GlooIngressClass},
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
				time.Sleep(time.Second)
			})
			It("does not return an error", func() {
				select {
				case <-time.After(time.Second * 1):
					// passed without error
				case err := <-ingressCtl.Error():
					Expect(err).NotTo(HaveOccurred())
					Fail("err passed, but was nil")
				}
			})
			It("creates the default virtualhost for the ingress", func() {
				glooRoute := &v1.Route{
					Matcher: &v1.Matcher{
						Path: &v1.Matcher_PathPrefix{
							PathPrefix: "/",
						},
					},
					SingleDestination: &v1.Destination{
						DestinationType: &v1.Destination_Upstream{
							Upstream: &v1.UpstreamDestination{
								Name: upstreamName(createdIngress.Namespace, *createdIngress.Spec.Backend),
							},
						},
					},
				}
				glooVirtualHost := &v1.VirtualHost{
					Name:    defaultVirtualHost,
					Domains: []string{"*"},
					Routes:  []*v1.Route{glooRoute},
				}
				virtualHostList, err := glooClient.V1().VirtualHosts().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(virtualHostList).To(HaveLen(1))
				virtualHost := virtualHostList[0]
				Expect(virtualHost.Name).To(Equal(defaultVirtualHost))
				// have to set metadata
				glooVirtualHost.Metadata = virtualHost.Metadata
				Expect(virtualHost).To(Equal(glooVirtualHost))
			})
		})
		Context("an ingress is created with multiple rules", func() {
			var (
				createdIngress *v1beta1.Ingress
				err            error
			)
			BeforeEach(func() {
				// add an ingress
				ingress := &v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "ingress-",
						Namespace:    namespace,
						Annotations:  map[string]string{"kubernetes.io/ingress.class": GlooIngressClass},
					},
					Spec: v1beta1.IngressSpec{
						TLS: []v1beta1.IngressTLS{
							{
								Hosts:      []string{"host1"},
								SecretName: "my-secret-1",
							},
							{
								Hosts:      []string{"host2"},
								SecretName: "my-secret-2",
							},
						},
						Rules: []v1beta1.IngressRule{
							{
								Host: "host1",
								IngressRuleValue: v1beta1.IngressRuleValue{
									HTTP: &v1beta1.HTTPIngressRuleValue{
										Paths: []v1beta1.HTTPIngressPath{
											{
												Path: "/foo/bar",
												Backend: v1beta1.IngressBackend{
													ServiceName: "service1",
													ServicePort: intstr.FromInt(1234),
												},
											},
											{
												Path: "/foo/baz",
												Backend: v1beta1.IngressBackend{
													ServiceName: "service2",
													ServicePort: intstr.FromInt(3456),
												},
											},
										},
									},
								},
							},
							{
								Host: "host2",
								IngressRuleValue: v1beta1.IngressRuleValue{
									HTTP: &v1beta1.HTTPIngressRuleValue{
										Paths: []v1beta1.HTTPIngressPath{
											{
												Path: "/foo/bar",
												Backend: v1beta1.IngressBackend{
													ServiceName: "service3",
													ServicePort: intstr.FromInt(1234),
												},
											},
											{
												Path: "/straw/berry",
												Backend: v1beta1.IngressBackend{
													ServiceName: "service4",
													ServicePort: intstr.FromString("foo"),
												},
											},
											{
												Path: "/bat/girl",
												Backend: v1beta1.IngressBackend{
													ServiceName: "service4",
													ServicePort: intstr.FromString("foo"),
												},
											},
										},
									},
								},
							},
						},
					},
				}
				createdIngress, err = kubeClient.ExtensionsV1beta1().Ingresses(namespace).Create(ingress)
				Must(err)
				time.Sleep(time.Second * 2)
			})
			AfterEach(func() {
				err = kubeClient.ExtensionsV1beta1().Ingresses(namespace).Delete(createdIngress.Name, nil)
				Must(err)
			})
			It("does not return an error", func() {
				select {
				case <-time.After(time.Second * 2):
					// passed without error
				case err := <-ingressCtl.Error():
					Expect(err).NotTo(HaveOccurred())
					Fail("err passed, but was nil")
				}
			})
			It("should de-duplicate repeated upstreams", func() {
				time.Sleep(time.Second * 3)
				expectedUpstreams := make(map[string]*v1.Upstream)
				for _, rule := range createdIngress.Spec.Rules {
					for _, path := range rule.HTTP.Paths {
						spec, _ := kubeplugin.EncodeUpstreamSpec(kubeplugin.UpstreamSpec{
							ServiceName:      path.Backend.ServiceName,
							ServiceNamespace: namespace,
							ServicePort:      path.Backend.ServicePort.String(),
						})
						expectedUpstreams[upstreamName(createdIngress.Namespace, path.Backend)] = &v1.Upstream{
							Name: upstreamName(createdIngress.Namespace, path.Backend),
							Type: kubeplugin.UpstreamTypeKube,
							Spec: spec,
						}
					}
				}
				upstreams, err := glooClient.V1().Upstreams().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(upstreams).To(HaveLen(len(expectedUpstreams)))
				for _, us := range upstreams {
					Expect(expectedUpstreams).To(HaveKey(us.Name))
					// ignore metadata, as teh client sets this
					us.Metadata = nil
					Expect(expectedUpstreams[us.Name]).To(Equal(us))
				}
			})
			It("create a route for every path", func() {
				time.Sleep(4 * time.Second)
				expectedVirtualHosts := map[string]*v1.VirtualHost{
					"host2": {
						Name:    "host2",
						Domains: []string{"host2"},
					},
					"host1": {
						Name:    "host1",
						Domains: []string{"host1"},
					},
				}

				for _, rule := range createdIngress.Spec.Rules {
					for _, path := range rule.HTTP.Paths {
						vHost := expectedVirtualHosts[rule.Host]
						vHost.Routes = append(vHost.Routes, &v1.Route{
							Matcher: &v1.Matcher{
								Path: &v1.Matcher_PathRegex{PathRegex: path.Path},
							},
							SingleDestination: &v1.Destination{
								DestinationType: &v1.Destination_Upstream{
									Upstream: &v1.UpstreamDestination{
										Name: upstreamName(createdIngress.Namespace, path.Backend),
									},
								},
							},
						})
						sortRoutes(vHost.Routes)
						expectedVirtualHosts[rule.Host] = vHost
					}
				}

				for _, tls := range createdIngress.Spec.TLS {
					for _, host := range tls.Hosts {
						vHost := expectedVirtualHosts[host]
						vHost.SslConfig = &v1.SSLConfig{
							SecretRef: tls.SecretName,
						}
						expectedVirtualHosts[host] = vHost
					}
				}
				virtuavirtualHostList, err := glooClient.V1().VirtualHosts().List()
				Expect(err).NotTo(HaveOccurred())
				Expect(virtuavirtualHostList).To(HaveLen(len(expectedVirtualHosts)))
				for _, virtualHost := range virtuavirtualHostList {
					// ignore metadata
					virtualHost.Metadata = nil
					Expect(expectedVirtualHosts).To(ContainElement(virtualHost))
				}
			})
		})
	})
})
