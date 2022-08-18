package extauth

import (
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"

	"k8s.io/apimachinery/pkg/util/sets"
)

// ConfigStore holds state for the ExtAuth plugin that must be shared across functions
type ConfigStore struct {
	httpListenerEntries map[*v1.HttpListener]*HttpListenerEntry
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		httpListenerEntries: map[*v1.HttpListener]*HttpListenerEntry{},
	}
}

type HttpListenerEntry struct {
	extAuthFilters        []*envoyauth.ExtAuthz
	filterConstructionErr error
	metadataNamespaces    sets.String
}

func (c *ConfigStore) getStagedFiltersForHttpListener(httpListener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	entry := c.httpListenerEntries[httpListener]
	if entry == nil {
		return nil, nil
	}

	if entry.filterConstructionErr != nil {
		return nil, entry.filterConstructionErr
	}

	return extauth.BuildStagedHttpFilters(func() ([]*envoyauth.ExtAuthz, error) {
		finalizedFilters := make([]*envoyauth.ExtAuthz, 0, len(entry.extAuthFilters))

		// We intentionally sort the additional metadata namespaces to keep envoy snapshot deterministic
		sortedMetadataNamespaces := entry.metadataNamespaces.List()

		for _, f := range entry.extAuthFilters {
			finalizedFilter := f

			finalizedFilter.MetadataContextNamespaces = append(finalizedFilter.MetadataContextNamespaces, sortedMetadataNamespaces...)
			finalizedFilters = append(finalizedFilters, finalizedFilter)
		}

		return finalizedFilters, nil
	}, extauth.FilterStage)

}

func (c *ConfigStore) getFiltersForHttpListener(httpListener *v1.HttpListener) ([]*envoyauth.ExtAuthz, error) {
	entry := c.httpListenerEntries[httpListener]
	if entry != nil {
		return entry.extAuthFilters, entry.filterConstructionErr
	}
	return nil, nil
}

func (c *ConfigStore) hasFiltersForHttpListener(httpListener *v1.HttpListener) bool {
	entry := c.httpListenerEntries[httpListener]
	return entry != nil && len(entry.extAuthFilters) > 0

}

func (c *ConfigStore) setFiltersForHttpListener(httpListener *v1.HttpListener, filters []*envoyauth.ExtAuthz, err error) {
	entry := c.httpListenerEntries[httpListener]
	if entry != nil {
		entry.extAuthFilters = filters
		entry.filterConstructionErr = err
		return
	}

	c.httpListenerEntries[httpListener] = &HttpListenerEntry{
		extAuthFilters:        filters,
		filterConstructionErr: err,
		metadataNamespaces:    sets.NewString(),
	}
}

func (c *ConfigStore) appendMetadataNamespacesForHttpListener(httpListener *v1.HttpListener, ns []string) {
	entry := c.httpListenerEntries[httpListener]
	if entry != nil {
		entry.metadataNamespaces.Insert(ns...)
		return
	}

	c.httpListenerEntries[httpListener] = &HttpListenerEntry{
		extAuthFilters:        nil,
		filterConstructionErr: nil,
		metadataNamespaces:    sets.NewString(ns...),
	}
}
