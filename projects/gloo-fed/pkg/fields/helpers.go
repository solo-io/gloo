package fields

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ClusterIndex = "spec.cluster"
)

func BuildClusterFieldMatcher(cluster string) client.MatchingFields {
	return map[string]string{
		ClusterIndex: cluster,
	}
}

/*
	copied from: https://github.com/kubernetes-sigs/controller-runtime/blob/587ca7c550ef21706453ea71793c1b26d5666836/pkg/cache/internal/cache_reader.go#L171

	In order to set up indexers for the default controller-runtime cache, these string values need to be appended which
	are only contained in an internal package. So we have to copy paste them here.
*/

// FieldIndexName constructs the name of the index over the given field,
// for use with an indexer.
func FieldIndexName(field string) string {
	return "field:" + field
}

// noNamespaceNamespace is used as the "namespace" when we want to list across all namespaces
const allNamespacesNamespace = "__all_namespaces"

// KeyToNamespacedKey prefixes the given index key with a namespace
// for use in field selector indexes.
func KeyToNamespacedKey(ns string, baseKey string) string {
	if ns != "" {
		return ns + "/" + baseKey
	}
	return allNamespacesNamespace + "/" + baseKey
}
