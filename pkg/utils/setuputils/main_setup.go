package setuputils

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/namespaces"

	"github.com/go-logr/zapr"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	kube2 "github.com/solo-io/gloo/pkg/bootstrap/leaderelector/kube"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"
	"github.com/solo-io/gloo/pkg/version"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaputil "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type SetupOpts struct {
	LoggerName string
	// logged as the version of Gloo currently executing
	Version     string
	SetupFunc   SetupFunc
	ExitOnError bool
	CustomCtx   context.Context

	// optional - if present, add these values in each JSON log line in the gloo pod.
	// By default, we already log the gloo version.
	LoggingPrefixVals []interface{}
	// optional - if present, report usage with the payload this discovers
	// should really only provide it in very intentional places- in the gloo pod, and in glooctl
	// otherwise, we'll provide redundant copies of the usage data

	ElectionConfig *leaderelector.ElectionConfig
}

var once sync.Once

// Main is the main entrypoint for running Gloo Edge components
// It works by performing the following:
//  1. Initialize a SettingsClient backed either by Kubernetes or a File
//  2. Run an event loop, watching events on the Settings resource, and executing the
//     opts.SetupFunc whenever settings change
//
// This allows Gloo components to automatically receive updates to Settings and reload their
// configuration, without needing to restart the container
func Main(opts SetupOpts) error {
	// prevent panic if multiple flag.Parse called concurrently
	once.Do(func() {
		flag.Parse()
	})

	ctx := opts.CustomCtx
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = contextutils.WithLogger(ctx, opts.LoggerName)
	loggingContext := append([]interface{}{"version", opts.Version}, opts.LoggingPrefixVals...)
	ctx = contextutils.WithLoggerValues(ctx, loggingContext...)

	logger := contextutils.LoggerFrom(ctx)
	logger.Debugw("Starting main setup function",
		"issue", "8539",
		"loggerName", opts.LoggerName,
		"version", opts.Version,
		"exitOnError", opts.ExitOnError,
		"hasElectionConfig", opts.ElectionConfig != nil)

	logger.Debug("Initializing settings client", "issue", "8539")
	settingsClient, err := fileOrKubeSettingsClient(ctx, setupNamespace, setupDir)
	if err != nil {
		logger.Errorw("Failed to create settings client", "issue", "8539", "error", err)
		return err
	}
	logger.Debug("Settings client created successfully", "issue", "8539")

	if err := settingsClient.Register(); err != nil {
		logger.Errorw("Failed to register settings client", "issue", "8539", "error", err)
		return err
	}
	logger.Debug("Settings client registered successfully", "issue", "8539")

	logger.Debug("Starting leader election process", "issue", "8539")
	identity, err := startLeaderElection(ctx, setupDir, opts.ElectionConfig)
	if err != nil {
		logger.Errorw("Leader election failed to start", "issue", "8539", "error", err)
		return err
	}
	logger.Debugw("Leader election started",
		"issue", "8539",
		"isLeader", identity.IsLeader(),
		"hasElectionConfig", opts.ElectionConfig != nil,
		"settingsDir", setupDir)

	// Wait a moment to see if we become leader and log the result
	go func() {
		time.Sleep(100 * time.Millisecond)
		logger.Debugw("Initial leadership status", "issue", "8539", "isLeader", identity.IsLeader())

		// Monitor for leadership changes
		go func() {
			<-identity.Elected()
			logger.Infow("Component elected as leader", "issue", "8539")
		}()
	}()

	logger.Debug("Creating namespace client", "issue", "8539")
	namespaceClient, err := namespaces.NewKubeNamespaceClient(ctx)
	// If there is any error when creating a KubeNamespaceClient (RBAC issues) default to a fake client
	if err != nil {
		logger.Warnw("Failed to create KubeNamespaceClient, using NoOp client", "issue", "8539", "error", err)
		namespaceClient = &namespaces.NoOpKubeNamespaceWatcher{}
	} else {
		logger.Debug("Namespace client created successfully", "issue", "8539")
	}

	logger.Debug("Setting up event loop", "issue", "8539")
	// settings come from the ResourceClient in the settingsClient
	// the eventLoop will Watch the emitter's settingsClient to receive settings from the ResourceClient
	emitter := v1.NewSetupEmitter(settingsClient, namespaceClient)
	settingsRef := &core.ResourceRef{Namespace: setupNamespace, Name: setupName}
	logger.Debugw("Created settings reference",
		"issue", "8539",
		"namespace", setupNamespace,
		"name", setupName)

	eventLoop := v1.NewSetupEventLoop(emitter, NewSetupSyncer(settingsRef, opts.SetupFunc, identity))
	logger.Debugw("Starting event loop",
		"issue", "8539",
		"watchNamespaces", []string{setupNamespace},
		"refreshRate", time.Second)

	errs, err := eventLoop.Run([]string{setupNamespace}, clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: time.Second,
	})
	if err != nil {
		logger.Errorw("Failed to start event loop", "issue", "8539", "error", err)
		return err
	}
	logger.Debug("Event loop started successfully", "issue", "8539")

	for err := range errs {
		if opts.ExitOnError {
			logger.Fatalw("Fatal error in setup", "issue", "8539", "error", err, "exitOnError", true)
		}
		logger.Errorw("Error in setup", "issue", "8539", "error", err)
	}
	logger.Debug("Main setup function completed", "issue", "8539")
	return nil
}

