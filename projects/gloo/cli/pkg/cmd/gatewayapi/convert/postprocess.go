package convert

import (
	"fmt"
	"strings"

	"k8s.io/utils/ptr"
	"sigs.k8s.io/gateway-api/apisx/v1alpha1"

	"github.com/golang/protobuf/proto"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/gatewayapi/convert/domain"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func (g *GatewayAPIOutput) PostProcess(opts *Options) error {

	// complete delegation
	if err := g.finishDelegation(); err != nil {
		return err
	}

	if opts.CombineRouteOptions {
		fmt.Printf("Combining route options...\n")
		g.combineRouteOptions()
	}
	if opts.IncludeUnknownResources {
		g.gatewayAPICache.YamlObjects = g.edgeCache.YamlObjects
	}

	// fix all the cel validation issues
	if err := g.celValidationCorrections(); err != nil {
		return err
	}

	return nil
}

// This function fixes all the cel validation rules that are cumbersome.
// All route rules with URL Rewrite need to be their own separate rules
// ListenerSets can only have a max of 64 listeners
func (g *GatewayAPIOutput) celValidationCorrections() error {
	fmt.Printf("Fixing CEL validations...\n")
	g.fixRewritesPerMatch()

	g.splitListenerSets()

	g.splitHTTPRouteRules()

	return nil
}

func (g *GatewayAPIOutput) splitHTTPRouteRules() {
	var httpRoutesToDelete []string
	var updatedHTTPRoutes []*domain.HTTPRouteWrapper
	for httpRouteKey, httpRoute := range g.gatewayAPICache.HTTPRoutes {
		if len(httpRoute.Spec.Rules) > 16 {
			// listener set needs to be broken up into multiple
			httpRoutesToDelete = append(httpRoutesToDelete, httpRouteKey)
			entries := splitRules(httpRoute.Spec.Rules, 16)
			g.AddErrorFromWrapper(ERROR_TYPE_CEL_VALIDATION_CORRECTION, httpRoute, "HTTPRoute contains too many route rules %d, splitting into %d new HTTPRoutes", len(httpRoute.Spec.Rules), len(entries))

			// for each entry set we create a new XListenerSet
			for i, entry := range entries {
				// new XListenerSet
				newHTTPRoute := httpRoute.DeepCopy()
				newHTTPRoute.Spec.Rules = entry
				newHTTPRoute.Name = fmt.Sprintf("%s-%d", httpRoute.Name, i)
				updatedHTTPRoutes = append(updatedHTTPRoutes, &domain.HTTPRouteWrapper{
					HTTPRoute:        newHTTPRoute,
					OriginalFileName: httpRoute.OriginalFileName,
				})
			}
		}
	}
	fmt.Printf("HTTPRoutes number of rules that required spliting: %d generated %d new routes\n", len(httpRoutesToDelete), len(updatedHTTPRoutes))

	for _, httpRouteKey := range httpRoutesToDelete {
		delete(g.gatewayAPICache.HTTPRoutes, httpRouteKey)
	}

	for _, httpRoute := range updatedHTTPRoutes {
		g.gatewayAPICache.AddHTTPRoute(httpRoute)
	}
}

func (g *GatewayAPIOutput) splitListenerSets() {
	var listenerSetsToDelete []string
	var updatedListenerSets []*domain.ListenerSetWrapper
	for listenerSetKey, listenerSet := range g.gatewayAPICache.ListenerSets {
		if len(listenerSet.Spec.Listeners) > 64 {
			// listener set needs to be broken up into multiple
			listenerSetsToDelete = append(listenerSetsToDelete, listenerSetKey)
			entries := splitListeners(listenerSet.Spec.Listeners, 64)
			g.AddErrorFromWrapper(ERROR_TYPE_CEL_VALIDATION_CORRECTION, listenerSet, "ListenerSet contains too many listeners %d, splitting into %d new ListenerSet", len(listenerSet.Spec.Listeners), len(entries))

			// for each entry set we create a new XListenerSet
			for i, entry := range entries {
				// new XListenerSet
				newListenerSet := listenerSet.DeepCopy()
				newListenerSet.Spec.Listeners = entry
				newListenerSet.Name = fmt.Sprintf("%s-%d", listenerSet.Name, i)
				updatedListenerSets = append(updatedListenerSets, &domain.ListenerSetWrapper{
					XListenerSet:     newListenerSet,
					OriginalFileName: listenerSet.OriginalFileName,
				})
			}
		}
	}
	fmt.Printf("ListenerSets number of listeners that required splitting: %d generated %d new listeners\n", len(listenerSetsToDelete), len(updatedListenerSets))

	for _, listenerSetKey := range listenerSetsToDelete {
		delete(g.gatewayAPICache.ListenerSets, listenerSetKey)
	}

	for _, listenerSet := range updatedListenerSets {
		g.gatewayAPICache.AddListenerSet(listenerSet)
	}
}

