package file

import (
	"os"
	"path/filepath"
	"time"

	"github.com/solo-io/glue-storage"
)

type Client struct {
	v1 *v1client
}

const upstreamsDir = "upstreams"
const virtualHostsDir = "virtualhosts"

func NewStorage(dir string, syncFrequency time.Duration) (storage.Storage, error) {
	if dir == "" {
		dir = GlueDefaultDirectory
	}
	return &Client{
		v1: &v1client{
			upstreams: &upstreamsClient{
				dir:           filepath.Join(dir, upstreamsDir),
				syncFrequency: syncFrequency,
			},
			virtualHosts: &virtualHostsClient{
				dir:           filepath.Join(dir, virtualHostsDir),
				syncFrequency: syncFrequency,
			},
		},
	}, nil
}

func (c *Client) V1() storage.V1 {
	return c.v1
}

type v1client struct {
	dir          string
	upstreams    *upstreamsClient
	virtualHosts *virtualHostsClient
}

func (c *v1client) Register() error {
	err := os.MkdirAll(c.dir, 0755)
	if err == os.ErrExist {
		return nil
	}
	return err
}

func (c *v1client) Upstreams() storage.Upstreams {
	return c.upstreams
}

func (c *v1client) VirtualHosts() storage.VirtualHosts {
	return c.virtualHosts
}
