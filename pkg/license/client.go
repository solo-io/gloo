package license

import "context"

type Client interface {
	IsLicenseValid() error
}

type client struct {
	ctx context.Context
}

func (c *client) IsLicenseValid() error {
	return LicenseStatus(c.ctx)
}

func NewClient(ctx context.Context) Client {
	return &client{ctx: ctx}
}
