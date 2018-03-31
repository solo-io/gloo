package faas_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo-plugins/kubernetes"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-function-discovery/internal/updater/faas"
	"github.com/solo-io/gloo-plugins/rest"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
)

func dummyListFuncs(gw string) (FaasFunctions, error) {
	funcs := FaasFunctions{
		{
			Name: "test",
		},
	}

	return funcs, nil
}

func getKubeUs(ns, n string) *v1.Upstream {
	if ns == "" {
		ns = "openfaas"
	}
	if n == "" {
		n = "gateway"
	}
	return &v1.Upstream{
		Name:     "openfaas-gateway-8080",
		Type:     kubernetes.UpstreamTypeKube,
		Metadata: &v1.Metadata{Namespace: "gloo-system"},
		Spec: kubernetes.EncodeUpstreamSpec(kubernetes.UpstreamSpec{
			ServiceName:      n,
			ServiceNamespace: ns,
			ServicePort:      8080,
		}),
	}
}

func getServiceUs(ns, n string) *v1.Upstream {
	if ns == "" {
		ns = "openfaas"
	}
	if n == "" {
		n = "gateway"
	}
	return &v1.Upstream{
		Name:     n,
		Type:     service.UpstreamTypeService,
		Metadata: &v1.Metadata{Namespace: ns},
		Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
			Hosts: []service.Host{{
				Addr: "gateway",
				Port: 80,
			}},
		}),
	}
}

func getServices(ns, n string) <-chan *v1.Upstream {
	c := make(chan *v1.Upstream)
	go func() {
		c <- getKubeUs(ns, n)
		c <- getServiceUs(ns, n)
		close(c)
	}()
	return c
}

var _ = Describe("Faas", func() {
	It("should get list of functions", func() {
		fr := FassRetriever{Lister: dummyListFuncs}

		for us := range getServices("", "") {

			funcs, err := fr.GetFuncs(us)

			Expect(err).NotTo(HaveOccurred())
			Expect(funcs).To(Not(BeEmpty()))
			Expect(funcs).To(HaveLen(1))

			outfunc := funcs[0]

			outspec, err := rest.DecodeFunctionSpec(outfunc.Spec)
			Expect(err).NotTo(HaveOccurred())

			Expect(outspec.Path).To(Equal("/function/test"))
		}

	})

	It("should ignore non gateway faas upstreams", func() {
		fr := FassRetriever{Lister: nil}

		for us := range getServices("", "not-the-gateway") {

			funcs, err := fr.GetFuncs(us)

			Expect(err).NotTo(HaveOccurred())
			Expect(funcs).To(BeEmpty())
		}
	})

	It("should ignore non faas upstreams", func() {
		fr := FassRetriever{Lister: nil}

		for us := range getServices("not-openfaas", "") {

			funcs, err := fr.GetFuncs(us)

			Expect(err).NotTo(HaveOccurred())
			Expect(funcs).To(BeEmpty())
		}
	})

	It("get correct gw kube", func() {
		var gw string
		f := func(gwinner string) (FaasFunctions, error) {
			gw = gwinner
			return nil, nil
		}
		fr := FassRetriever{Lister: f}
		us := getKubeUs("", "")

		fr.GetFuncs(us)
		Expect(gw).To(Equal("http://gateway.openfaas.svc.cluster.local:8080/"))
	})

	It("get correct gw service", func() {
		var gw string
		f := func(gwinner string) (FaasFunctions, error) {
			gw = gwinner
			return nil, nil
		}
		fr := FassRetriever{Lister: f}
		us := getServiceUs("", "")

		fr.GetFuncs(us)
		Expect(gw).To(Equal("http://gateway:80/"))
	})
})
