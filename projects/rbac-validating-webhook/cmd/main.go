package main

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	fed_bootstrap "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/bootstrap"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/bootstrap"
	rbacconfig "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/config"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/webhook"
	"github.com/solo-io/solo-projects/projects/rbac-validating-webhook/pkg/placement"
	"go.uber.org/zap"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	rootCtx := bootstrap.CreateRootContext(context.Background(), "gloo-fed-rbac-validation")
	mgr := fed_bootstrap.MustLocalManager(rootCtx)
	rbacCfg := rbacconfig.NewConfig()

	if err := webhook.InitializeWebhook(
		rootCtx,
		mgr,
		rbacCfg,
		placement.NewParser(mgr.GetScheme(), placement.NewTypedParser()),
	); err != nil {
		contextutils.LoggerFrom(rootCtx).Fatalw("Unable to start admission webhook", zap.Error(err))
	}

	err := mgr.Start(rootCtx)
	if err != nil {
		contextutils.LoggerFrom(rootCtx).Fatalw("An error occurred", zap.Error(err))
	}
	contextutils.LoggerFrom(rootCtx).Infow("Shutting down, root context cancelled.")
}
