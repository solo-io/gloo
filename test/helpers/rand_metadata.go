package helpers

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/bxcodec/faker"
)

func NewRandomMetadata() core.Metadata {
	meta := core.Metadata{}
	faker.FakeData(&meta)
	// dns label stuff
	meta.Name = "a"+RandString(6)+"a"
	meta.Namespace = "a"+RandString(6)+"a"
	return meta
}