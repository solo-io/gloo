package mesh_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	prommodel "github.com/prometheus/common/model"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = PDescribe("Mesh", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		testClients services.TestClients

		pathToService string
		envoyFactory  *services.EnvoyFactory
		promFactory   *services.PrometheusFactory

		svcs []*services.Service
		mesh services.QuoteUnquoteMesh
	)

	BeforeSuite(func() {
		var err error
		pathToService, err = gexec.Build("github.com/solo-io/solo-projects/projects/glooshot/test/svc")
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		testClients = services.RunGateway(ctx, true)

		var err error
		envoyFactory, err = services.NewEnvoyFactory()
		Expect(err).NotTo(HaveOccurred())
		promFactory, err = services.NewPrometheusFactory()
		Expect(err).NotTo(HaveOccurred())

		svcs = []*services.Service{{
			Name:    "svc1",
			Process: pathToService,
		}, {
			Name:    "svc2",
			Process: pathToService,
		},
		}
		mesh = services.QuoteUnquoteMesh{}
		mesh.Start(envoyFactory, testClients, svcs)
	})

	AfterEach(func() {
		envoyFactory.Clean()
		promFactory.Clean()
		cancel()
		gexec.TerminateAndWait(2 * time.Second)
	})

	It("service mesh happy path", func() {
		// contact svc1 on the mesh port
		Eventually(func() (string, error) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/hello", svcs[0].MeshPort))
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			return string(body), nil

		}, "5s", "0.1s").Should(Equal("hellohello"))
	})

	It("service mesh fail path", func() {
		// test that adding fault works
		mesh.AddFault(1, 100.0)

		Eventually(func() error {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/hello", svcs[1].MeshPort))
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusServiceUnavailable {
				return fmt.Errorf("got %d; fault not engaged", resp.StatusCode)
			}
			return nil

		}, "5s", "0.1s").Should(Not(HaveOccurred()))

		Eventually(func() (string, error) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/hello", svcs[0].MeshPort))
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			return string(body), nil

		}, "5s", "0.1s").Should(ContainSubstring("fault filter abort"))

		// test that removing faults work
		mesh.RemoveFault(1)
		Eventually(func() error {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/hello", svcs[1].MeshPort))
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				return errors.New("fault is still engaged")
			}
			return nil
		}, "5s", "0.1s").Should(Not(HaveOccurred()))

	})

	It("service mesh fail with prom", func() {
		promInst := promFactory.NewPrometheusInstance()
		promInst.AddMesh(&mesh)

		// wait for prom to start
		Eventually(func() error {
			_, err := http.Get(fmt.Sprintf("http://localhost:%d/", promInst.Port))
			return err
		}, "5s", "0.1s").Should(Not(HaveOccurred()))

		// test that adding fault works
		mesh.AddFault(1, 100.0)
		// Make some faults
		Eventually(func() (string, error) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/hello", svcs[0].MeshPort))
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			return string(body), nil

		}, "5s", "0.1s").Should(ContainSubstring("fault filter abort"))

		client := promInst.Client()
		// see that promethues is up & scraped something
		Eventually(func() (string, error) {
			v, err := client.Query(ctx, "envoy_cluster_mesh_local_1_upstream_rq_xx", time.Time{})
			if err != nil {
				return "", err
			}
			return v.String(), nil
		}, "15s", "0.1s").Should(ContainSubstring("envoy_cluster_mesh_local_1_upstream_rq_xx"))

		value, err := client.Query(ctx, "envoy_cluster_mesh_local_1_upstream_rq_xx{envoy_response_code_class=\"5\"}", time.Time{})
		Expect(err).NotTo(HaveOccurred())

		Expect(value.Type()).To(Equal(prommodel.ValVector))
		vectorSample := value.(prommodel.Vector)

		Expect(vectorSample).To(Not(BeEmpty()))
		Expect(vectorSample[0].Value).To(BeNumerically(">", 0))

	})

})
