package testutils

import (
	"context"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/check"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/printers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// CheckResourcesOk checks if the resources are ok, similar to `glooctl check`
func CheckResourcesOk(ctx context.Context, installNamespace string) error {
	contextWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()
	opts := &options.Options{
		Metadata: core.Metadata{
			Namespace: installNamespace,
		},
		Top: options.Top{
			Ctx: contextWithCancel,
		},
	}
	return check.CheckResources(contextWithCancel, printers.P{}, opts)
}

// IsGlooInstalled checks if gloo is installed by checking the resources
func IsGlooInstalled(ctx context.Context, installNamespace string) bool {
	return CheckResourcesOk(ctx, installNamespace) == nil
}
