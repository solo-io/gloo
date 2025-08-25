package selectionutils

import (
	"context"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	sk_errors "github.com/solo-io/solo-kit/pkg/errors"
)

type RouteTableSelector interface {
	SelectOrBuildRouteTable(ctx context.Context, ref *core.ResourceRef) (*gatewayv1.RouteTable, error)
}

type routeTableSelector struct {
	client       gatewayv1.RouteTableClient
	podNamespace string
}

var _ RouteTableSelector = &routeTableSelector{}

func NewRouteTableSelector(client gatewayv1.RouteTableClient, podNamespace string) *routeTableSelector {
	return &routeTableSelector{
		client:       client,
		podNamespace: podNamespace,
	}
}

func (s *routeTableSelector) SelectOrBuildRouteTable(ctx context.Context, ref *core.ResourceRef) (*gatewayv1.RouteTable, error) {
	logger := contextutils.LoggerFrom(ctx)
	logger.Debugw("Selecting or building RouteTable",
		"issue", "8539",
		"ref", ref.String(),
		"podNamespace", s.podNamespace)

	// Read or build route table
	// unlike virtual service, name must be provided as there is no "default" virtual service
	name := ref.GetName()
	if name == "" {
		logger.Debugw("RouteTable name is empty",
			"issue", "8539",
			"ref", ref.String())
		return nil, eris.New("must provide a name for the target route table")
	}

	ns := ref.GetNamespace()
	if ns == "" {
		ns = defaults.GlooSystem
		logger.Debugw("RouteTable namespace was empty, using default",
			"issue", "8539",
			"originalRef", ref.String(),
			"defaultNamespace", ns)
	}

	logger.Debugw("Attempting to read RouteTable from client",
		"issue", "8539",
		"namespace", ns,
		"name", name,
		"fullRef", ref.String())

	found, err := s.client.Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: ctx})
	if err != nil && !sk_errors.IsNotExist(err) {
		logger.Debugw("RouteTable client read failed with error",
			"issue", "8539",
			"ref", ref.String(),
			"error", err.Error(),
			"isNotExistError", sk_errors.IsNotExist(err))
		return nil, err
	}
	if found != nil {
		logger.Debugw("RouteTable found successfully",
			"issue", "8539",
			"ref", ref.String(),
			"foundNamespace", found.GetMetadata().GetNamespace(),
			"foundName", found.GetMetadata().GetName())
		return found, nil
	}

	logger.Debugw("RouteTable not found, building default",
		"issue", "8539",
		"ref", ref.String())
	// Build a new default route table object
	return s.build(ctx, ref)
}

func (s *routeTableSelector) build(_ context.Context, ref *core.ResourceRef) (*gatewayv1.RouteTable, error) {
	routeTable := &gatewayv1.RouteTable{
		Metadata: &core.Metadata{
			Namespace: ref.GetNamespace(),
			Name:      ref.GetName(),
		},
	}
	if routeTable.GetMetadata().GetNamespace() == "" {
		routeTable.GetMetadata().Namespace = s.podNamespace
	}
	if routeTable.GetMetadata().GetName() == "" {
		routeTable.GetMetadata().Name = "default"
	}

	return routeTable, nil
}
