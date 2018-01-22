package pkglog // import "go.pedge.io/pkg/log"

import (
	"fmt"
	"log/syslog"
	"os"
	"path/filepath"
	"strings"

	"go.pedge.io/lion"
	"go.pedge.io/lion/syslog"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Env defines a struct for environment variables that can be parsed with go.pedge.io/env.
type Env struct {
	// DisableStderrLog says to disable logging to stderr.
	DisableStderrLog bool `env:"DISABLE_STDERR_LOG"`
	// The directory to write rotating logs to.
	// If not set and SyslogNetwork and SyslogAddress not set, logs will be to stderr.
	LogDir string `env:"LOG_DIR"`
	// The syslog network, either udp or tcp.
	// Must be set with SyslogAddress.
	// If not set and LogDir not set, logs will be to stderr.
	SyslogNetwork string `env:"SYSLOG_NETWORK"`
	// The syslog host:port.
	// Must be set with SyslogNetwork.
	// If not set and LogDir not set, logs will be to stderr.
	SyslogAddress string `env:"SYSLOG_ADDRESS"`
	// The level to log at, must be one of DEBUG, INFO, WARN, ERROR, FATAL, PANIC.
	LogLevel string `env:"LOG_LEVEL"`
}

// SetupLogging sets up logging.
func SetupLogging(appName string, env Env) error {
	var pushers []lion.Pusher
	if !env.DisableStderrLog {
		pushers = append(
			pushers,
			lion.NewTextWritePusher(
				os.Stderr,
			),
		)
	}
	if env.LogDir != "" {
		pushers = append(
			pushers,
			lion.NewTextWritePusher(
				&lumberjack.Logger{
					Filename:   filepath.Join(env.LogDir, fmt.Sprintf("%s.log", appName)),
					MaxBackups: 3,
				},
			),
		)
	}
	if env.SyslogNetwork != "" && env.SyslogAddress != "" {
		writer, err := syslog.Dial(
			env.SyslogNetwork,
			env.SyslogAddress,
			syslog.LOG_INFO,
			appName,
		)
		if err != nil {
			return err
		}
		pushers = append(
			pushers,
			sysloglion.NewPusher(
				writer,
			),
		)
	}
	if len(pushers) > 0 {
		lion.SetLogger(
			lion.NewLogger(
				lion.NewMultiPusher(
					pushers...,
				),
			),
		)
	} else {
		lion.SetLogger(
			lion.DiscardLogger,
		)
	}
	lion.RedirectStdLogger()
	if env.LogLevel != "" {
		level, err := lion.NameToLevel(strings.ToUpper(env.LogLevel))
		if err != nil {
			return err
		}
		lion.SetLevel(level)
	}
	return nil
}
