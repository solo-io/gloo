package v1

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type ApiSnapshot struct {
	Artifacts ArtifactsByNamespace
	Endpoints EndpointsByNamespace
	Proxies   ProxiesByNamespace
	Secrets   SecretsByNamespace
	Upstreams UpstreamsByNamespace
}

func (s ApiSnapshot) Clone() ApiSnapshot {
	return ApiSnapshot{
		Artifacts: s.Artifacts.Clone(),
		Endpoints: s.Endpoints.Clone(),
		Proxies:   s.Proxies.Clone(),
		Secrets:   s.Secrets.Clone(),
		Upstreams: s.Upstreams.Clone(),
	}
}

func (s ApiSnapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, artifact := range snapshotForHashing.Artifacts.List() {
		resources.UpdateMetadata(artifact, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	for _, endpoint := range snapshotForHashing.Endpoints.List() {
		resources.UpdateMetadata(endpoint, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	for _, proxy := range snapshotForHashing.Proxies.List() {
		resources.UpdateMetadata(proxy, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		proxy.SetStatus(core.Status{})
	}
	for _, secret := range snapshotForHashing.Secrets.List() {
		resources.UpdateMetadata(secret, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	for _, upstream := range snapshotForHashing.Upstreams.List() {
		resources.UpdateMetadata(upstream, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		upstream.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}
