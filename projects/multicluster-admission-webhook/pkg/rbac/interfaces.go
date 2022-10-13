package rbac

import (
	"context"

	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
)

//go:generate mockgen -source ./interfaces.go -destination ./mocks/mock_placement.go

/*
	Extracts a list of Placements from a given raw k8s object.
	A valid placement must have non-empty namespace and cluster list, or it will be considered invalid.

	This is generated with skv2 using the template defined in ./codegen/parser.gotmpl.
*/
type Parser interface {
	Parse(ctx context.Context, rawObj []byte) ([]*multicluster_types.Placement, error)
}
