package syncer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana/template"

	"github.com/golang/protobuf/ptypes"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils"
	"github.com/solo-io/gloo/pkg/utils/setuputils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-projects/pkg/version"
	v1 "github.com/solo-io/solo-projects/projects/observability/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"
	"k8s.io/client-go/rest"
)

const (
	jsonExtention = ".json"

	observability                 = "observability"
	grafanaUrl                    = "GRAFANA_URL"
	grafanaUsername               = "GRAFANA_USERNAME"
	grafanaPassword               = "GRAFANA_PASSWORD"
	grafanaApiKey                 = "GRAFANA_API_KEY"
	grafanaCaCrt                  = "GRAFANA_CA_BUNDLE"
	upstreamDashboardTemplatePath = "/observability/dashboard-template.json"
	defaultDashboardDir           = "/observability/defaults"
)

func Main() error {
	return setuputils.Main(setuputils.SetupOpts{
		LoggerName:  observability,
		Version:     version.Version,
		ExitOnError: true,
		SetupFunc:   Setup,
	})
}

func Setup(ctx context.Context, kubeCache kube.SharedCache, inMemoryCache memory.InMemoryResourceCache, settings *gloov1.Settings) error {
	var (
		cfg *rest.Config
	)
	params := bootstrap.NewConfigFactoryParams(settings, inMemoryCache, kubeCache, &cfg, nil)
	upstreamFactory, err := bootstrap.ConfigFactoryForSettings(params, gloov1.UpstreamCrd)
	if err != nil {
		return err
	}

	refreshRate, err := ptypes.Duration(settings.GetRefreshRate())
	if err != nil {
		return err
	}

	writeNamespace := settings.DiscoveryNamespace
	if writeNamespace == "" {
		writeNamespace = gloodefaults.GlooSystem
	}
	watchNamespaces := utils.ProcessWatchNamespaces(settings.WatchNamespaces, writeNamespace)

	defaultDashboardFolderId := generalFolderId
	// check if the user inputted a default dashboard id.
	if obsOpts := settings.GetObservabilityOptions(); obsOpts != nil {
		if grafInt := obsOpts.GetGrafanaIntegration(); grafInt != nil {
			if rawFolderId := grafInt.DefaultDashboardFolderId; rawFolderId != nil {
				logger := contextutils.LoggerFrom(ctx)
				logger.Infof("Using inputted default folder id: %d", rawFolderId.GetValue())
				// just accept it if it exists. Validation happens later.
				defaultDashboardFolderId = uint(rawFolderId.GetValue())
			}
		}
	}

	opts := Opts{
		WriteNamespace:  writeNamespace,
		WatchNamespaces: watchNamespaces,
		Upstreams:       upstreamFactory,
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: refreshRate,
		},
		DevMode:                  true,
		DefaultDashboardFolderId: defaultDashboardFolderId,
	}

	return RunObservability(opts)
}

