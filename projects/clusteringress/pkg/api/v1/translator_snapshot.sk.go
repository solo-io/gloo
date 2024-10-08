// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"fmt"
	"hash"
	"hash/fnv"
	"log"

	github_com_solo_io_gloo_projects_clusteringress_pkg_api_external_knative "github.com/solo-io/gloo/projects/clusteringress/pkg/api/external/knative"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type TranslatorSnapshot struct {
	Clusteringresses github_com_solo_io_gloo_projects_clusteringress_pkg_api_external_knative.ClusterIngressList
}

func (s TranslatorSnapshot) Clone() TranslatorSnapshot {
	return TranslatorSnapshot{
		Clusteringresses: s.Clusteringresses.Clone(),
	}
}

func (s TranslatorSnapshot) Hash(hasher hash.Hash64) (uint64, error) {
	if hasher == nil {
		hasher = fnv.New64()
	}
	if _, err := s.hashClusteringresses(hasher); err != nil {
		return 0, err
	}
	return hasher.Sum64(), nil
}

func (s TranslatorSnapshot) hashClusteringresses(hasher hash.Hash64) (uint64, error) {
	return hashutils.HashAllSafe(hasher, s.Clusteringresses.AsInterfaces()...)
}

func (s TranslatorSnapshot) HashFields() []zap.Field {
	var fields []zap.Field
	hasher := fnv.New64()
	ClusteringressesHash, err := s.hashClusteringresses(hasher)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	fields = append(fields, zap.Uint64("clusteringresses", ClusteringressesHash))
	snapshotHash, err := s.Hash(hasher)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	return append(fields, zap.Uint64("snapshotHash", snapshotHash))
}

func (s *TranslatorSnapshot) GetResourcesList(resource resources.Resource) (resources.ResourceList, error) {
	switch resource.(type) {
	case *github_com_solo_io_gloo_projects_clusteringress_pkg_api_external_knative.ClusterIngress:
		return s.Clusteringresses.AsResources(), nil
	default:
		return resources.ResourceList{}, eris.New("did not contain the input resource type returning empty list")
	}
}

func (s *TranslatorSnapshot) RemoveFromResourceList(resource resources.Resource) error {
	refKey := resource.GetMetadata().Ref().Key()
	switch resource.(type) {
	case *github_com_solo_io_gloo_projects_clusteringress_pkg_api_external_knative.ClusterIngress:

		for i, res := range s.Clusteringresses {
			if refKey == res.GetMetadata().Ref().Key() {
				s.Clusteringresses = append(s.Clusteringresses[:i], s.Clusteringresses[i+1:]...)
				break
			}
		}
		return nil
	default:
		return eris.Errorf("did not remove the resource because its type does not exist [%T]", resource)
	}
}

func (s *TranslatorSnapshot) RemoveMatches(predicate core.Predicate) {
	var Clusteringresses github_com_solo_io_gloo_projects_clusteringress_pkg_api_external_knative.ClusterIngressList
	for _, res := range s.Clusteringresses {
		if matches := predicate(res.GetMetadata()); !matches {
			Clusteringresses = append(Clusteringresses, res)
		}
	}
	s.Clusteringresses = Clusteringresses
}

func (s *TranslatorSnapshot) UpsertToResourceList(resource resources.Resource) error {
	refKey := resource.GetMetadata().Ref().Key()
	switch typed := resource.(type) {
	case *github_com_solo_io_gloo_projects_clusteringress_pkg_api_external_knative.ClusterIngress:
		updated := false
		for i, res := range s.Clusteringresses {
			if refKey == res.GetMetadata().Ref().Key() {
				s.Clusteringresses[i] = typed
				updated = true
			}
		}
		if !updated {
			s.Clusteringresses = append(s.Clusteringresses, typed)
		}
		s.Clusteringresses.Sort()
		return nil
	default:
		return eris.Errorf("did not add/replace the resource type because it does not exist %T", resource)
	}
}

type TranslatorSnapshotStringer struct {
	Version          uint64
	Clusteringresses []string
}

func (ss TranslatorSnapshotStringer) String() string {
	s := fmt.Sprintf("TranslatorSnapshot %v\n", ss.Version)

	s += fmt.Sprintf("  Clusteringresses %v\n", len(ss.Clusteringresses))
	for _, name := range ss.Clusteringresses {
		s += fmt.Sprintf("    %v\n", name)
	}

	return s
}

func (s TranslatorSnapshot) Stringer() TranslatorSnapshotStringer {
	snapshotHash, err := s.Hash(nil)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	return TranslatorSnapshotStringer{
		Version:          snapshotHash,
		Clusteringresses: s.Clusteringresses.NamespacesDotNames(),
	}
}

var TranslatorGvkToHashableResource = map[schema.GroupVersionKind]func() resources.HashableResource{
	github_com_solo_io_gloo_projects_clusteringress_pkg_api_external_knative.ClusterIngressGVK: github_com_solo_io_gloo_projects_clusteringress_pkg_api_external_knative.NewClusterIngressHashableResource,
}
