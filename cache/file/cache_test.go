package file

import (
	"io/ioutil"
	"os"

	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/glue/pkg/api/types"
	. "github.com/solo-io/glue/test/helpers"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Cache", func() {
	Describe("FileCache", func() {
		var (
			dir string
			err error
		)
		BeforeEach(func() {
			dir, err = ioutil.TempDir("", "filecachetest")
			Must(err)
		})
		AfterEach(func() {
			os.RemoveAll(dir)
		})
		Describe("SetConfigChangeHandler(handler ConfigChangeHandler)", func() {
			for _, description := range []string{
				"a config file is created",
				"a config file is updated",
				"a config file is deleted",
			} {
				Context(description, func() {
					It("reads the file in as a new config", func() {
						cache := NewFileCache(dir, time.Millisecond)
						resultsChan := make(chan types.Config)
						cache.SetConfigChangeHandler(func(new types.Config) {
							resultsChan <- new
						})
						cfgA := newTestConfig()
						data, err := yaml.Marshal(cfgA)
						Expect(err).NotTo(HaveOccurred())
						err = ioutil.WriteFile(filepath.Join(dir, "config.yml"), data, 0644)
						Expect(err).NotTo(HaveOccurred())
						select {
						case result := <-resultsChan:
							Expect(result).To(Equal(cfgA))
						case <-time.After(time.Second):
							Fail("expected new config to be read in before 1s")
						}
					})
				})
			}
		})
		//Describe("SetConfig", func() {
		//	Context("passed a set of resources", func() {
		//		It("creates the resources as yaml files in the specified directory", func() {
		//			cache := NewFileCache(dir, time.Millisecond)
		//			cfg := newTestConfig()
		//			err = cache.SetConfig(cfg)
		//			Expect(err).NotTo(BeNil())
		//			for _, resource := range []string{routeDirName, upstreamDirName, virtualhostDirName}{
		//				_, err = os.Stat(filepath.Join(dir, ))
		//				Expect(err).NotTo(BeNil())
		//			}
		//			for _, route := range cfg.Routes {
		//				filename := route.Name()+".yml"
		//				raw :=
		//			}
		//		})
		//	})
		//})
	})
})

func NewTestConfig() types.Config {
	routes := []types.Route{
		{
			Matcher: types.Matcher{
				Path: types.PrefixPathMatcher{
					Prefix: "/foo",
				},
				Headers:     map[string]string{"x-foo-bar": ""},
				Verbs:       []string{"GET", "POST"},
				VirtualHost: "my_vhost",
			},
			Destination: types.FunctionDestination{
				FunctionName: "aws.foo",
			},
			Plugins: map[string]types.Spec{
				"auth": {"username": "alice", "password": "bob"},
			},
		},
		{
			Matcher: types.Matcher{
				Path: types.ExactPathMatcher{
					Path: "/bar",
				},
				Verbs: []string{"GET", "POST"},
			},
			Destination: types.UpstreamDestination{
				UpstreamName:  "my_upstream",
				RewritePrefix: "/baz",
			},
			Plugins: map[string]types.Spec{
				"auth": {"username": "alice", "password": "bob"},
			},
		},
	}
	upstreams := []types.Upstream{
		{
			Name: "aws",
			Type: "lambda",
			Spec: types.Spec{"region": "us_east_1", "secret_key_ref": "my-aws-secret-key", "access_key_ref": "my0aws-access-key"},
			Functions: []types.Function{
				{
					Name: "my_lambda_function",
					Spec: types.Spec{
						"context_parameter": map[string]string{"KEY": "VAL"},
					},
				},
			},
		},
		{
			Name: "my_upstream",
			Type: "service",
			Spec: types.Spec{"url": "https://myapi.example.com"},
		},
	}
	virtualhosts := []types.VirtualHost{
		{
			Domains:   []string{"*.example.io"},
			SSLConfig: types.SSLConfig{},
		},
	}
	return types.Config{
		Routes:       routes,
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}
