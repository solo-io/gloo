package module

import (
	"encoding/json"

	"sync"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/glue/config"
	"github.com/solo-io/glue/pkg/log"
)

type mockModule struct{}

func (m *mockModule) Identifier() string {
	return "fake"
}

func (m *mockModule) SecretsToWatch(_ []byte) (SecretNames []string) {
	return nil
}

func (m *mockModule) Translate(_ map[string]string, configBlob []byte) (config.EnvoyResources, error) {
	log.Printf("called!")
	var mockRules []struct {
		VirtualHost string `json:"v_host"`
	}
	err := json.Unmarshal(configBlob, &mockRules)
	Expect(err).NotTo(HaveOccurred())
	var resources config.EnvoyResources
	for _, rule := range mockRules {
		resources.Routes = append(resources.Routes, config.RouteWrapper{VirtualHost: rule.VirtualHost})
	}
	return resources, nil
}

var _ = Describe("Registry", func() {
	BeforeEach(func() {
		ready = sync.WaitGroup{}
		ready.Add(1)
	})
	Describe("Global Lock", func() {
		It("waits for a config object to be loaded", func() {
			finished := make(chan bool)
			go func() {
				ready.Wait()
				finished <- true
			}()
			Init(config.NewConfig())
			Eventually(func() bool {
				return <-finished
			}).Should(Equal(true))
		})
	})
	Describe("Register", func() {
		var yml = `- fake:
    v_host: something
- fake:
    v_host: something_else`
		It("adds a module to the global registry", func() {
			Init(config.NewConfig())
			Register(&mockModule{})
			time.Sleep(time.Millisecond * 250)
			err := globalRegistry.cfg.Update([]byte(yml))
			Expect(err).NotTo(HaveOccurred())
			resources := globalRegistry.cfg.GetResources()
			Expect(len(resources)).To(Equal(1))
			Expect(resources[0].Routes[0].VirtualHost).To(ContainSubstring("something"))
		})
	})
})
