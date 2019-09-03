package rawgetter

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

//go:generate mockgen -destination mocks/mock_raw_getter.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter RawGetter

type RawGetter interface {
	GetRaw(ctx context.Context, in resources.InputResource, resourceCrd crd.Crd) *v1.Raw

	// init `emptyInputResource` from the resource provided in the yaml
	// enforces that the name/namespace in the yaml don't differ from the provided ref, and that the resource version is provided
	InitResourceFromYamlString(ctx context.Context, yamlString string, refToValidate *core.ResourceRef, emptyInputResource resources.InputResource) error
}
