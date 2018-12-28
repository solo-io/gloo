package unit_test

import (
	"os"
	"time"

	"github.com/solo-io/solo-projects/projects/licensingserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/licensingserver/pkg/server"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type MockLicensingClient struct{}

func (mls *MockLicensingClient) Validate(key string) (bool, error) {
	if key == "5" {
		return true, nil
	} else {
		return false, nil
	}
}

var _ = Describe("Unit", func() {

	var (
		lvs *server.LicenseValidationServer
		lsc *licensingserver.LicensingServerClient
	)

	BeforeSuite(func() {
		var err error
		settings := server.Settings{
			HealthPort: 9001,
			GrpcPort:   9000,
		}
		lvs, err = server.NewServer(&MockLicensingClient{}, settings, nil)
		Expect(err).NotTo(HaveOccurred())
		go func() {
			err := lvs.Start()
			Expect(err).NotTo(HaveOccurred())
		}()
	})

	BeforeEach(func() {
		if _, present := os.LookupEnv("CLOUD_BUILD"); present {
			time.Sleep(time.Second)
		}
		var err error
		lsc, err = licensingserver.NewLicensingServerClient("", "9000", "9001")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := lsc.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("Health Check GRPC", func() {
		resp, err := lsc.HealthCheckGRPC()
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(healthpb.HealthCheckResponse_SERVING))
	})

	It("Health check HTTP", func() {
		resp, err := lsc.HealthCheckHTTP()
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(healthpb.HealthCheckResponse_SERVING))

	})

	It("produces an error when an empty key is passed", func() {
		_, err := lsc.Validate("")
		Expect(err).To(HaveOccurred())
	})

	It("Validate correct key", func() {
		valid, err := lsc.Validate("5")
		Expect(err).NotTo(HaveOccurred())
		Expect(valid).To(BeTrue())
	})

	It("validate incorrect key", func() {
		valid, err := lsc.Validate("6")
		Expect(err).NotTo(HaveOccurred())
		Expect(valid).To(BeFalse())
	})

	AfterSuite(func() {
		lvs.Close()
	})
})