func splitRules(slice []gwv1.HTTPRouteRule, maxLen int) [][]gwv1.HTTPRouteRule {
	var result [][]gwv1.HTTPRouteRule
	for maxLen < len(slice) {
		slice, result = slice[maxLen:], append(result, slice[0:maxLen:maxLen])
	}
	result = append(result, slice)
	return result
}
func splitListeners(slice []v1alpha1.ListenerEntry, maxLen int) [][]v1alpha1.ListenerEntry {
	var result [][]v1alpha1.ListenerEntry
	for maxLen < len(slice) {
		slice, result = slice[maxLen:], append(result, slice[0:maxLen:maxLen])
	}
	result = append(result, slice)
	return result
}

func (g *GatewayAPIOutput) fixRewritesPerMatch() {
	var updatedHTTPRoutes []*domain.HTTPRouteWrapper
	for _, httpRoute := range g.gatewayAPICache.HTTPRoutes {
		//
		var updatedRules []gwv1.HTTPRouteRule
		for _, rr := range httpRoute.Spec.Rules {
			for _, filter := range rr.Filters {
				if filter.Type == gwv1.HTTPRouteFilterURLRewrite {
					// if there is more than one match then we need to split out the rules
					if len(rr.Matches) > 1 {

						// create new rules and add them to the updatedRules
						var splitRules []gwv1.HTTPRouteRule
						for i, match := range rr.Matches {
							ruleName := ptr.To(gwv1.SectionName(fmt.Sprintf("%v-%d", rr.Name, i)))
							if rr.Name == nil {
								ruleName = nil
							}
							splitRules = append(splitRules, gwv1.HTTPRouteRule{
								Name:               ruleName,
								Matches:            []gwv1.HTTPRouteMatch{match},
								Filters:            rr.Filters,
								BackendRefs:        rr.BackendRefs,
								Timeouts:           rr.Timeouts,
								Retry:              rr.Retry,
								SessionPersistence: rr.SessionPersistence,
							})
						}
						updatedRules = append(updatedRules, splitRules...)
					}
				}
			}
		}
		if len(updatedRules) > 0 {
			g.AddErrorFromWrapper(ERROR_TYPE_CEL_VALIDATION_CORRECTION, httpRoute, "updating HTTPRoute URLRewrite rules %d to conform to one rule per match, total new rules %d", len(httpRoute.Spec.Rules), len(updatedRules))
			httpRoute.Spec.Rules = updatedRules
			updatedHTTPRoutes = append(updatedHTTPRoutes, httpRoute)
		}
	}
	// update the routes
	for _, httpRoute := range updatedHTTPRoutes {
		g.gatewayAPICache.AddHTTPRoute(httpRoute)
	}
}

func (g *GatewayAPIOutput) finishDelegation() error {

	// for all edge routetables we need to go and update labels on the httproutes to support delegation
	updatedHTTPRoutes := map[string]*domain.HTTPRouteWrapper{}
	for _, rtt := range g.edgeCache.RouteTables {
		routesToUpdate := g.processRouteForDelegation(rtt.Spec.Routes)

		for _, r := range routesToUpdate {
			// check to see if we already matched on this httproute
			updatedHTTPRoute, found := updatedHTTPRoutes[domain.NameNamespaceIndex(r.Name, r.Namespace)]

			if found {
				delegateValue := updatedHTTPRoute.Labels["delegation.gateway.solo.io/label"]
				if delegateValue != r.Labels["delegation.gateway.solo.io/label"] {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, rtt, "HTTPRoute already selected by delegation label delegation.gateway.solo.io/label: %s but also selected by delegation.gateway.solo.io/label: %s", r.Labels["delegation.gateway.solo.io/label"], delegateValue)
					continue
				}
			}
			updatedHTTPRoutes[domain.NameNamespaceIndex(r.Name, r.Namespace)] = r
		}
	}
	for _, vs := range g.edgeCache.VirtualServices {
		routesToUpdate := g.processRouteForDelegation(vs.Spec.VirtualHost.Routes)

		for _, r := range routesToUpdate {
			// check to see if we already matched on this httproute
			updatedHTTPRoute, found := updatedHTTPRoutes[domain.NameNamespaceIndex(r.Name, r.Namespace)]

			if found {
				delegateValue := updatedHTTPRoute.Labels["delegation.gateway.solo.io/label"]
				if delegateValue != r.Labels["delegation.gateway.solo.io/label"] {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, r, "HTTPRoute already selected by delegation label delegation.gateway.solo.io/label: %s but also selected by delegation.gateway.solo.io/label: %s", r.Labels["delegation.gateway.solo.io/label"], delegateValue)
					continue
				}
			}
			updatedHTTPRoutes[domain.NameNamespaceIndex(r.Name, r.Namespace)] = r
		}

	}
	for name, route := range updatedHTTPRoutes {
		g.gatewayAPICache.HTTPRoutes[name] = route
	}

	return nil
}

