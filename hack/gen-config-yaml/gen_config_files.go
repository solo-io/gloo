package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/pkg/protoutil"

	"github.com/solo-io/glue/internal/plugins/service"
	"github.com/solo-io/glue/pkg/api/types/v1"
)

var upstreamAddr string

var upstreamHost string
var upstreamPort uint32

var upstreamName = "my-upstream"

func main() {
	flag.StringVar(&upstreamAddr, "addr", "localhost:8080", "upstream addr")
	flag.Parse()
	parts := strings.Split(upstreamAddr, ":")
	upstreamHost = parts[0]
	p, err := strconv.Atoi(parts[1])
	must(err)
	upstreamPort = uint32(p)
	cfg := NewTestConfig()
	outDir := "_glue_config"
	err = os.MkdirAll(filepath.Join(outDir, "upstreams"), 0755)
	must(err)
	err = os.MkdirAll(filepath.Join(outDir, "virtualhosts"), 0755)
	must(err)
	for _, upstream := range cfg.Upstreams {
		jsn, err := protoutil.Marshal(upstream)
		must(err)
		data, err := yaml.JSONToYAML(jsn)
		must(err)
		filename := filepath.Join(outDir, "upstreams", fmt.Sprintf("upstream-%v.yml", upstream.Name))
		err = ioutil.WriteFile(filename, data, 0644)
		must(err)
	}
	for _, virtualHost := range cfg.VirtualHosts {
		jsn, err := protoutil.Marshal(virtualHost)
		must(err)
		data, err := yaml.JSONToYAML(jsn)
		must(err)
		log.GreyPrintf("%s", jsn)
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
	upstreams := []*v1.Upstream{
		{
			Name: "localhost-python",
			Type: service.UpstreamTypeService,
			Spec: service.EncodeUpstreamSpec(service.UpstreamSpec{
				Hosts: []service.Host{
					{Addr: upstreamAddr, Port: upstreamPort},
				},
			}),
		},
	}
	virtualhosts := []*v1.VirtualHost{
		NewTestVirtualHost("localhost-app", NewTestRoute(), NewTestRouteMultiDest()),
	}
	return v1.Config{
		Upstreams:    upstreams,
		VirtualHosts: virtualhosts,
	}
}

func NewTestVirtualHost(name string, routes ...*v1.Route) *v1.VirtualHost {
	return &v1.VirtualHost{
		Name:   name,
		Routes: routes,
	}
}

func NewTestRoute() *v1.Route {
	return &v1.Route{
		Matcher: &v1.Matcher{
			Path: &v1.Matcher_PathPrefix{
				PathPrefix: "/foo",
			},
			Headers: map[string]string{"x-foo-bar": ""},
			Verbs:   []string{"GET", "POST"},
		},
		SingleDestination: &v1.Destination{
			DestinationType: &v1.Destination_Upstream{
				Upstream: &v1.UpstreamDestination{
					Name: upstreamName,
				},
			},
		},
	}
}

func NewTestRouteMultiDest() *v1.Route {
	return &v1.Route{
		Matcher: &v1.Matcher{
			Path: &v1.Matcher_PathPrefix{
				PathPrefix: "/foo",
			},
			Headers: map[string]string{"x-foo-bar": ""},
			Verbs:   []string{"GET", "POST"},
		},
		MultipleDestinations: []*v1.WeightedDestination{
			{
				Destination: &v1.Destination{
					DestinationType: &v1.Destination_Upstream{
						Upstream: &v1.UpstreamDestination{
							Name: upstreamName,
						},
					},
				},
				Weight: 5,
			},
			{
				Destination: &v1.Destination{
					DestinationType: &v1.Destination_Upstream{
						Upstream: &v1.UpstreamDestination{
							Name: upstreamName,
						},
					},
				},
				Weight: 10,
			},
		},
	}
}
