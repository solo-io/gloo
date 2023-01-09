package runner

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/solo-io/solo-projects/pkg/version"

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

	loggingContext := []interface{}{"version", version.Version}
	ctx := contextutils.WithLoggerValues(context.Background(), loggingContext...)
	if err := RunWithSettings(ctx, settings); err != nil {
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
		if err := cleanUnixSocket(ctx, addr); err != nil {
			return nil, err
		}
	} else {
		addr = fmt.Sprintf(":%d", settings.ExtAuthSettings.ServerPort)
		runMode = "gRPC"
		network = "tcp"
	}

	var listener net.Listener
	if settings.TlsEnabled {
		logger.Info("TLS is enabled, loading certificates")

		tlsMode = "secure"
		keyPair, keyErr := tls.LoadX509KeyPair(settings.CertPath, settings.KeyPath)
		if keyErr != nil {
			logger.Warnw("failed to load certificates from files, trying to load them from environment variables",
				zap.Error(keyErr), zap.String("CertPath", settings.CertPath), zap.String("KeyPath", settings.KeyPath))
			if len(settings.Cert) == 0 || len(settings.Key) == 0 {
				return nil, eris.New("must provide Cert and Key when running in TLS mode")
			}
			keyPair, keyErr = tls.X509KeyPair(settings.Cert, settings.Key)
			if keyErr != nil {
				return nil, keyErr
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

// cleanUnixSocket will remove the socket file if it exists
func cleanUnixSocket(ctx context.Context, addr string) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		select {
		case <-sigs:
		case <-ctx.Done():
		}
		err := os.RemoveAll(addr)
		if err != nil {
			contextutils.LoggerFrom(ctx).Warnf("server failed to clean unix socket: %s", err.Error())
		}
	}()
	// Remove the unix socket file, because it could already be in use.
	// this is safe because this is the only process that should be using this file as a server
	// clients are still able to connect
	return os.RemoveAll(addr)
}
