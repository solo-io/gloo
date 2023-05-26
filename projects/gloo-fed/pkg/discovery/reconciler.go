package discovery

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/pkg/ezkube"
	"github.com/solo-io/skv2/pkg/verifier"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/input"
	hub_v1sets "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/sets"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery/translator"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	SnapshotName = "discovery"
)

type GlooResourceReconciler interface {
	ReconcileAll(id ezkube.ClusterResourceId) (bool, error)
}

type glooResourceReconcilerImpl struct {
	ctx                context.Context
	glooInstanceClient fedv1.GlooInstanceClient
	fedInputSnapshot   input.Builder
	translator         translator.GlooInstanceTranslator
	verifier           verifier.ServerResourceVerifier
}

func (g *glooResourceReconcilerImpl) ReconcileAll(
	_ ezkube.ClusterResourceId,
) (bool, error) { // The first bool indicates whether we should retry
	// Build a snapshot that gets everything across all clusters.
	// Translate that snapshot into a list of Gloo Instances.
	logger := contextutils.LoggerFrom(g.ctx)
	logger.Debug("Gloo Fed discovery - ReconcileAll")
	// the verifier allows us to ignore missing CRDs on remote clusters when we're gathering all the
	// resources to produce the input snapshot (which will be used to construct the GlooInstances)
	discoverySnapshot, err := g.fedInputSnapshot.BuildSnapshot(g.ctx, SnapshotName, input.BuildOptions{
		Gateways:              input.ResourceBuildOptions{Verifier: g.verifier},
		MatchableHttpGateways: input.ResourceBuildOptions{Verifier: g.verifier},
		MatchableTcpGateways:  input.ResourceBuildOptions{Verifier: g.verifier},
		VirtualServices:       input.ResourceBuildOptions{Verifier: g.verifier},
		RouteTables:           input.ResourceBuildOptions{Verifier: g.verifier},
		Upstreams:             input.ResourceBuildOptions{Verifier: g.verifier},
		UpstreamGroups:        input.ResourceBuildOptions{Verifier: g.verifier},
		Settings:              input.ResourceBuildOptions{Verifier: g.verifier},
		Proxies:               input.ResourceBuildOptions{Verifier: g.verifier},
		AuthConfigs:           input.ResourceBuildOptions{Verifier: g.verifier},
		RateLimitConfigs:      input.ResourceBuildOptions{Verifier: g.verifier},
	})
	if err != nil {
		logger.Errorw("Gloo Fed discovery - BuildSnapshot failed", zap.Error(err))
		return true, err
	}
	instances := g.translator.FromSnapshot(g.ctx, discoverySnapshot)
	newSet := hub_v1sets.NewGlooInstanceSet(instances...)

	existingInstances, err := g.glooInstanceClient.ListGlooInstance(g.ctx)
	if err != nil {
		logger.Errorw("Gloo Fed discovery - ListGlooInstance failed", zap.Error(err))
		return true, err
	}

	existingSet := hub_v1sets.NewGlooInstanceSet()
	for _, existingIter := range existingInstances.Items {
		existing := existingIter
		existingSet.Insert(&existing)
	}

	for _, instance := range newSet.List() {
		if err = g.glooInstanceClient.UpsertGlooInstance(g.ctx, instance); IgnoreAlreadyExists(err) != nil {
			logger.Errorw("Gloo Fed discovery - UpsertGlooInstance failed", zap.Error(err))
			return true, err
		}
	}

	for _, deletedIter := range existingSet.Difference(newSet).List() {
		deleted := deletedIter
		if err = g.glooInstanceClient.DeleteGlooInstance(g.ctx, client.ObjectKey{
			Namespace: deleted.GetNamespace(),
			Name:      deleted.GetName(),
		}); err != nil {
			logger.Errorw("Gloo Fed discovery - DeleteGlooInstance failed", zap.Error(err))
			return true, err
		}
	}
	// return false because there is no need to requeue if the reconcile is successful
	// (if set to true, will keep re-reconciling every <reconcileInterval> as specified in
	// RegisterMultiClusterReconciler)
	return false, nil
}

// Visible for testing
func NewGlooResourceReconciler(
	ctx context.Context,
	glooInstanceClient fedv1.GlooInstanceClient,
	snapshotBuilder input.Builder,
	translator translator.GlooInstanceTranslator,
) GlooResourceReconciler {
	return &glooResourceReconcilerImpl{
		ctx:                ctx,
		glooInstanceClient: glooInstanceClient,
		fedInputSnapshot:   snapshotBuilder,
		translator:         translator,
		verifier:           buildVerifier(ctx),
	}
}

var (
	// If a Gloo Edge CRD does not exist on a remote Gloo Edge cluster, we don't want to fail the
	// entire reconcile loop. Instead, just log a debug message. Below are all the gvks from the
	// input snapshot for which we can safely ignore errors (excludes built-in k8s resources as well
	// as Settings which are used to construct the GlooInstances)
	gvksToSkipVerify = []schema.GroupVersionKind{
		{
			Group:   "gateway.solo.io",
			Version: "v1",
			Kind:    "Gateway",
		},
		{
			Group:   "gateway.solo.io",
			Version: "v1",
			Kind:    "MatchableHttpGateway",
		},
		{
			Group:   "gateway.solo.io",
			Version: "v1",
			Kind:    "MatchableTcpGateway",
		},
		{
			Group:   "gateway.solo.io",
			Version: "v1",
			Kind:    "VirtualService",
		},
		{
			Group:   "gateway.solo.io",
			Version: "v1",
			Kind:    "RouteTable",
		},
		{
			Group:   "gloo.solo.io",
			Version: "v1",
			Kind:    "Upstream",
		},
		{
			Group:   "gloo.solo.io",
			Version: "v1",
			Kind:    "UpstreamGroup",
		},
		{
			Group:   "gloo.solo.io",
			Version: "v1",
			Kind:    "Proxy",
		},
		{
			Group:   "enterprise.gloo.solo.io",
			Version: "v1",
			Kind:    "AuthConfig",
		},
		{
			Group:   "ratelimit.api.solo.io",
			Version: "v1alpha1",
			Kind:    "RateLimitConfig",
		},
	}
)

func buildVerifier(ctx context.Context) verifier.ServerResourceVerifier {
	verifyOpts := map[schema.GroupVersionKind]verifier.ServerVerifyOption{}
	for _, gvk := range gvksToSkipVerify {
		verifyOpts[gvk] = verifier.ServerVerifyOption_LogDebugIfNotPresent
	}
	return verifier.NewVerifier(ctx, verifyOpts)
}
