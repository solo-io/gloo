package placement

import (
	"context"

	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/test/internal/api/test.multicluster.solo.io/v1alpha1"
)

type typedParser struct{}

func NewTypedParser() TypedParser {
	return &typedParser{}
}

func (t *typedParser) ParseTest(_ context.Context, obj *v1alpha1.Test) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{
		{
			Namespaces: obj.Spec.GetNamespaces(),
			Clusters:   obj.Spec.GetClusters(),
		},
	}, nil
}
