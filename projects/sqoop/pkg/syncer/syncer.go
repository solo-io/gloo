package syncer

import (
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"context"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/translator"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
)

type Syncer struct {
	writeNamespace  string
	reporter        reporter.Reporter
	writeErrs       chan error
	proxyReconciler gloov1.ProxyReconciler
}

func (s *Syncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {
	ctx = contextutils.WithLogger(ctx, "syncer")

	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin sync %v (%v schemas %v resolverMaps)",
		snap.Hash(),
		len(snap.Schemas),
		len(snap.ResolverMaps),
	)
	defer logger.Infof("end sync %v", snap.Hash())
	logger.Debugf("%v", snap)

	proxy, resourceErrs := translator.Translate(s.writeNamespace, snap)
	if err := s.reporter.WriteReports(ctx, resourceErrs); err != nil {
		return errors.Wrapf(err, "writing reports")
	}
	if err := resourceErrs.Validate(); err != nil {
		logger.Warnf("gateway %v was rejected due to invalid config: %v\nxDS cache will not be updated.", err)
		return nil
	}
	// proxy was deleted / none desired
	if proxy == nil {
		return s.proxyReconciler.Reconcile(s.writeNamespace, nil, nil, clients.ListOpts{})
	}
	logger.Debugf("creating proxy %v", proxy.Metadata.Ref())
	if err := s.proxyReconciler.Reconcile(s.writeNamespace, gloov1.ProxyList{proxy}, nil, clients.ListOpts{}); err != nil {
		return err
	}

	/*
	1. create new graphql endpoints
	2. update router with graphql endpoints
	3. resourceErrs := configErrs / writeReports(resourceErrs)
	4. configure gloo
	# tip: use multierror instead of early returns (when possible)
	 */
	var endpoints []*router.Endpoint
	for _, schema := range snap.Schemas.List() {
		endpoint, err := createEndpointForSchema()
		if err != nil {
			// multierror / report
		}
		endpoints = append(endpoints, endpoint)
	}
	s.router.SetEndpoints(endpoints)

	// start propagating for new set of resources
	// TODO(ilackarms): reinstate propagator
	return nil // s.propagator.PropagateStatuses(snap, proxy, clients.WatchOpts{Ctx: ctx})

}

func desiredResources() (*gloov1.Proxy, []*v1.ResolverMap, error) {}

// for supporting multiple schemas on the same port, essentially
func createRouterEndpoints()                     {}
func createEndpointForSchema()                   {}
func generateSkeletonResolvermap()               {}
func generateEmptyResolvermap(schema *v1.Schema) {}
func configureGloo()                             {}
