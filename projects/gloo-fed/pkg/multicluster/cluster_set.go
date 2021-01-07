package multicluster

import (
	"context"
	"sort"
	"sync"

	"github.com/solo-io/skv2/pkg/multicluster"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

//go:generate mockgen -source ./cluster_set.go -destination ./mocks/mock_cluster_set.go

// ClusterSet records the names of clusters currently available via the skv2 Cluster Watcher it is registered with.
type ClusterSet interface {
	multicluster.ClusterHandler
	multicluster.ClusterRemovedHandler
	// ListClusters returns a list of all available clusters.
	ListClusters() []string
	// Exists returns true if the given cluster is available via the Cluster Watcher
	Exists(cluster string) bool
}

type set struct {
	clusters map[string]struct{}
	mutex    sync.RWMutex
}

func NewClusterSet() ClusterSet {
	return &set{
		clusters: make(map[string]struct{}),
		mutex:    sync.RWMutex{},
	}
}

func (s *set) AddCluster(_ context.Context, cluster string, _ manager.Manager) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.clusters[cluster] = struct{}{}
}

func (s *set) RemoveCluster(cluster string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.clusters, cluster)
}

func (s *set) ListClusters() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var output []string
	for cluster := range s.clusters {
		output = append(output, cluster)
	}
	sort.Strings(output)
	return output
}

func (s *set) Exists(cluster string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.clusters[cluster]
	return ok
}
