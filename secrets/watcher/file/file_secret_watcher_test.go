package file_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/secrets/watcher"
	. "github.com/solo-io/glue/secrets/watcher/file"
	. "github.com/solo-io/glue/test/helpers"
)

var _ = Describe("FileSecretWatcher", func() {
	var (
		file  string
		err   error
		watch watcher.Watcher
	)
	BeforeEach(func() {
		f, err := ioutil.TempFile("", "filesecrettest")
		Must(err)
		file = f.Name()
		watch, err = NewFileSecretWatcher(file, time.Millisecond)
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
				case <-watch.Secrets():
					Fail("config was received, expected error")
				case err := <-watch.Error():
					Expect(err).To(HaveOccurred())
				case <-time.After(time.Second * 1):
					Fail("expected err to have occurred before 1s")
				}
			})
		})
		Context("no secrets wanted", func() {
			It("doesnt send anything on any channel", func() {
				missingSecrets := []byte("foo:\n  bar: baz\n  qux: qaz")
				err = ioutil.WriteFile(file, missingSecrets, 0644)
				Expect(err).NotTo(HaveOccurred())
				select {
				case <-watch.Secrets():
					Fail("config was received, expected timeout")
				case err := <-watch.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 1):
					// passed
				}
			})
		})
		Context("want secrets that the file doesn't contain", func() {
			It("sends an error on the Error() channel", func() {
				missingSecrets := []byte("foo:\n  bar: baz\n  qux: qaz")
				err = ioutil.WriteFile(file, missingSecrets, 0644)
				Expect(err).NotTo(HaveOccurred())
				go watch.UpdateRefs([]string{"this key really should not be in the secretmap"})
				select {
				case <-watch.Secrets():
					Fail("config was received, expected error")
				case err := <-watch.Error():
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("secretmap not found"))
				case <-time.After(time.Second * 1):
					Fail("expected err to have occurred before 1s")
				}
			})
		})
		Context("a valid config is written to a file", func() {
			It("sends a corresponding config on the Config()", func() {
				secrets := NewTestSecrets()
				yml, err := yaml.Marshal(secrets)
				Must(err)
				err = ioutil.WriteFile(file, yml, 0644)
				Must(err)
				var key string
				for k := range secrets {
					key = k
					break
				}
				go watch.UpdateRefs([]string{key})
				select {
				case parsedSecrets := <-watch.Secrets():
					Expect(parsedSecrets).To(Equal(secrets))
				case err := <-watch.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 5):
					Fail("expected new secrets to be read in before 1s")
				}
			})
		})
	})
})