func RunObservability(opts Opts) error {
	opts.WatchOpts = opts.WatchOpts.WithDefaults()
	opts.WatchOpts.Ctx = contextutils.WithLogger(opts.WatchOpts.Ctx, "observability")

	upstreamClient, err := gloov1.NewUpstreamClient(opts.WatchOpts.Ctx, opts.Upstreams)
	if err != nil {
		return err
	}
	if err := upstreamClient.Register(); err != nil {
		return err
	}

	grafanaApiUrl, err := getGrafanaApiUrl()
	if err != nil {
		return err
	}

	restClient, err := buildRestClient(opts.WatchOpts.Ctx, grafanaApiUrl)
	if err != nil {
		return err
	}

	dashboardClient := grafana.NewDashboardClient(restClient)
	snapshotClient := grafana.NewSnapshotClient(restClient)

	files, err := ioutil.ReadDir(defaultDashboardDir)
	if err != nil {
		return err
	}

	defaultDashboardUids := make(map[string]struct{})
	for _, file := range files {
		filename := file.Name()
		if !strings.HasSuffix(filename, jsonExtention) {
			continue
		}
		defaultJsonStr, err := getDashboardJson(opts.WatchOpts.Ctx, filepath.Join(defaultDashboardDir, filename))
		if err != nil {
			return err
		}

		uid := strings.TrimSuffix(filename, jsonExtention)
		templateGenerator := template.NewDefaultJsonGenerator(uid, defaultJsonStr)

		loadDefaultDashboard(opts.WatchOpts.Ctx, templateGenerator, opts.DefaultDashboardFolderId, dashboardClient)
		defaultDashboardUids[uid] = struct{}{}
	}

	dashboardJsonTemplate, err := getDashboardJson(opts.WatchOpts.Ctx, upstreamDashboardTemplatePath)
	if err != nil {
		return err
	}

	dashSyncer := NewGrafanaDashboardSyncer(dashboardClient, snapshotClient, dashboardJsonTemplate, opts.DefaultDashboardFolderId, defaultDashboardUids)

	emitter := v1.NewDashboardsEmitter(upstreamClient)
	eventLoop := v1.NewDashboardsEventLoop(emitter, dashSyncer)
	writeErrs := make(chan error)
	eventLoopErrs, err := eventLoop.Run(opts.WatchNamespaces, opts.WatchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(opts.WatchOpts.Ctx, writeErrs, eventLoopErrs, "event_loop")

	go func() {
		for {
			select {
			case err := <-writeErrs:
				contextutils.LoggerFrom(opts.WatchOpts.Ctx).Errorf("error: %v", err)
			case <-opts.WatchOpts.Ctx.Done():
				close(writeErrs)
				return
			}
		}
	}()
	return nil
}

func buildRestClient(ctx context.Context, grafanaApiUrl string) (grafana.RestClient, error) {
	var (
		caCrtEnvValue = os.Getenv(grafanaCaCrt)
		logger        = contextutils.LoggerFrom(ctx)
	)

	creds, err := buildRestCredentials(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("Setting up HTTPS connection to grafana")

	if strings.HasPrefix(grafanaApiUrl, "https") {
		httpClient := &http.Client{}

		if caCrtEnvValue != "" {
			pool := x509.NewCertPool()
			caCert, err := ioutil.ReadFile(caCrtEnvValue)
			if err != nil {
				caCert = []byte(caCrtEnvValue)
				caCrtEnvValue = "customGrafana.caBundle"
			}
			if ok := pool.AppendCertsFromPEM(caCert); !ok {
				return nil, errors.Errorf("Unable to parse PEM encoded certificate in %s", caCrtEnvValue)
			}
			httpClient = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: pool,
					},
				},
			}
		}
		return grafana.NewRestClient(grafanaApiUrl, httpClient, creds), nil
	}

	logger.Info("Setting up HTTP connection to grafana")
	httpClient := http.DefaultClient
	return grafana.NewRestClient(grafanaApiUrl, httpClient, creds), nil
}

func buildRestCredentials(ctx context.Context) (*grafana.Credentials, error) {
	var (
		username = os.Getenv(grafanaUsername)
		password = os.Getenv(grafanaPassword)
		apiKey   = os.Getenv(grafanaApiKey)
		logger   = contextutils.LoggerFrom(ctx)
	)

	if apiKey != "" {
		logger.Info("Using api key for authentication to grafana")
		return &grafana.Credentials{
			BasicAuth: nil,
			ApiKey:    apiKey,
		}, nil
	} else if username != "" && password != "" {
		logger.Info("Using basic auth for authentication to grafana")
		return &grafana.Credentials{
			BasicAuth: &grafana.BasicAuth{
				Username: username,
				Password: password,
			},
			ApiKey: "",
		}, nil
	} else {
		return nil, grafana.IncompleteGrafanaCredentials
	}
}

func getGrafanaApiUrl() (string, error) {
	url := os.Getenv(grafanaUrl)
	if url == "" {
		return "", NoGrafanaUrl(grafanaUrl)
	}

	return url, nil
}

func getDashboardJson(ctx context.Context, filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnf("Error reading file %s - %s", filename, err.Error())
		return "", nil
	}

	return string(bytes), nil
}
