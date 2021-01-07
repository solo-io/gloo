package failover

import (
	"context"

	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
)

const (
	PortNumber           = 15443
	PortName             = "failover"
	DownstreamSecretName = "failover-downstream"
	UpstreamSecretName   = "failover-upstream"
)

//go:generate mockgen -source ./interfaces.go -destination ./mocks/mock_interfaces.go

// FailoverProcessor calculates the output Upstream for a given change/event of a FailoverProcessor.
type FailoverProcessor interface {
	// ProcessFailoverUpdate returns a status similar to a typical error return value, if nil it can be ignored,
	// otherwise it should be saved back to the parent object.
	ProcessFailoverUpdate(
		ctx context.Context,
		obj *fedv1.FailoverScheme,
	) (*gloov1.Upstream, *fed_types.FailoverSchemeStatus)
	// ProcessFailoverDelete does not return a status as the object is being deleted.
	// May return nil, nil if primary is nil, as nothing can be done, and there is no "error"
	ProcessFailoverDelete(ctx context.Context, obj *fedv1.FailoverScheme) (*gloov1.Upstream, error)
}
