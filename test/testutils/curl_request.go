package testutils

import (
	"fmt"
	"net/http"
)

// CurlRequestBuilder simplifies the process of generating curl requests in tests
type CurlRequestBuilder struct {
	verbose           bool
	withoutStats      bool
	connectionTimeout int // seconds
	returnHeaders     bool

	method string

	protocol string
}

// DefaultCurlRequestBuilder returns a CurlRequestBuilder with some default values
func DefaultCurlRequestBuilder() *CurlRequestBuilder {
	return &CurlRequestBuilder{
		verbose:           false,
		withoutStats:      false,
		connectionTimeout: 3,
		returnHeaders:     false,
		method:            http.MethodGet,
	}
}

func (c *CurlRequestBuilder) VerboseOutput() *CurlRequestBuilder {
	c.verbose = true
	return c
}

func (c *CurlRequestBuilder) WithoutStatus() *CurlRequestBuilder {
	c.withoutStats = true
	return c
}

func (c *CurlRequestBuilder) WithConnectionTimeout(seconds int) *CurlRequestBuilder {
	c.withoutStats = true
	return c
}

func (c *CurlRequestBuilder) BuildArgs() []string {
	args := []string{"curl"}

	if c.verbose {
		args = append(args, "-v")
	}
	if c.withoutStats {
		args = append(args, "-s")
	}
	if c.connectionTimeout > 0 {
		seconds := fmt.Sprintf("%v", c.connectionTimeout)
		args = append(args, "--connect-timeout", seconds, "--max-time", seconds)
	}
	if c.returnHeaders {
		args = append(args, "-I")
	}

	if c.method != http.MethodGet && c.method != "" {
		args = append(args, "-X"+c.method)
	}

	if opts.Host != "" {
		args = append(args, "-H", "Host: "+opts.Host)
	}
	if opts.CaFile != "" {
		args = append(args, "--cacert", opts.CaFile)
	}
	if opts.Body != "" {
		args = append(args, "-H", "Content-Type: application/json")
		args = append(args, "-d", opts.Body)
	}
	for h, v := range opts.Headers {
		args = append(args, "-H", fmt.Sprintf("%v: %v", h, v))
	}
	if opts.AllowInsecure {
		args = append(args, "-k")
	}

	port := opts.Port
	if port == 0 {
		port = 8080
	}
	protocol := opts.Protocol
	if protocol == "" {
		protocol = "http"
	}
	service := opts.Service
	if service == "" {
		service = "test-ingress"
	}
	if opts.SelfSigned {
		args = append(args, "-k")
	}
	if opts.Sni != "" {
		sniResolution := fmt.Sprintf("%s:%d:%s", opts.Sni, port, service)
		fullAddress := fmt.Sprintf("%s://%s:%d", protocol, opts.Sni, port)
		args = append(args, "--resolve", sniResolution, fullAddress)
	} else {
		args = append(args, fmt.Sprintf("%v://%s:%v%s", protocol, service, port, opts.Path))
	}

	return nil
}
