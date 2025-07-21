package convert

import "math/rand"

const (
	perConnectionBufferLimit = "kgateway.dev/per-connection-buffer-limit"
	routeWeight              = "kgateway.dev/route-weight"
	delegationLabelGroup     = "delegation.kgateway.dev"
	delegationLabel          = delegationLabelGroup + "/label"
	RandomSuffix             = 4
	RandomSeed               = 1
)

func (g *GatewayAPIOutput) Convert(opts *Options) error {

	// Convert upstreams to backends first so that we can reference them in the Settings and Routes
	for _, upstream := range g.edgeCache.Upstreams() {
		g.convertUpstreamToBackend(upstream)
	}

	for _, settings := range g.edgeCache.Settings() {
		// We only translate virtual services for ones that match a gateway selector
		// TODO in the future we could blindly convert VS and not attach it to anything
		err := g.convertSettings(settings)
		if err != nil {
			return err
		}
	}

	for _, gateway := range g.edgeCache.GlooGateways() {

		// we first need to generate Gateway objects with the correct names based on proxy Names
		// spec.proxyNames
		g.generateGatewaysFromProxyNames(gateway)

		// We only translate virtual services for ones that match a gateway selector
		// TODO(nick) - in the future we could blindly convert VS and not attach it to anything
		err := g.convertVirtualServices(gateway, opts.DisableListenerSets)
		if err != nil {
			return err
		}
	}

	for _, routeTable := range g.edgeCache.RouteTables() {
		err := g.convertRouteTableToHTTPRoute(routeTable)
		if err != nil {
			return err
		}
	}

	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var r = rand.New(rand.NewSource(RandomSeed))

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}
