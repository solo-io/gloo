package webhook

import (
	"context"

	"github.com/go-logr/zapr"
	"github.com/solo-io/go-utils/contextutils"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	rbacconfig "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/config"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/rbac"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/webhook"
	test_v1alpha1 "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/test/internal/api/test.multicluster.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/test/internal/pkg/placement"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaputil "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func Start(rootCtx context.Context, rbacCfg *rbacconfig.Config) error {
	atomicLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	baseLogger := zaputil.NewRaw(
		zaputil.Level(&atomicLevel),
	)
	// klog
	zap.ReplaceGlobals(baseLogger)
	// controller-runtime
	log.SetLogger(zapr.NewLogger(baseLogger))
	// go-utils
	contextutils.SetFallbackLogger(baseLogger.Sugar())
	ctx := CreateRootContext(rootCtx, "multicluster-admission-webhook")
	logger := contextutils.LoggerFrom(ctx)

	mgr, err := initMasterManager(rbacCfg)
	if err != nil {
		return err
	}

	if err := webhook.InitializeWebhook(ctx, mgr, rbacCfg, initPlacementParser(mgr)); err != nil {
		return err
	}

	if err := mgr.Start(ctx); err != nil {
		logger.Fatalf("error while running multicluster-admission-webhook: %+v", err)
	}
	return nil
}

func initPlacementParser(mgr manager.Manager) rbac.Parser {
	return placement.NewParser(mgr.GetScheme(), placement.NewTypedParser())
}

// get the manager for the local cluster; we will use this as our "master" cluster
func initMasterManager(rbacCfg *rbacconfig.Config) (manager.Manager, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		return nil, err
	}
	// Add MultiClusterRole/RoleBinding scheme
	if err := multicluster_v1alpha1.SchemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
		return nil, err
	}
	// Add Test scheme
	if err := test_v1alpha1.SchemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
		return nil, err
	}
	return mgr, nil
}

func CreateRootContext(customCtx context.Context, name string) context.Context {
	rootCtx := customCtx
	if rootCtx == nil {
		rootCtx = context.Background()
	}
	setupLogging(name)
	return rootCtx
}

func setupLogging(name string) {
	// Default to info level logging
	level := zapcore.DebugLevel
	atomicLevel := zap.NewAtomicLevelAt(level)
	baseLogger := zaputil.NewRaw(zaputil.Level(&atomicLevel)).Named(name)
	// klog
	zap.ReplaceGlobals(baseLogger)
	// controller-runtime
	log.SetLogger(zapr.NewLogger(baseLogger))
	// go-utils
	contextutils.SetFallbackLogger(baseLogger.Sugar())
}
