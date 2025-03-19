package setuputils

import (
	"context"
	"flag"
	"io"
	"net/http"
	"os"
	"strings"
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
	"crypto/tls"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/pkg/errors"
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
	logger.Infof("Waiting for Kubernetes API server to be healthy...")
	// Wait for Kubernetes API server to be healthy
	if err := waitForKubeApiServer(ctx); err != nil {
		return err
	}

	settingsClient, err := fileOrKubeSettingsClient(ctx, setupNamespace, setupDir)
	if err != nil {
		return err
	}

	if err := settingsClient.Register(); err != nil {
		return err
	}

	identity, err := startLeaderElection(ctx, setupDir, opts.ElectionConfig)
	if err != nil {
		return err
	}

	namespaceClient, err := namespaces.NewKubeNamespaceClient(ctx)
	// If there is any error when creating a KubeNamespaceClient (RBAC issues) default to a fake client
	if err != nil {
		namespaceClient = &namespaces.NoOpKubeNamespaceWatcher{}
	}

	// settings come from the ResourceClient in the settingsClient
	// the eventLoop will Watch the emitter's settingsClient to receive settings from the ResourceClient
	emitter := v1.NewSetupEmitter(settingsClient, namespaceClient)
	settingsRef := &core.ResourceRef{Namespace: setupNamespace, Name: setupName}
	eventLoop := v1.NewSetupEventLoop(emitter, NewSetupSyncer(settingsRef, opts.SetupFunc, identity))
	errs, err := eventLoop.Run([]string{setupNamespace}, clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: time.Second,
	})
	if err != nil {
		return err
	}
	for err := range errs {
		if opts.ExitOnError {
			contextutils.LoggerFrom(ctx).Fatalf("error in setup: %v", err)
		}
		contextutils.LoggerFrom(ctx).Errorf("error in setup: %v", err)
	}
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
	if electionConfig == nil || settingsDir != "" || leaderelector.IsDisabled() {
		// If a component does not contain election config, it does not support HA
		// If the settingsDir is non-empty, it means that Settings are not defined in Kubernetes and therefore we can't use the
		// leader election library which depends on Kubernetes
		// If leader election is explicitly disabled, it means a user has decided not to opt-into HA
		return singlereplica.NewElectionFactory().StartElection(ctx, electionConfig)
	}

	cfg, err := kubeutils.GetRestConfigWithKubeContext("")
	if err != nil {
		return nil, err
	}
	return kube2.NewElectionFactory(cfg).StartElection(ctx, electionConfig)
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

// waitForKubeApiServer polls the Kubernetes API server until it's healthy or the context is canceled
func waitForKubeApiServer(ctx context.Context) error {
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("Waiting for Kubernetes API server to be healthy...")
	
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		return errors.Wrap(err, "building kube config")
	}
	
	// Create a client for the /healthz endpoint
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.TLSClientConfig.Insecure,
			},
		},
	}
	
	// Construct the API server health check URL
	healthzURL := config.Host
	if !strings.HasSuffix(healthzURL, "/") {
		healthzURL += "/"
	}
	healthzURL += "healthz"
	
	// Poll until healthy
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := http.NewRequestWithContext(ctx, "GET", healthzURL, nil)
			if err != nil {
				logger.Warnf("Error creating request to check API server health: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			
			// Add auth if needed
			if config.BearerToken != "" {
				req.Header.Set("Authorization", "Bearer "+config.BearerToken)
			}
			
			resp, err := client.Do(req)
			if err != nil {
				logger.Debugf("API server health check failed: %v, retrying in 5 seconds", err)
				time.Sleep(5 * time.Second)
				continue
			}
			defer resp.Body.Close()
			
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				if string(body) == "ok" {
					logger.Infof("Kubernetes API server is healthy")
					return nil
				}
			}
			
			logger.Debugf("API server returned non-OK status: %d, retrying in 5 seconds", resp.StatusCode)
			time.Sleep(5 * time.Second)
		}
	}
}
