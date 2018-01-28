package file_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/glue/config/watcher"
	. "github.com/solo-io/glue/config/watcher/file"
	"github.com/solo-io/glue/pkg/log"
	. "github.com/solo-io/glue/test/helpers"
)

var _ = Describe("Watcher", func() {
	var (
		dir   string
		err   error
		watch watcher.Watcher
	)
	BeforeEach(func() {
		dir, err = ioutil.TempDir("", "filecachetest")
		Must(err)
		watch, err = NewFileWatcher(dir, time.Millisecond)
		Must(err)
	})
	AfterEach(func() {
		log.Printf("removing " + dir)
		os.RemoveAll(dir)
	})
	Describe("watching directory", func() {
		Context("an invalid config is written to a file", func() {
			It("sends an error on the Error() channel", func() {
				invalidConfig := []byte("in: valid")
				err = ioutil.WriteFile(filepath.Join(dir, "config.yml"), invalidConfig, 0644)
				Expect(err).NotTo(HaveOccurred())
				select {
				case <-watch.Config():
					Fail("config was received, expected error")
				case err := <-watch.Error():
					Expect(err).To(HaveOccurred())
				case <-time.After(time.Second):
					Fail("expected new config to be read in before 1s")
				}
			})
		})
		Context("a valid config is written to a file", func() {
			It("sends a corresponding config on the Config()", func() {
				cfg := NewTestConfig()
				m := jsonpb.Marshaler{}
				str, err := m.MarshalToString(cfg)
				Must(err)
				jsn := []byte(str)
				yml, err := yaml.JSONToYAML(jsn)
				Must(err)
				err = ioutil.WriteFile(filepath.Join(dir, "config.yml"), yml, 0644)
				Must(err)
				select {
				case parsedCfg := <-watch.Config():
					Expect(parsedCfg).To(Equal(cfg))
				case err := <-watch.Error():
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 30):
					Fail("expected new config to be read in before 1s")
				}
			})
		})
	})
})
