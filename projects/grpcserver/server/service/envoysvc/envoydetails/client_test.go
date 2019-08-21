package envoydetails_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails/mocks"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	mockCtrl                 *gomock.Controller
	podsGetter               *mocks.MockPodsGetter
	podNamespacePodInterface *mocks.MockPodInterface
	httpGetter               *mocks.MockHttpGetter
	client                   envoydetails.Client
	podNamespace             = "test-pod-ns"
	testErr                  = errors.New("test-error")
)

var _ = Describe("Envoy Details Client Test", func() {

	getPod := func(name, ip, port, dump, id string) kubev1.Pod {
		annotations := make(map[string]string, 2)
		if port != "" {
			annotations[envoydetails.ReadConfigPortAnnotationKey] = port
		}
		if dump != "" {
			annotations[envoydetails.ReadConfigConfigDumpAnnotationKey] = dump
		}
		labels := make(map[string]string, 2)
		if id != "" {
			labels[envoydetails.GatewayProxyIdLabel] = id
		}
		labels["gloo"] = "gateway-proxy"

		return kubev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:   podNamespace,
				Name:        name,
				Annotations: annotations,
				Labels:      labels,
			},
			Status: kubev1.PodStatus{
				PodIP: ip,
			},
		}
	}

	getDetails := func(envoyName, fileContent, fileError string) *v1.EnvoyDetails {
		return &v1.EnvoyDetails{
			Name: envoyName,
			Raw: &v1.Raw{
				FileName:           envoyName + ".json",
				Content:            fileContent,
				ContentRenderError: fileError,
			},
		}
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		podsGetter = mocks.NewMockPodsGetter(mockCtrl)
		podNamespacePodInterface = mocks.NewMockPodInterface(mockCtrl)
		httpGetter = mocks.NewMockHttpGetter(mockCtrl)

		podsGetter.EXPECT().Pods(podNamespace).Return(podNamespacePodInterface)

		client = envoydetails.NewClient(podsGetter, httpGetter)
	})

	Describe("List", func() {
		Context("works", func() {
			It("uses gateway proxy ID for envoy name when present, else falls back to pod name", func() {
				podA := getPod("a", "ipa", "porta", "config-dump-a", "proxy-a")
				podB := getPod("b", "ipb", "portb", "config-dump-b", "")
				podList := &kubev1.PodList{Items: []kubev1.Pod{podA, podB}}

				dumpA := "a config dump"
				dumpB := "b config dump"

				podNamespacePodInterface.EXPECT().
					List(metav1.ListOptions{LabelSelector: envoydetails.GatewayProxyLabelSelector}).
					Return(podList, nil)
				httpGetter.EXPECT().Get("ipa", "porta", "config-dump-a").Return(dumpA, nil)
				httpGetter.EXPECT().Get("ipb", "portb", "config-dump-b").Return(dumpB, nil)

				actual, err := client.List(context.Background(), podNamespace)
				Expect(err).NotTo(HaveOccurred())
				expected := []*v1.EnvoyDetails{
					getDetails("proxy-a", dumpA, ""),
					getDetails("b", dumpB, ""),
				}
				Expect(actual).To(Equal(expected))
			})

			It("skips pods that don't have a port annotation", func() {
				podA := getPod("a", "ipa", "porta", "config-dump-a", "proxy-a")
				podB := getPod("b", "ipb", "", "config-dump-b", "proxy-b")
				podList := &kubev1.PodList{Items: []kubev1.Pod{podA, podB}}

				dumpA := "a config dump"

				podNamespacePodInterface.EXPECT().
					List(metav1.ListOptions{LabelSelector: envoydetails.GatewayProxyLabelSelector}).
					Return(podList, nil)
				httpGetter.EXPECT().Get("ipa", "porta", "config-dump-a").Return(dumpA, nil)

				actual, err := client.List(context.Background(), podNamespace)
				Expect(err).NotTo(HaveOccurred())
				expected := []*v1.EnvoyDetails{
					getDetails("proxy-a", dumpA, ""),
				}
				Expect(actual).To(Equal(expected))
			})

			It("skips pods that don't have a config dump path annotation", func() {
				podA := getPod("a", "ipa", "porta", "", "proxy-a")
				podB := getPod("b", "ipb", "portb", "config-dump-b", "proxy-b")
				podList := &kubev1.PodList{Items: []kubev1.Pod{podA, podB}}

				dumpB := "b config dump"

				podNamespacePodInterface.EXPECT().
					List(metav1.ListOptions{LabelSelector: envoydetails.GatewayProxyLabelSelector}).
					Return(podList, nil)
				httpGetter.EXPECT().Get("ipb", "portb", "config-dump-b").Return(dumpB, nil)

				actual, err := client.List(context.Background(), podNamespace)
				Expect(err).NotTo(HaveOccurred())
				expected := []*v1.EnvoyDetails{
					getDetails("proxy-b", dumpB, ""),
				}
				Expect(actual).To(Equal(expected))
			})

			It("displays a ContentRenderError on the raw object when the http getter errors", func() {
				pod := getPod("a", "ipa", "porta", "config-dump-a", "proxy-a")
				podList := &kubev1.PodList{Items: []kubev1.Pod{pod}}

				podNamespacePodInterface.EXPECT().
					List(metav1.ListOptions{LabelSelector: envoydetails.GatewayProxyLabelSelector}).
					Return(podList, nil)
				httpGetter.EXPECT().Get("ipa", "porta", "config-dump-a").Return("", testErr)

				actual, err := client.List(context.Background(), podNamespace)
				Expect(err).NotTo(HaveOccurred())
				expected := []*v1.EnvoyDetails{
					getDetails("proxy-a", "", envoydetails.FailedToGetEnvoyConfig(podNamespace, "a")),
				}
				Expect(actual).To(Equal(expected))
			})
		})

		Context("errors", func() {
			It("when the pod interface errors", func() {
				podNamespacePodInterface.EXPECT().
					List(metav1.ListOptions{LabelSelector: envoydetails.GatewayProxyLabelSelector}).
					Return(nil, testErr)

				_, err := client.List(context.Background(), podNamespace)
				Expect(err).To(HaveOccurred())
				expectedErr := envoydetails.FailedToListPodsError(testErr, podNamespace, envoydetails.GatewayProxyLabelSelector)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
	})
})
