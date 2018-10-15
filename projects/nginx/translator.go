package nginx

import (
	"github.com/ahmetb/go-linq"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gloo_solo_io1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

func virtualHostsFromVirtualServices(virtualServices []v1.VirtualService) []*gloo_solo_io1.VirtualHost {
	virtualServiceQuery := linq.From(virtualServices)
	virtualHosts := virtualServiceQuery.SelectT(func(virtualService v1.VirtualService) *gloo_solo_io1.VirtualHost {
		return virtualService.VirtualHost
	})

	var result []*gloo_solo_io1.VirtualHost
	virtualHosts.ToSlice(&result)
	return result
}

func locationsFromVirtualHosts(virtualHosts []*gloo_solo_io1.VirtualHost) []Location {
	prefixes := linq.From(prefixesFromVirtualHosts(virtualHosts))
	locations := prefixes.SelectT(locationFromPrefix)

	var result []Location
	locations.ToSlice(&result)
	return result
}

func prefixesFromVirtualHosts(virtualHosts []*gloo_solo_io1.VirtualHost) []string {
	virtualHostQuery := linq.From(virtualHosts)
	routes := virtualHostQuery.SelectManyT(func(virtualHost *gloo_solo_io1.VirtualHost) linq.Query {
		return linq.From(virtualHost.Routes)
	})
	matchers := routes.SelectT(func(route *gloo_solo_io1.Route) *gloo_solo_io1.Matcher {
		return route.Matcher
	})
	prefixes := matchers.SelectT(func(matcher *gloo_solo_io1.Matcher) string {
		return matcher.GetPrefix() // TODO(talnordan) What if the `Matcher` is not a prefix?
	})

	var result []string
	prefixes.ToSlice(&result)
	return result
}

func locationFromPrefix(prefix string) Location {
	return Location{
		Prefix: prefix,
	}
}

func Translate(gateway *v1.Gateway, virtualServices []v1.VirtualService) (Server, error) {
	virtualHosts := virtualHostsFromVirtualServices(virtualServices)
	locations := locationsFromVirtualHosts(virtualHosts)
	server := Server{
		Locations: locations,
	}
	return server, nil
}
