package syncer

import (
	"context"

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
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/todo"
)

type Syncer struct {
	writeNamespace    string
	reporter          reporter.Reporter
	writeErrs         chan error
	proxyReconciler   gloov1.ProxyReconciler
	resolverMapClient v1.ResolverMapClient
	engine            *engine.Engine
	router            *router.Router
}

func NewSyncer(writeNamespace string,
	reporter reporter.Reporter,
	writeErrs chan error,
	proxyReconciler gloov1.ProxyReconciler,
	resolverMapClient v1.ResolverMapClient,
	engine *engine.Engine,
	router *router.Router) v1.ApiSyncer {
	s := &Syncer{
		writeNamespace:    writeNamespace,
		reporter:          reporter,
		writeErrs:         writeErrs,
		proxyReconciler:   proxyReconciler,
		resolverMapClient: resolverMapClient,
		engine:            engine,
		router:            router,
	}
	return s
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

	proxy := translator.Translate(s.writeNamespace, snap, resourceErrs)
	if err := s.reporter.WriteReports(ctx, resourceErrs); err != nil {
		return errors.Wrapf(err, "writing reports")
	}
	if err := resourceErrs.Validate(); err != nil {
		logger.Errorf("snapshot %v was rejected due to invalid config: %v", err)
		return nil
	}
	logger.Infof("creating proxy %v", proxy.Metadata.Ref())
	if err := s.proxyReconciler.Reconcile(s.writeNamespace, gloov1.ProxyList{proxy}, TODO.TransitionFunction, clients.ListOpts{
		Ctx: ctx,
		Selector: map[string]string{
			"created_by": "sqoop",
		},
	}); err != nil {
		return err
	}

	var endpoints []*router.Endpoint
	var resolverMapsToGenerate v1.ResolverMapList
	var schemasToUpdate v1.SchemaList
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
			schemasToUpdate = append(schemasToUpdate, schema)

			// nothing to do for this schema yet, need to receive some resolvers
			continue
		}
		resourceErrs.Accept(schema)

		resolverMap, err := snap.ResolverMaps.List().Find(schema.ResolverMap.Strings())
		if err != nil {
			resourceErrs.AddError(schema, errors.Wrapf(err, "finding resolvermap for schema"))
			continue
		}

		resourceErrs.Accept(resolverMap)

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
	for _, rm := range resolverMapsToGenerate {
		if _, err := s.resolverMapClient.Write(rm, clients.WriteOpts{}); err != nil {
			return errors.Wrapf(err, "writing generated resolver maps to storage")
		}
	}

	// start propagating for new set of resources
	// TODO(ilackarms): reinstate propagator
	return nil // s.propagator.PropagateStatuses(snap, proxy, clients.WatchOpts{Ctx: ctx})

}

func parseSchemaString(sch *v1.Schema) (*schema.Schema, error) {
	parsedSchema := schema.New()
	return parsedSchema, parsedSchema.Parse(sch.InlineSchema)
}
