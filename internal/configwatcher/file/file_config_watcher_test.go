package configwatcher_test

import (
	"io/ioutil"
	"os"
	"time"

	"encoding/json"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/glue/implemented_modules/file/configwatcher"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/pkg/module"
	. "github.com/solo-io/glue/test/helpers"
)

var _ = Describe("FileConfigWatcher", func() {
	var (
		file  string
		err   error
		watch module.ConfigWatcher
	)
	BeforeEach(func() {
		f, err := ioutil.TempFile("", "filecachetest")
		Must(err)
		file = f.Name()
		watch, err = NewFileConfigWatcher(file, time.Millisecond)
		Must(err)
	})
	AfterEach(func() {
		log.Printf("removing " + file)
		os.RemoveAll(file)
	})
	Describe("watching file", func() {
		Context("an invalid config is written to a file", func() {
			It("sends an error on the Error() channel", func() {
				invalidConfig := []byte("wdf112 1`12")
				err = ioutil.WriteFile(file, invalidConfig, 0644)
				Expect(err).NotTo(HaveOccurred())
				select {
				case <-watch.Config():
					Fail("config was received, expected error")
				case err := <-watch.Error():
					Expect(err).To(HaveOccurred())
				case <-time.After(time.Second * 1):
					Fail("expected err to have occurred before 1s")
				}
			})
		})
		Context("a valid config is written to a file", func() {
			It("sends a corresponding config on the Config()", func() {
				cfg := NewTestConfig()
				yml, err := yaml.Marshal(cfg)
				Must(err)
				err = ioutil.WriteFile(file, yml, 0644)
				Must(err)
				var expectedCfg v1.Config
				data, err := json.Marshal(cfg)
				Expect(err).To(BeNil())
				err = json.Unmarshal(data, &expectedCfg)
				Expect(err).To(BeNil())
				select {
				case parsedCfg := <-watch.Config():
					Expect(*parsedCfg).To(Equal(expectedCfg))
				case err := <-watch.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 1):
					Fail("expected new config to be read in before 1s")
				}
			})
		})
	})
})