func fileOrKubeSettingsClient(ctx context.Context, setupNamespace, settingsDir string) (v1.SettingsClient, error) {
	if settingsDir != "" {
		contextutils.LoggerFrom(ctx).Infow("using filesystem for settings", zap.String("directory", settingsDir))
		return v1.NewSettingsClient(ctx, &factory.FileResourceClientFactory{
			RootDir: settingsDir,
		})
	}

	cfg, err := kubeutils.GetRestConfigWithKubeContext("")
	if err != nil {
		return nil, err
	}
	return v1.NewSettingsClient(ctx, &factory.KubeResourceClientFactory{
		Crd:                v1.SettingsCrd,
		Cfg:                cfg,
		SharedCache:        kube.NewKubeCache(ctx),
		NamespaceWhitelist: []string{setupNamespace},
	})
}

func startLeaderElection(ctx context.Context, settingsDir string, electionConfig *leaderelector.ElectionConfig) (leaderelector.Identity, error) {
	logger := contextutils.LoggerFrom(ctx)

	logger.Debugw("Evaluating leader election requirements",
		"issue", "8539",
		"hasElectionConfig", electionConfig != nil,
		"settingsDir", settingsDir,
		"leaderElectionDisabled", leaderelector.IsDisabled())

	if electionConfig == nil || settingsDir != "" || leaderelector.IsDisabled() {
		// If a component does not contain election config, it does not support HA
		// If the settingsDir is non-empty, it means that Settings are not defined in Kubernetes and therefore we can't use the
		// leader election library which depends on Kubernetes
		// If leader election is explicitly disabled, it means a user has decided not to opt-into HA

		var reason string
		if electionConfig == nil {
			reason = "no election config provided - component does not support HA"
		} else if settingsDir != "" {
			reason = "using file-based settings - kubernetes leader election not available"
		} else if leaderelector.IsDisabled() {
			reason = "leader election explicitly disabled"
		}

		logger.Infow("Using single replica election (no HA)", "issue", "8539", "reason", reason)
		return singlereplica.NewElectionFactory().StartElection(ctx, electionConfig)
	}

	logger.Debugw("Using Kubernetes-based leader election",
		"issue", "8539",
		"electionId", electionConfig.Id,
		"namespace", electionConfig.Namespace)

	cfg, err := kubeutils.GetRestConfigWithKubeContext("")
	if err != nil {
		logger.Errorw("Failed to get kubernetes config for leader election", "issue", "8539", "error", err)
		return nil, err
	}

	logger.Debug("Kubernetes config obtained, starting election factory", "issue", "8539")
	identity, err := kube2.NewElectionFactory(cfg).StartElection(ctx, electionConfig)
	if err != nil {
		logger.Errorw("Failed to start kubernetes leader election", "issue", "8539", "error", err)
		return nil, err
	}

	logger.Debug("Kubernetes leader election factory started successfully", "issue", "8539")
	return identity, nil
}

// SetupLogging sets up controller-runtime logging
func SetupLogging(ctx context.Context, loggerName string) {
	level := zapcore.InfoLevel
	// if log level is set in env, use that
	if envLogLevel := os.Getenv(contextutils.LogLevelEnvName); envLogLevel != "" {
		if err := (&level).Set(envLogLevel); err != nil {
			contextutils.LoggerFrom(ctx).Infof("Could not set log level from env %s=%s, available levels "+
				"can be found here: https://pkg.go.dev/go.uber.org/zap/zapcore?tab=doc#Level",
				contextutils.LogLevelEnvName,
				envLogLevel,
				zap.Error(err),
			)
		}
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	baseLogger := zaputil.NewRaw(
		zaputil.Level(&atomicLevel),
		zaputil.RawZapOpts(zap.Fields(zap.String("version", version.Version))),
	).Named(loggerName)

	// controller-runtime
	log.SetLogger(zapr.NewLogger(baseLogger))
}
