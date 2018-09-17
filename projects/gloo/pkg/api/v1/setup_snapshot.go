package v1

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type SetupSnapshot struct {
	Settings SettingsByNamespace
}

func (s SetupSnapshot) Clone() SetupSnapshot {
	return SetupSnapshot{
		Settings: s.Settings.Clone(),
	}
}

func (s SetupSnapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, settings := range snapshotForHashing.Settings.List() {
		resources.UpdateMetadata(settings, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}
