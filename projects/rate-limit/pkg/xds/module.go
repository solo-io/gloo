package xds

import (
	"context"
	"fmt"
	"os"

	"github.com/solo-io/solo-projects/pkg/xds"

	glooRatelimitSyncer "github.com/solo-io/gloo/projects/gloo/pkg/syncer/ratelimit"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	solo_api_rl "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/rate-limiter/pkg/modules"
	ratelimit "github.com/solo-io/rate-limiter/pkg/service"
	"github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"go.uber.org/zap"
)

var (
	rlConnectedStateDescription = "zero indicates ratelimit is unable to connect to connect to gloo to get its configuration"
	mRlConnectedState           = stats.Int64("glooe.ratelimit/xds_client_connected_state", rlConnectedStateDescription, "1")

	rlConnectedStateView = &view.View{
		Name:        "glooe.ratelimit/connected_state",
		Measure:     mRlConnectedState,
		Description: rlConnectedStateDescription,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{},
	}

	rlConnectedStateCounterDescription = "number of client connections ratelimit has made to gloo to get its configuration"
	mRlConnectedStateCounter           = stats.Int64("glooe.ratelimit/xds_client_connect_counter", rlConnectedStateCounterDescription, "1")

	rlConnectedStateCounterView = &view.View{
		Name:        "glooe.ratelimit/connected_state_counter",
		Measure:     mRlConnectedStateCounter,
		Description: rlConnectedStateCounterDescription,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	_ = view.Register(rlConnectedStateView, rlConnectedStateCounterView)
}

const (
	// The rate limit server sends xDS discovery requests to Gloo to get its configuration from Gloo. This constant determines
	// the value of the nodeInfo.Metadata.role field that the server sends along to retrieve its configuration snapshot,
	// similarly to how the regular Gloo gateway-proxies do.
	ServerRole = glooRatelimitSyncer.ServerRole
)

var ConfigGenErr = func(err error, domain string) error {
	return eris.Wrapf(err, "failed to generate configuration for domain [%s]", domain)
}

// This `RunnableModule` implementation uses xDS to get rate limit configuration from the GlooE control plane.
type configSource struct {
	glooAddress     string
	domainGenerator shims.RateLimitDomainGenerator
}

func NewConfigSource(settings Settings, domainGenerator shims.RateLimitDomainGenerator) modules.RunnableModule {
	return &configSource{
		glooAddress:     settings.GlooAddress,
		domainGenerator: domainGenerator,
	}
}

func (*configSource) Name() string {
	return "xDS"
}

func (x *configSource) Run(ctx context.Context, service ratelimit.RateLimitServiceServer) error {
	nodeInfo, err := x.getNodeInfo()

	xDSClientLoopFunc := func(ctx context.Context) error {
		stats.Record(ctx, mRlConnectedStateCounter.M(int64(1)))

		logger := contextutils.LoggerFrom(ctx)

		client := v1.NewRateLimitConfigClient(&nodeInfo, func(version string, resources []*v1.RateLimitConfig) error {
			logger.Debugw("received new rate limit config", zap.Any("config", resources))

			for _, cfg := range resources {
				// Use the domain for configId for xds resources
				domain, err := x.domainGenerator.NewRateLimitDomain(ctx, cfg.Domain, cfg.Domain,
					&solo_api_rl.RateLimitConfigSpec_Raw{Descriptors: cfg.Descriptors, SetDescriptors: cfg.SetDescriptors})
				if err != nil {
					logger.Errorw("failed to generate configuration for domain", zap.String("domain", cfg.Domain), zap.Error(err))
					return ConfigGenErr(err, cfg.Domain)
				}
				service.SetDomain(domain)
			}
			return nil
		})

		conn, err := xds.GetXdsClientConnection(ctx, x.glooAddress)
		if err != nil {
			logger.Errorw("failed to establish connection to gloo xDS server", zap.String("glooAddress", x.glooAddress), zap.Error(err))
			return err
		}

		stats.Record(ctx, mRlConnectedState.M(int64(1)))
		defer func() {
			err = conn.Close()
			if err != nil {
				contextutils.LoggerFrom(ctx).Errorw("failed to close grpc connection", zap.Any("error", err))
			} else {
				contextutils.LoggerFrom(ctx).Infow("closed grpc connection", zap.Any("address", x.glooAddress))
			}
			stats.Record(ctx, mRlConnectedState.M(int64(0)))
		}()

		return client.Start(ctx, conn)
	}

	err = contextutils.NewExponentioalBackoff(contextutils.ExponentioalBackoff{}).Backoff(ctx, xDSClientLoopFunc)
	if err == context.Canceled {
		return nil
	}
	return err
}

func (x *configSource) getNodeInfo() (core.Node, error) {
	var nodeInfo core.Node
	var err error

	nodeInfo.Id, err = os.Hostname()
	// TODO(yuval-k): unhardcode this
	if err != nil {
		nodeInfo.Id = fmt.Sprintf("%s-unknown", ServerRole)
	}
	nodeInfo.Cluster = ServerRole
	nodeInfo.Metadata = &_struct.Struct{
		Fields: map[string]*_struct.Value{
			"role": {
				Kind: &_struct.Value_StringValue{
					StringValue: ServerRole,
				},
			},
		},
	}
	return nodeInfo, err
}
