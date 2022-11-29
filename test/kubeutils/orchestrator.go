package kubeutils

import (
	"context"
	"fmt"

	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/solo-projects/test/services"
)

var _ Orchestrator = new(KindOrchestrator)

type Orchestrator interface {
	SetClusterContext(ctx context.Context, clusterName string) error
}

type KindOrchestrator struct {
}

func NewKindOrchestrator() *KindOrchestrator {
	return &KindOrchestrator{}
}

func (k *KindOrchestrator) SetClusterContext(ctx context.Context, clusterName string) error {
	contextutils.LoggerFrom(ctx).Debugf("Setting kube context to %s-%s", "kind", clusterName)
	return services.Kubectl("config", "use-context", fmt.Sprintf("kind-%s", clusterName))
}

type NoOpOrchestrator struct {
}

func NewNoOpOrchestrator() *NoOpOrchestrator {
	return &NoOpOrchestrator{}
}

func (n *NoOpOrchestrator) SetClusterContext(_ context.Context, _ string) error {
	return nil
}
