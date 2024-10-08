// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"fmt"
	"hash"
	"hash/fnv"
	"log"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type EdsSnapshot struct {
	Upstreams UpstreamList
}

func (s EdsSnapshot) Clone() EdsSnapshot {
	return EdsSnapshot{
		Upstreams: s.Upstreams.Clone(),
	}
}

func (s EdsSnapshot) Hash(hasher hash.Hash64) (uint64, error) {
	if hasher == nil {
		hasher = fnv.New64()
	}
	if _, err := s.hashUpstreams(hasher); err != nil {
		return 0, err
	}
	return hasher.Sum64(), nil
}

func (s EdsSnapshot) hashUpstreams(hasher hash.Hash64) (uint64, error) {
	return hashutils.HashAllSafe(hasher, s.Upstreams.AsInterfaces()...)
}

func (s EdsSnapshot) HashFields() []zap.Field {
	var fields []zap.Field
	hasher := fnv.New64()
	UpstreamsHash, err := s.hashUpstreams(hasher)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	fields = append(fields, zap.Uint64("upstreams", UpstreamsHash))
	snapshotHash, err := s.Hash(hasher)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	return append(fields, zap.Uint64("snapshotHash", snapshotHash))
}

func (s *EdsSnapshot) GetResourcesList(resource resources.Resource) (resources.ResourceList, error) {
	switch resource.(type) {
	case *Upstream:
		return s.Upstreams.AsResources(), nil
	default:
		return resources.ResourceList{}, eris.New("did not contain the input resource type returning empty list")
	}
}

func (s *EdsSnapshot) RemoveFromResourceList(resource resources.Resource) error {
	refKey := resource.GetMetadata().Ref().Key()
	switch resource.(type) {
	case *Upstream:

		for i, res := range s.Upstreams {
			if refKey == res.GetMetadata().Ref().Key() {
				s.Upstreams = append(s.Upstreams[:i], s.Upstreams[i+1:]...)
				break
			}
		}
		return nil
	default:
		return eris.Errorf("did not remove the resource because its type does not exist [%T]", resource)
	}
}

func (s *EdsSnapshot) RemoveMatches(predicate core.Predicate) {
	var Upstreams UpstreamList
	for _, res := range s.Upstreams {
		if matches := predicate(res.GetMetadata()); !matches {
			Upstreams = append(Upstreams, res)
		}
	}
	s.Upstreams = Upstreams
}

func (s *EdsSnapshot) UpsertToResourceList(resource resources.Resource) error {
	refKey := resource.GetMetadata().Ref().Key()
	switch typed := resource.(type) {
	case *Upstream:
		updated := false
		for i, res := range s.Upstreams {
			if refKey == res.GetMetadata().Ref().Key() {
				s.Upstreams[i] = typed
				updated = true
			}
		}
		if !updated {
			s.Upstreams = append(s.Upstreams, typed)
		}
		s.Upstreams.Sort()
		return nil
	default:
		return eris.Errorf("did not add/replace the resource type because it does not exist %T", resource)
	}
}

type EdsSnapshotStringer struct {
	Version   uint64
	Upstreams []string
}

func (ss EdsSnapshotStringer) String() string {
	s := fmt.Sprintf("EdsSnapshot %v\n", ss.Version)

	s += fmt.Sprintf("  Upstreams %v\n", len(ss.Upstreams))
	for _, name := range ss.Upstreams {
		s += fmt.Sprintf("    %v\n", name)
	}

	return s
}

func (s EdsSnapshot) Stringer() EdsSnapshotStringer {
	snapshotHash, err := s.Hash(nil)
	if err != nil {
		log.Println(eris.Wrapf(err, "error hashing, this should never happen"))
	}
	return EdsSnapshotStringer{
		Version:   snapshotHash,
		Upstreams: s.Upstreams.NamespacesDotNames(),
	}
}

var EdsGvkToHashableResource = map[schema.GroupVersionKind]func() resources.HashableResource{
	UpstreamGVK: NewUpstreamHashableResource,
}
