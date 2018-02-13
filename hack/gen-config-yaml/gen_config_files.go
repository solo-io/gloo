package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/solo-io/glue/internal/plugins/service"
	"github.com/solo-io/glue/pkg/api/types/v1"
)

func main() {
	cfg := NewTestConfig()
	outDir := "_glue_config"
	err := os.MkdirAll(filepath.Join(outDir, "upstreams"), 0755)
	must(err)
	err = os.MkdirAll(filepath.Join(outDir, "virtualhosts"), 0755)
	must(err)
	for _, upstream := range cfg.Upstreams {
		data, err := yaml.Marshal(upstream)
		must(err)
		filename := filepath.Join(outDir, "upstreams", fmt.Sprintf("upstream-%v.yml", upstream.Name))
		err = ioutil.WriteFile(filename, data, 0644)
		must(err)
	}
	for _, virtualHost := range cfg.VirtualHosts {
		data, err := yaml.Marshal(virtualHost)
		must(err)
		filename := filepath.Join(outDir, "virtualhosts", fmt.Sprintf("virtualhost-%v.yml", virtualHost.Name))
		err = ioutil.WriteFile(filename, data, 0644)
		must(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func NewTestConfig() v1.Config {
	upstreams := []v1.Upstream{
		{
			Name: "localhost-python",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{Addr: "localhost", Port: 8000},
				},
			}),
		},
	}
	virtualhosts := []v1.VirtualHost{
		NewTestVirtualHost("localhost-app", NewTestRoute()),
	}
	return v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func NewTestVirtualHost(name string, routes ...v1.Route) v1.VirtualHost {
	return v1.VirtualHost{
		Name:   name,
		Routes: routes,
	}
}

func NewTestRoute() v1.Route {
	return v1.Route{
		Matcher: v1.Matcher{
			Path: v1.Path{
				Regex: "/",
			},
			Verbs: []string{"GET", "POST"},
		},
		Destination: v1.Destination{
			SingleDestination: v1.SingleDestination{
				UpstreamDestination: &v1.UpstreamDestination{
					UpstreamName: "localhost-python",
				},
			},
		},
	}
}
