package envoy

import (
	"fmt"
	v2 "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/envoy/gloo-mesh-client-go/networking.gloo.solo.io/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"log"
	"os"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	"strings"
)

var runtimeScheme *runtime.Scheme
var codecs serializer.CodecFactory
var decoder runtime.Decoder

var SchemeBuilder = runtime.SchemeBuilder{
	// Solo Gloo Mesh API resources
	v2.AddToScheme,
}

func init() {
	runtimeScheme = runtime.NewScheme()
	if err := SchemeBuilder.AddToScheme(runtimeScheme); err != nil {
		log.Fatal(err)
	}
	codecs = serializer.NewCodecFactory(runtimeScheme)
	decoder = codecs.UniversalDeserializer()
}

// TODO try and combine httproutes by some hostname matching
// TODO try and join route options with ones that have the same RouteOption settings

func (o *Outputs) PostProcess(routeTablesFile string) error {
	log.Printf("Current Routes %d\n", len(o.httpRoutes))
	if routeTablesFile != "" {
		err := o.CombineByGlooMeshRouteTableHosts(routeTablesFile)
		if err != nil {
			return err
		}

	}

	////Try and combine httproutes by domain
	//if err := o.combineHTTPRoutes(4); err != nil {
	//	return err
	//}
	return nil
}

// This doesnt do anything as the envoy host matchers already are the same, this is just reverence code  if needed later.
// please delete me in the future
func (o *Outputs) CombineByGlooMeshRouteTableHosts(routeTablesFile string) error {
	// load all the route tables
	data, err := os.ReadFile(routeTablesFile)
	if err != nil {
		return err
	}
	var routeTables []*v2.RouteTable
	for _, resourceYAML := range strings.Split(string(data), "---") {
		// yaml to object
		obj, _, err := decoder.Decode([]byte(resourceYAML), nil, nil)
		if err != nil {
			if runtime.IsNotRegisteredError(err) {
				// we just want to add the yaml and move on
				fmt.Print("object not registered", resourceYAML)
				continue
			}

			// TODO if we cant decode it, don't do anything and continue
			// log.Printf("# Skipping object due to error file parsing error %s", err)
			continue
		}
		switch o := obj.(type) {
		case *v2.RouteTable:

			routeTables = append(routeTables, o)
		default:
			//skipping object
		}
	}
	newHTTPRoutes := map[string]*gwv1.HTTPRoute{}
	// for every routeTable in Gloo Mesh, create a HTTPRoute if the host exists
	for _, routeTable := range routeTables {
		for _, host := range routeTable.Spec.Hosts {
			// for each host in the route table find the httproute,
			//if one exists merge it with the existing one or create a new one
			httpRoute, key := o.FindHTTPRouteByHost(host)
			if httpRoute == nil {
				fmt.Printf("host %s not found in httproutes %s\n", host, routeTable.Name)
				continue
			}

			// check to see
			routeName := fmt.Sprintf("%s/%s", httpRoute.Namespace, httpRoute.Name)
			if newHTTPRoutes[routeName] != nil {
				newHTTPRoutes[routeName] = mergeHTTPRoutes(httpRoute.Name, newHTTPRoutes[routeName], httpRoute)
			} else {
				newHTTPRoutes[routeName] = &gwv1.HTTPRoute{
					TypeMeta: metav1.TypeMeta{
						Kind:       "HTTPRoute",
						APIVersion: "gateway.networking.k8s.io/v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      httpRoute.Name,
						Namespace: "gloo-system-nick",
					},
					Spec: gwv1.HTTPRouteSpec{
						CommonRouteSpec: httpRoute.Spec.CommonRouteSpec,
						Hostnames:       httpRoute.Spec.Hostnames,
						Rules:           httpRoute.Spec.Rules,
					},
				}
			}
			delete(o.httpRoutes, key)
		}

	}

	log.Printf("New Routes %d\n", len(newHTTPRoutes))
	log.Printf("Existing Routes After %d\n", len(o.httpRoutes))
	for _, newHTTPRoute := range newHTTPRoutes {
		o.AddRoute(newHTTPRoute)
	}
	log.Printf("TotalRoutes %d\n", len(o.httpRoutes))
	return nil
}

func mergeHTTPRoutes(routeName string, httpRoute1 *gwv1.HTTPRoute, httpRoute2 *gwv1.HTTPRoute) *gwv1.HTTPRoute {
	// combines hostnames and rules
	routeName1 := fmt.Sprintf("%s/%s", httpRoute1.Namespace, httpRoute1.Name)
	routeName2 := fmt.Sprintf("%s/%s", httpRoute2.Namespace, httpRoute2.Name)
	fmt.Printf("merging (%s) http route %s with route %s\n", routeName, routeName1, routeName2)
	newHTTPRoute := &gwv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: "gateway.networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      routeName,
			Namespace: "gloo-system-nick",
		},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: httpRoute1.Spec.CommonRouteSpec,
			//Hostnames:       httpRoute1.Spec.Hostnames, //hostnames should be the same
			Rules: httpRoute1.Spec.Rules,
		},
	}

	newHTTPRoute.Spec.Hostnames = append(newHTTPRoute.Spec.Hostnames, httpRoute2.Spec.Hostnames...)
	newHTTPRoute.Spec.Rules = append(newHTTPRoute.Spec.Rules, httpRoute2.Spec.Rules...)

	return newHTTPRoute

}

func (o *Outputs) FindHTTPRouteByHost(host string) (*gwv1.HTTPRoute, string) {
	for key, route := range o.httpRoutes {
		for _, h := range route.Spec.Hostnames {
			if string(h) == host {
				return route, key
			}
		}
	}
	return nil, ""
}

// depending on the subdomain index we will combine routes with a matching subdomain
func (o *Outputs) combineHTTPRoutes(subDomainIndex int) error {

	hostnamesMap := map[string][]string{}
	log.Printf("Current Routes %d\n", len(o.httpRoutes))

	for _, route := range o.httpRoutes {
		for _, h := range route.Spec.Hostnames {
			// Example subdomain index 3
			//something.subdomain.google.com
			hostSplit := strings.Split(string(h), ".")
			if len(hostSplit) < subDomainIndex {
				//its in its own hostnames map
				if hostnamesMap[string(h)] == nil {
					hostnamesMap[string(h)] = []string{string(h)}
				}
				hostnamesMap[string(h)] = append(hostnamesMap[string(h)], string(h))
				continue
			}
			//
			subdomain := strings.Join(hostSplit[len(hostSplit)-subDomainIndex:], ".")
			if hostnamesMap[subdomain] == nil {
				hostnamesMap[subdomain] = []string{subdomain}
			}
			hostnamesMap[subdomain] = append(hostnamesMap[subdomain], string(h))
		}
	}
	log.Printf("Routes after combination %d\n", len(hostnamesMap))
	return nil
}

func (o *Outputs) GenerateGenericRouteOptions() error {

	// its very common to have retries or extauth disabled, if we find these we will separate them out
	return nil
}
