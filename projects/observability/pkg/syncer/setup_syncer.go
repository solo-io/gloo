package syncer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/avast/retry-go/v4"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	"golang.org/x/exp/slices"

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
	tmplExtension = ".json.tmpl"

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

func Setup(ctx context.Context, kubeCache kube.SharedCache, inMemoryCache memory.InMemoryResourceCache, settings *gloov1.Settings, _ leaderelector.Identity) error {
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

	// Set the prefix to max 20 characters so the names of the dashboards can be unique
	dashboardPrefix, extraDashboardTags, err := generateDashboardPrefixAndTags(settings.GetObservabilityOptions().GetGrafanaIntegration().GetDashboardPrefix())
	if err != nil {
		return err
	}

	opts := Opts{
		WriteNamespace:  writeNamespace,
		WatchNamespaces: watchNamespaces,
		Upstreams:       upstreamFactory,
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: refreshRate,
		},
		DevMode:                    true,
		DefaultDashboardFolderId:   defaultDashboardFolderId,
		DashboardPrefix:            dashboardPrefix,
		ExtraDashboardTags:         extraDashboardTags,
		ExtraMetricQueryParameters: settings.GetObservabilityOptions().GetGrafanaIntegration().GetExtraMetricQueryParameters(),
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

	files, err := os.ReadDir(defaultDashboardDir)
	if err != nil {
		return err
	}

	// The observability deployment may finish and this Sync may start running before the Grafana deployment is ready
	// Retry the Sync until Grafana is ready
	// https://github.com/solo-io/solo-projects/issues/307
	var folders []grafana.FolderProperties
	err = retry.Do(func() error {
		folders, err = dashboardClient.GetAllFolderIds()
		if err != nil && strings.Contains(err.Error(), "connection refused") {
			return errors.Errorf("Grafana is not yet up : %+v", err)
		}
		return nil
	}, grafanaSyncRetryOpts...)
	if err != nil {
		return err
	}

	// By default, the static folders are created in the `gloo` folder. Keep this behavour.
	// Ref: https://github.com/solo-io/solo-projects/blame/42e057e6039e65f79a61f11c56bc3a046fd5b6df/install/helm/gloo-ee/values-template.yaml#L460
	folderID := generalFolderId
	if opts.DefaultDashboardFolderId == generalFolderId {
		for _, folder := range folders {
			if folder.Title == "gloo" {
				folderID = folder.ID
				break
			}
		}
	} else {
		for _, folder := range folders {
			if folder.ID == opts.DefaultDashboardFolderId {
				folderID = folder.ID
				break
			}
		}
	}

	defaultDashboardUids := make(map[string]struct{})
	templateGeneratorOpts := []template.Option{
		template.WithDashboardPrefix(opts.DashboardPrefix),
		template.WithExtraDashboardTags(opts.ExtraDashboardTags),
		template.WithExtraMetricQueryParameters(opts.ExtraMetricQueryParameters)}

	for _, file := range files {
		filename := file.Name()
		if !strings.HasSuffix(filename, tmplExtension) {
			continue
		}
		defaultJsonStr, err := getDashboardJson(opts.WatchOpts.Ctx, filepath.Join(defaultDashboardDir, filename))
		if err != nil {
			return err
		}

		uid := template.ToUID(opts.DashboardPrefix) + strings.TrimSuffix(filename, tmplExtension)
		templateGenerator := template.NewDefaultDashboardTemplateGenerator(uid, defaultJsonStr, templateGeneratorOpts...)
		loadDefaultDashboard(opts.WatchOpts.Ctx, templateGenerator, folderID, dashboardClient)
		defaultDashboardUids[uid] = struct{}{}
	}

	dashboardJsonTemplate, err := getDashboardJson(opts.WatchOpts.Ctx, upstreamDashboardTemplatePath)
	if err != nil {
		return err
	}

	dashSyncer := NewGrafanaDashboardSyncer(dashboardClient, snapshotClient, dashboardJsonTemplate, defaultDashboardUids,
		WithDefaultDashboardFolderId(opts.DefaultDashboardFolderId),
		WithDashboardPrefix(opts.DashboardPrefix),
		WithExtraDashboardTags(opts.ExtraDashboardTags),
		WithExtraMetricQueryParameters(opts.ExtraMetricQueryParameters))

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
			caCert, err := os.ReadFile(caCrtEnvValue)
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
	bytes, err := os.ReadFile(filename)
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnf("Error reading file %s - %s", filename, err.Error())
		return "", nil
	}

	return string(bytes), nil
}

// generateDashboardPrefixAndTags converts the dashboard prefix into a grafana compatible UID.
// It also adds the prefix to the list of additional tags to be used to tag the dashboard
func generateDashboardPrefixAndTags(dashboardPrefix string) (string, []string, error) {
	extraTags := append([]string{}, defaultTags...)
	if dashboardPrefix != "" {
		if len(dashboardPrefix) > MaxPrefixLength {
			return "", nil, fmt.Errorf("dashboard prefix [%s] exceeds the maximum allowed length of 20 characters", dashboardPrefix)
		}
		dashboardPrefix = template.ToUID(dashboardPrefix)
		// Since a dashboard prefix is provided, add it to the list of tags. We do this to ensure we only manage dashboards with the prefix. Now when querying resources, we only do so for resources with matching tags
		if !slices.Contains(extraTags, dashboardPrefix) {
			extraTags = append(extraTags, dashboardPrefix)
		}
	} else {
		// We still tag dashboards that do not have a prefix to ensure we do not fetch dashboards with a prefix
		// Eg: If there are two edge installations that manage dashboards on an external grafana instance, one with no prefix and another with a prefix 'prefix',
		// when the no prefix installation would fetch dashboards with the default tags ['gloo'], it would still return dashboards with tags ['gloo', 'prefix']
		// To avoid this we add the defaultPrefixTag to the list of defaultTags. This is a workaround until we find a way to uniquely identify an edge installation across multiple clusters (with the same release name and namespace)
		if !slices.Contains(extraTags, defaultPrefixTag) {
			extraTags = append(extraTags, defaultPrefixTag)
		}
	}
	return dashboardPrefix, extraTags, nil
}
