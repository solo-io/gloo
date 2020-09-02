package failover

import (
	"context"
	"net"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/consul"
)

func NewDnsResolver() consul.DnsResolver {
	return &failoverDnsResolverImpl{}
}

type failoverDnsResolverImpl struct{}

func (f *failoverDnsResolverImpl) Resolve(ctx context.Context, dnsName string) ([]net.IPAddr, error) {
	res := net.Resolver{
		PreferGo: true, // otherwise we may use cgo which doesn't resolve on my mac in testing
		Dial: func(ctx context.Context, network, address string) (conn net.Conn, err error) {
			// DNS typically uses UDP and falls back to TCP if the response size is greater than one packet
			// (originally 512 bytes). we use TCP to ensure we receive all IPs in a large DNS response
			return net.Dial("tcp", address)
		},
	}
	return res.LookupIPAddr(ctx, dnsName)
}
