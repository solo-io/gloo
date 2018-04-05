package file

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"

	"path/filepath"

	. "github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/secretwatcher"
	. "github.com/solo-io/gloo/pkg/secretwatcher/file"
)

var _ = Describe("FileSecretWatcher", func() {
	var (
		dir   string
		file  string
		ref   string
		err   error
		watch secretwatcher.Interface
	)
	BeforeEach(func() {
		ref = "secrets.yml"
		dir, err = ioutil.TempDir("", "filesecrettest")
		Must(err)
		file = filepath.Join(dir, ref)
		watch, err = NewSecretWatcher(dir, time.Millisecond)
		Must(err)
	})
	AfterEach(func() {
		log.Debugf("removing " + dir)
		os.RemoveAll(dir)
	})
	Describe("watching file", func() {
		Context("no secrets wanted", func() {
			It("doesnt send anything on any channel", func() {
				missingSecrets := map[string]map[string][]byte{"another-key": {"foo": []byte("bar"), "baz": []byte("qux")}}
				data, err := json.Marshal(missingSecrets)
				Expect(err).NotTo(HaveOccurred())
				err = ioutil.WriteFile(file, data, 0644)
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
		Context("valid secrets are written to the ref file", func() {
			It("sends a corresponding secretmap on Secrets()", func() {

				secrets := map[string]string{"username": "me@example.com", "password": "foobar"}
				yml, err := yaml.Marshal(secrets)
				Must(err)
				err = ioutil.WriteFile(file, yml, 0644)
				Must(err)
				go watch.TrackSecrets([]string{ref})
				select {
				case parsedSecrets := <-watch.Secrets():
					Expect(parsedSecrets).To(Equal(secretwatcher.SecretMap{
						ref: secrets,
					}))
				case err := <-watch.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 5):
					Fail("expected new secrets to be read in before 1s")
				}
			})
		})
	})
})
