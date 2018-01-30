package discovery_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"

	. "github.com/solo-io/glue/implemented_modules/file/discovery"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/pkg/module"
	. "github.com/solo-io/glue/test/helpers"
)

var _ = Describe("FileSecretWatcher", func() {
	var (
		file      string
		err       error
		discovery module.Discovery
	)
	BeforeEach(func() {
		f, err := ioutil.TempFile("", "filesecrettest")
		Must(err)
		file = f.Name()
		discovery, err = NewServiceDiscovery(file, time.Millisecond)
		Must(err)
	})
	AfterEach(func() {
		log.Printf("removing " + file)
		os.RemoveAll(file)
	})
	Describe("watching file", func() {
		Context("an invalid structure is written to a file", func() {
			It("sends an error on the Error() channel", func() {
				invalidData := []byte("foo: bar")
				err = ioutil.WriteFile(file, invalidData, 0644)
				Expect(err).NotTo(HaveOccurred())
				select {
				case <-discovery.Endpoints():
					Fail("secretmap was received, expected error")
				case err := <-discovery.Error():
					Expect(err).To(HaveOccurred())
				case <-time.After(time.Second * 1):
					Fail("expected err to have occurred before 1s")
				}
			})
		})
		Context("no secrets wanted", func() {
			It("doesnt send anything on any channel", func() {
				missingSecrets := map[string]map[string][]byte{"another-key": {"foo": []byte("bar"), "baz": []byte("qux")}}
				data, err := json.Marshal(missingSecrets)
				Expect(err).NotTo(HaveOccurred())
				err = ioutil.WriteFile(file, data, 0644)
				Expect(err).NotTo(HaveOccurred())
				select {
				case <-discovery.Endpoints():
					Fail("secretmap was received, expected timeout")
				case err := <-discovery.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 1):
					// passed
				}
			})
		})
		Context("a valid config is written to a file", func() {
			It("sends a corresponding secretmap on Secrets()", func() {
				cfg := NewTestConfig()
				yml, err := yaml.Marshal(cfg)
				Must(err)
				err = ioutil.WriteFile(file, yml, 0644)
				Must(err)
				upstreams := make([]v1.Upstream, 1)
				for _, upstream := range cfg.Upstreams {
					upstreams[0] = upstream
					break
				}
				go discovery.TrackUpstreams(upstreams)
				select {
				case parsedSecrets := <-discovery.Endpoints():
					Expect(parsedSecrets).To(Equal(cfg))
				case err := <-discovery.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 5):
					Fail("expected new secrets to be read in before 1s")
				}
			})
		})
	})
})
