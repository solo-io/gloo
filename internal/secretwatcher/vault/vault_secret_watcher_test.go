package vault_test

import (
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	. "github.com/solo-io/gloo-testing/helpers"
	. "github.com/solo-io/gloo/internal/secretwatcher/vault"
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

var _ = Describe("watching file", func() {
	var (
		containerName, token string
		err                  error
		watch                secretwatcher.Interface
		vaultAddr            = "http://127.0.0.1:8200"
		client               *vaultapi.Client
		secretPath           = "/secret/something"
	)
	BeforeEach(func() {
		containerName = "vault-test-" + RandString(4)
		token = "asdf" //RandString(8)
		err := DockerRunVault(containerName, token)
		Must(err)
		time.Sleep(time.Second)
		cfg := vaultapi.DefaultConfig()
		cfg.Address = vaultAddr
		client, err = vaultapi.NewClient(cfg)
		Must(err)
		client.SetToken(token)
		watch, err = NewVaultSecretWatcher(time.Second, 1, vaultAddr, token, make(chan struct{}))
		Must(err)
	})
	AfterEach(func() {
		err := DockerRm(containerName)
		Must(err)
	})
	Context("no secrets wanted", func() {
		It("doesnt send anything on any channel", func() {
			client.Logical().Write(secretPath, map[string]interface{}{"some": "secret"})
			Expect(err).NotTo(HaveOccurred())
			select {
			case <-watch.Secrets():
				Fail("secretmap was received, expected timeout")
			case err := <-watch.Error():
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 1):
				// passed
			}
		})
	})
	Context("want secrets that the file doesn't contain", func() {
		It("sends nothing", func() {
			client.Logical().Write(secretPath, map[string]interface{}{"some": "secret"})
			Expect(err).NotTo(HaveOccurred())
			go watch.TrackSecrets([]string{"this key really should not be in the secretmap"})
			select {
			case <-watch.Secrets():
				Fail("secretmap was received, expected error")
			case err := <-watch.Error():
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 1):
				// passed
			}
		})
	})
	Context("a valid config is written to a file", func() {
		It("sends a corresponding secretmap on Secrets()", func() {
			secrets := map[string]interface{}{
				"password": "foo",
			}
			stringSecrets := map[string]string{
				"password": "foo",
			}
			_, err := client.Logical().Write(secretPath, secrets)
			Must(err)
			go watch.TrackSecrets([]string{secretPath})
			select {
			case parsedSecrets := <-watch.Secrets():
				Expect(parsedSecrets).To(Equal(secretwatcher.SecretMap{secretPath: stringSecrets}))
			case err := <-watch.Error():
				Expect(err).NotTo(HaveOccurred())
			case <-time.After(time.Second * 5):
				Fail("expected new secrets to be read in before 1s")
			}
		})
	})
})
