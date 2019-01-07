package keygen_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/licensingserver"
	"github.com/solo-io/solo-projects/projects/licensingserver/pkg/clients"
	"github.com/solo-io/solo-projects/projects/licensingserver/pkg/server"
)

const (
	RUN_KEYGEN_TESTS = "RUN_KEYGEN_TESTS"
)

var _ = Describe("e2e", func() {

	var (
		auth *clients.KeygenAuthConfig
		lvs  *server.LicenseValidationServer
		lsc  *licensingserver.LicensingServerClient
	)

	decodeAuth := func() error {
		var tempAuth clients.KeygenAuthConfig
		keygenCredsFile := os.ExpandEnv(clients.KEYGEN_CREDENTIALS_FILE)
		if keygenCredsFile == "" {
			err := envconfig.Process(clients.KEYGEN_ENV_EXPANSION, &tempAuth)
			if err != nil {
				return err
			}
		} else {
			byt, err := ioutil.ReadFile(keygenCredsFile)
			if err != nil {
				return err
			}
			err = json.Unmarshal(byt, &tempAuth)
			if err != nil {
				return err
			}
		}
		auth = &tempAuth
		return nil
	}

	BeforeSuite(func() {

		Expect(decodeAuth()).NotTo(HaveOccurred())
		var err error
		settings := server.Settings{
			GrpcPort:   9000,
			HealthPort: 9001,
		}
		client, err := clients.NewKeygenLicensingClient(auth)
		Expect(err).NotTo(HaveOccurred())
		lvs, err = server.NewServer(client, settings, nil)
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

	Context("License Validation", func() {
		It("can quick validate a valid license", func() {
			sampleLicense := "706715dd-92e2-4213-8923-997a9c2e31d1"
			valid, err := lsc.Validate(sampleLicense)
			Expect(err).NotTo(HaveOccurred())
			Expect(valid).To(BeTrue())
		})

		It("can quick validate an invalid license", func() {
			valid, err := lsc.Validate("123456")
			Expect(err).NotTo(HaveOccurred())
			Expect(valid).To(BeFalse())
		})

	})

	AfterEach(func() {
		err := lsc.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		if lvs != nil {
			lvs.Close()
		}
	})
})
