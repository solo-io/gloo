package runner

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"

	"go.uber.org/zap"

	"github.com/rotisserie/eris"
	"github.com/solo-io/ext-auth-service/pkg/server"
	"github.com/solo-io/go-utils/contextutils"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
)

func init() {
	_ = view.Register(ocgrpc.DefaultServerViews...)
}

func Run() {
	settings := NewSettings()

	if err := RunWithSettings(context.Background(), settings); err != nil {
		log.Fatalf("server stopped with unexpected error: %s", err.Error())
	}
}

// Blocks until context is cancelled or termination signals are received.
// Normal shutdown will return nil error.
func RunWithSettings(ctx context.Context, settings Settings) error {

	if settings.HeadersToRedact != "" && settings.HeadersToRedact != "-" {
		settings.ExtAuthSettings.HeadersToRedact = strings.Fields(settings.HeadersToRedact)
	}

	listener, err := buildListener(ctx, settings)
	if err != nil {
		return err
	}

	xdsConfigSource := NewConfigSource(settings)

	extAuthServer, err := server.NewServerBuilder().
		WithSettings(settings.ExtAuthSettings).
		WithModule(xdsConfigSource).
		WithListener(listener).
		Build()

	if err != nil {
		return eris.Wrap(err, "failed to build ext auth server instance")
	}

	return extAuthServer.Start(ctx)
}

func buildListener(ctx context.Context, settings Settings) (net.Listener, error) {
	logger := contextutils.LoggerFrom(ctx)

	var (
		addr, runMode, network, tlsMode string
		err                             error
	)

	if settings.ServerUDSAddr != "" {
		addr = settings.ServerUDSAddr
		runMode = "unixDomainSocket"
		network = "unix"
	} else {
		addr = fmt.Sprintf(":%d", settings.ExtAuthSettings.ServerPort)
		runMode = "gRPC"
		network = "tcp"
	}

	var listener net.Listener
	if settings.TlsEnabled {
		logger.Info("TLS is enabled, loading certificates")

		tlsMode = "secure"
		keyPair, err := tls.LoadX509KeyPair(settings.CertPath, settings.KeyPath)
		if err != nil {
			logger.Warnw("failed to load certificates from files, trying to load them from environment variables",
				zap.Error(err), zap.String("CertPath", settings.CertPath), zap.String("KeyPath", settings.KeyPath))
			if len(settings.Cert) == 0 || len(settings.Key) == 0 {
				return nil, eris.New("must provide Cert and Key when running in TLS mode")
			}
			keyPair, err = tls.X509KeyPair(settings.Cert, settings.Key)
			if err != nil {
				return nil, err
			}
		}
		cfg := &tls.Config{Certificates: []tls.Certificate{keyPair}}
		listener, err = tls.Listen(network, addr, cfg)
	} else {
		tlsMode = "insecure"
		listener, err = net.Listen(network, addr)
	}

	if err != nil {
		return nil, eris.Wrap(err, "failed to announce on network")
	}

	logger.Infof("external auth server running in [%s] [%s] mode, listening at [%s]", tlsMode, runMode, addr)

	return listener, nil
}