func (g *GatewayAPIOutput) processRouteForDelegation(routes []*v1.Route) []*domain.HTTPRouteWrapper {
	var routesToUpdate []*domain.HTTPRouteWrapper
	for _, rt := range routes {
		if rt.GetDelegateAction() != nil && rt.GetDelegateAction().GetSelector() != nil {
			selector := rt.GetDelegateAction().GetSelector()
			// go find the RouteTabels by this selector
			namespaces := selector.GetNamespaces()
			if namespaces != nil || len(selector.GetNamespaces()) == 0 {
				// default namespace is gloo-system
				namespaces = []string{"gloo-system"}
			}

			// get all http routes
			for _, httpRoute := range g.gatewayAPICache.HTTPRoutes {
				value, matches := routeMatchSelector(httpRoute, selector)
				if matches {

					if httpRoute.Labels == nil {
						httpRoute.Labels = map[string]string{}
					}
					// this value should match the backend ref name
					//   # Delegate to routes with the label delegation.gateway.solo.io/label:foobar
					//   - group: delegation.gateway.solo.io
					//     kind: label # no other value is allowed
					//     name: single.example.com # label value
					//     namespace: httpbin # defaults to parent's namespace if unset
					httpRoute.Labels["delegation.gateway.solo.io/label"] = value
					routesToUpdate = append(routesToUpdate, httpRoute)
				}
			}
		}
	}
	return routesToUpdate
}

func routeMatchSelector(route *domain.HTTPRouteWrapper, selector *v1.RouteTableSelector) (string, bool) {

	//check namespace first
	if namespaceMatch(route.Namespace, selector.Namespaces) {
		// check to see if any of the labels match the selector
		for k, v := range selector.Labels {
			// see if the route has the label key in the selector
			value, match := route.Labels[k]
			if match {
				// check to see if the values are the same
				if v == value {
					return v, true
				}
			}
		}
	}

	return "", false
}

func namespaceMatch(namespace string, namespaces []string) bool {
	for _, ns := range namespaces {
		if ns == "*" {
			return true
		}
		if ns == namespace {
			return true
		}
	}
	return false
}

// TODO we should only combine route options in the same namespace
func (g *GatewayAPIOutput) combineRouteOptions() {
	totalRouteOptions := len(g.gatewayAPICache.RouteOptions)
	var routeOptionKeys []string
	for key, _ := range g.gatewayAPICache.RouteOptions {
		routeOptionKeys = append(routeOptionKeys, key)
	}
	duplicates := map[string][]string{}

	// go through each namespace and only work on ones that match
	for _, primaryKey := range routeOptionKeys {
		primaryKeyNs := strings.Split(primaryKey, "/")

		for _, secondaryKey := range routeOptionKeys {
			secondaryKeyNs := strings.Split(secondaryKey, "/")
			if primaryKeyNs[0] != secondaryKeyNs[0] {
				// skip all keys not in the same namespace
				continue
			}
			if primaryKey == secondaryKey {
				// skip if its the same primaryKey
				continue
			}
			ro, found1 := g.gatewayAPICache.RouteOptions[primaryKey]
			if !found1 {
				// this primary primaryKey has already been removed
				//fmt.Printf("primary key %s not found\n", primaryKey)
				break
			}
			ro2, found2 := g.gatewayAPICache.RouteOptions[secondaryKey]
			if !found2 {
				// move on to the next secondaryKey
				continue
			}

			if proto.Equal(&ro.Spec, &ro2.Spec) {
				duplicates[primaryKey] = append(duplicates[primaryKey], secondaryKey)
				//fmt.Printf("Route Option %s matches %s\n", primaryKey, secondaryKey)
				// remove both of them from the list
				delete(g.gatewayAPICache.RouteOptions, secondaryKey)
			}
		}
	}

	// for every duplicate we need to create a new name and then do a replace
	replacementMap := map[string]string{}
	combined := 0
	for primaryKey, dups := range duplicates {
		newName := fmt.Sprintf("shared-%s", RandStringRunes(8))

		replacementMap[primaryKey] = newName

		//fmt.Printf("Duplicate Key: %s\n", primaryKey)
		for _, dup := range dups {
			//fmt.Printf("\t%d. Match: %s\n", i, dup)
			replacementMap[dup] = newName
		}

		// create a new RouteOption with the new name
		existingRO := g.gatewayAPICache.RouteOptions[primaryKey]
		existingRO.Name = newName
		g.gatewayAPICache.AddRouteOption(existingRO)
		delete(g.gatewayAPICache.RouteOptions, primaryKey)
		combined++
	}

	for key, route := range g.gatewayAPICache.HTTPRoutes {
		var newRules []gwv1.HTTPRouteRule
		for _, rule := range route.Spec.Rules {
			for _, filter := range rule.Filters {
				if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == "RouteOption" {
					nameNamespace := fmt.Sprintf("%s/%s", route.Namespace, filter.ExtensionRef.Name)
					newName, found := replacementMap[nameNamespace]
					if found {
						filter.ExtensionRef.Name = gwv1.ObjectName(newName)
					}
				}
			}
			newRules = append(newRules, rule)
		}
		route.Spec.Rules = newRules
		g.gatewayAPICache.HTTPRoutes[key] = route
	}
	fmt.Printf("Initial %d RouteOptions combined to %d\n", totalRouteOptions, combined)
}
