package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type Client struct {
	v1 *v1client
}

// TODO: support basic auth and tls
func NewStorage(cfg *api.Config, rootPath string, syncFrequency time.Duration) (*Client, error) {
	// Get a new client
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating consul client")
	}

	return &Client{
		v1: &v1client{
			upstreams: &upstreamsClient{
				consul:        client,
				rootPath:      rootPath + "/upstreams",
				syncFrequency: syncFrequency,
			},
			//virtualHosts: &virtualHostsClient{
			//	dir:           filepath.Join(dir, virtualHostsDir),
			//	syncFrequency: syncFrequency,
			//},
		},
	}, nil
}

func (c *Client) V1() *v1client {
	return c.v1
}

type v1client struct {
	upstreams *upstreamsClient
	//virtualHosts *virtualHostsClient
}

func (c *v1client) Register() error {
	return nil
}

func (c *v1client) Upstreams() *upstreamsClient {
	return c.upstreams
}

//func (c *v1client) VirtualHosts() storage.VirtualHosts {
//	return c.virtualHosts
//}
