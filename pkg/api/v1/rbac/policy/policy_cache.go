package policy

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type Snapshot struct {
	PolicyList   []*Policy
	IdentityList []*Identity
}

func (s Snapshot) Clone() Snapshot {
	var policyList []*Policy
	for _, policy := range s.PolicyList {
		policyList = append(policyList, proto.Clone(policy).(*Policy))
	}
	var identityList []*Identity
	for _, identity := range s.IdentityList {
		identityList = append(identityList, proto.Clone(identity).(*Identity))
	}
	return Snapshot{
		PolicyList:   policyList,
		IdentityList: identityList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, policy := range snapshotForHashing.PolicyList {
		resources.UpdateMetadata(policy, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		policy.SetStatus(core.Status{})
	}
	for _, identity := range snapshotForHashing.IdentityList {
		resources.UpdateMetadata(identity, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		identity.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}

type Cache interface {
	Register() error
	Policy() PolicyClient
	Identity() IdentityClient
	Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(policyClient PolicyClient, identityClient IdentityClient) Cache {
	return &cache{
		policy:   policyClient,
		identity: identityClient,
	}
}

type cache struct {
	policy   PolicyClient
	identity IdentityClient
}

func (c *cache) Register() error {
	if err := c.policy.Register(); err != nil {
		return err
	}
	if err := c.identity.Register(); err != nil {
		return err
	}
	return nil
}

func (c *cache) Policy() PolicyClient {
	return c.policy
}

func (c *cache) Identity() IdentityClient {
	return c.identity
}

func (c *cache) Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
	snapshots := make(chan *Snapshot)
	errs := make(chan error)

	currentSnapshot := Snapshot{}

	sync := func(newSnapshot Snapshot) {
		if currentSnapshot.Hash() == newSnapshot.Hash() {
			return
		}
		currentSnapshot = newSnapshot
		snapshots <- &currentSnapshot
	}
	policyChan, policyErrs, err := c.policy.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Policy watch")
	}
	identityChan, identityErrs, err := c.identity.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting Identity watch")
	}

	go func() {
		for {
			select {
			case policyList := <-policyChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.PolicyList = policyList
				sync(newSnapshot)
			case identityList := <-identityChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.IdentityList = identityList
				sync(newSnapshot)
			case err := <-policyErrs:
				errs <- err
			case err := <-identityErrs:
				errs <- err
			}
		}
	}()
	return snapshots, errs, nil
}
