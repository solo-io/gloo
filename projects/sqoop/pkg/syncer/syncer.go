package syncer

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/qloo/pkg/graphql"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/engine"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/engine/router"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/translator"
	"github.com/vektah/gqlgen/neelance/schema"
)

type Syncer struct {
	writeNamespace        string
	reporter              reporter.Reporter
	writeErrs             chan error
	proxyReconciler       gloov1.ProxyReconciler
	resolverMapReconciler v1.ResolverMapReconciler
	engine                *engine.Engine
	router                *router.Router
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

	resourceErrs := make(reporter.ResourceErrors)
	resourceErrs.Initialize(snap.ResolverMaps.List().AsInputResources()...)
	resourceErrs.Initialize(snap.Schemas.List().AsInputResources()...)

	var endpoints []*router.Endpoint
	var resolverMapsToGenerate v1.ResolverMapList
	for _, schema := range snap.Schemas.List() {
		if schema.ResolverMap.Name == "" {
			newMeta := core.Metadata{
				Name:        schema.Metadata.Name,
				Namespace:   schema.Metadata.Namespace,
				Annotations: map[string]string{"created_for": schema.Metadata.Name},
			}
			parsedSchema, err := parseSchemaString(schema)
			if err != nil {
				resourceErrs.AddError(schema, errors.Wrapf(err, "failed to parse schema"))
				continue
			}

			rm := translator.GenerateResolverMapSkeleton(newMeta, parsedSchema)

			resolverMapsToGenerate = append(resolverMapsToGenerate, rm)

			// nothing to do for this schema yet, need to receive some resolvers
			continue
		}

		resolverMap, err := snap.ResolverMaps.List().Find(schema.ResolverMap.Strings())
		if err != nil {
			resourceErrs.AddError(schema, errors.Wrapf(err, "finding resolvermap for schema"))
			continue
		}

		endpoint, schemaErr, resolverErr := s.engine.CreateGraphqlEndpoint(schema, resolverMap)
		if schemaErr != nil {
			resourceErrs.AddError(schema, schemaErr)
		}
		if resolverErr != nil {
			resourceErrs.AddError(resolverMap, resolverErr)
		}
		if schemaErr != nil || resolverErr != nil {
			continue
		}
		endpoints = append(endpoints, endpoint)
	}
	s.router.UpdateEndpoints(endpoints)

	proxy, resolverMaps, resourceErrs := translator.Translate(s.writeNamespace, snap)
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
	for _, schema := range snap.Schemas.List() {
		resolverMap, err := snap.ResolverMaps.List().Find(schema.ResolverMap.Strings())
		if err != nil {

		}
	}

	// start propagating for new set of resources
	// TODO(ilackarms): reinstate propagator
	return nil // s.propagator.PropagateStatuses(snap, proxy, clients.WatchOpts{Ctx: ctx})

}

func (el *Syncer) update(cfg *v1.ApiSnapshot) error {
	endpoints, reports := el.createGraphqlEndpoints(cfg)
	el.router.UpdateEndpoints(endpoints...)
	errs := configErrs(reports)
	if err := el.reporter.WriteReports(reports); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := el.operator.ConfigureGloo(); err != nil {
		errs = multierror.Append(errs, err)
	}
	return errs
}

func (el *Syncer) createGraphqlEndpoints(cfg *v1.Config) ([]*graphql.Endpoint, []reporter.ConfigObjectReport) {
	var (
		endpoints          []*graphql.Endpoint
		schemaReports      []reporter.ConfigObjectReport
		resolverMapReports []reporter.ConfigObjectReport
	)
	resolverMapErrs := make(map[*v1.ResolverMap]error)

	for _, schema := range cfg.Schemas {
		schemaReport := reporter.ConfigObjectReport{
			CfgObject: schema,
		}
		// empty map means we should generate a skeleton and update the schema to point to it
		ep, schemaErr, resolverMapErr := el.handleSchema(schema, cfg.ResolverMaps)
		if schemaErr != nil {
			resolverMapErr.err = multierror.Append(resolverMapErr.err, errors.Wrap(schemaErr, "schema was not accepted"))
		}
		if resolverMapErr.resolverMap != nil {
			err := resolverMapErrs[resolverMapErr.resolverMap]
			if resolverMapErr.err != nil {
				err = multierror.Append(resolverMapErrs[resolverMapErr.resolverMap], resolverMapErr.err)
			}
			resolverMapErrs[resolverMapErr.resolverMap] = err
		}
		schemaReport.Err = schemaErr
		schemaReports = append(schemaReports, schemaReport)
		if ep == nil {
			continue
		}
		endpoints = append(endpoints, ep)
	}
	for resolverMap, err := range resolverMapErrs {
		resolverMapReports = append(resolverMapReports, reporter.ConfigObjectReport{
			CfgObject: resolverMap,
			Err:       err,
		})
	}
	return endpoints, append(schemaReports, resolverMapReports...)
}

type resolverMapError struct {
	resolverMap *v1.ResolverMap
	err         error
}

func (el *Syncer) handleSchema(schema *v1.Schema, resolvers []*v1.ResolverMap) (*graphql.Endpoint, error, resolverMapError) {
	if schema.ResolverMap == "" {
		return nil, el.createEmptyResolverMap(schema), resolverMapError{}
	}
	for _, resolverMap := range resolvers {
		if resolverMap.Name == schema.ResolverMap {
			ep, schemaErr, resolverErr := el.createGraphqlEndpoint(schema, resolverMap)
			return ep, schemaErr, resolverMapError{resolverMap: resolverMap, err: resolverErr}
		}
	}
	return nil, errors.Errorf("resolver map %v for schema %v not found", schema.ResolverMap, schema.Name), resolverMapError{}
}

func parseSchemaString(sch *v1.Schema) (*schema.Schema, error) {
	parsedSchema := schema.New()
	return parsedSchema, parsedSchema.Parse(sch.InlineSchema)
}
