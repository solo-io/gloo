package translator

import (
	"fmt"
	"sort"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"

	"github.com/envoyproxy/go-control-plane/api/filter/network"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/module"
	"github.com/solo-io/glue/pkg/translator/plugin"
)

type Translator struct {
	plugins []plugin.Plugin
}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) Translate(cfg v1.Config,
	clusters module.EndpointGroups,
	secretMap module.SecretMap) (envoycache.Snapshot, error) {

	/*

		create all clusters, and run the filters on all clusters.
		if from some reason a cluster has errored, send it back to user. and remove it
		from the list

	*/

	/*
		Create virtual hosts and ssl certificates and the such.
		for each virtual host, go over it's routes and:
			Create all routes inline, and then send them to be augmented by all filters
	*/

	// runTranslation

	// combine with cluster + endpoints
	// stable sort

	// computer snapshort version
	return envoycache.Snapshot{}, fmt.Errorf("not implemented")
}

func (t *Translator) getAllDependencies() []string {

}

func (t *Translator) runValidation() []error {

}

func (t *Translator) runTranslation(cfg v1.Config, secretMap module.SecretMap) envoycache.Snapshot {
	// compute virtual VirtualHosts
	// compute Routes
	// ...

	// do a stable sort
	var filters []plugin.FilterWrapper
	var routes []plugin.RouteWrapper
	var clusters []plugin.ClusterWrapper

	for _, plgin := range t.plugins {
		resource := plgin.Translate(cfg, secretMap)
		filters = append(filters, resource.Filters...)
		clusters = append(clusters, resource.Clusters...)
		routes = append(routes, resource.Routes...)
	}

	// for each route, find out which upstream it goes to and add metadata as appropriate.

	// sort out the filters
	sortedFilters := sortFilters(filters)

	snapshot := envoycache.NewSnapshot

}

func sortFilters(filters []plugin.FilterWrapper) []network.HttpFilter {
	// sort them accoirding to stage and then according to the name.
	less := func(i, j int) bool {
		filteri := filters[i]
		filterj := filters[j]
		if filteri.Stage != filterj.Stage {
			return filteri.Stage < filterj.Stage
		}
		return filteri.Filter.Name < filterj.Filter.Name
	}
	sort.Slice(filters, less)

	var sortedFilters []network.HttpFilter
	for _, filter := range filters {
		sortedFilters = append(sortedFilters, filter.Filter)
	}

	return sortedFilters
}
