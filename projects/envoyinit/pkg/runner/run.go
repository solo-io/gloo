package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/solo-io/gloo/pkg/utils/cmdutils"
	"github.com/solo-io/gloo/projects/envoyinit/pkg/downward"
	"github.com/solo-io/go-utils/contextutils"
)

// ErrValidationInterrupted means the envoy validation process did not complete: it was killed by
// context cancellation or an outside signal, such as SIGKILL from the kernel OOM killer, or it crashed
// (ErrValidationCrashed). It says nothing about the validity of the config.
var ErrValidationInterrupted = errors.New("envoy validation process was interrupted before it completed")

// ErrValidationCrashed means envoy crashed while validating the config. A crash can come from the config
// or from the pod's environment, so it wraps ErrValidationInterrupted: callers keep the config Envoy
// runs and retry. Callers that cache results can keep a crash for a short time, so the retries do not
// start envoy each time.
var ErrValidationCrashed = fmt.Errorf("envoy crashed while validating this configuration: %w", ErrValidationInterrupted)

// crashRetried tells an operator what happens after a crash.
const crashRetried = "Envoy keeps the configuration it was last sent for this proxy, and gloo retries."

const (
	// Environment variable for the file that is used to inject input configuration used to bootstrap envoy
	inputConfigPathEnv     = "INPUT_CONF"
	defaultInputConfigPath = "/etc/envoy/envoy.yaml"

	// Environment variable for the file that is written to with transformed bootstrap configuration
	outputConfigPathEnv     = "OUTPUT_CONF"
	defaultOutputConfigPath = "/tmp/envoy.yaml"

	// Environment variable for the path to the envoy executable
	envoyExecutableEnv     = "ENVOY"
	defaultEnvoyExecutable = "/usr/local/bin/envoy"
)

func RunEnvoyValidate(ctx context.Context, envoyExecutable, bootstrapConfig string) error {
	logger := contextutils.LoggerFrom(ctx)

	validateCmd := cmdutils.Command(ctx, envoyExecutable, "--mode", "validate", "--config-path", "/dev/fd/0",
		"-l", "critical", "--log-format", "%v")
	validateCmd = validateCmd.WithStdin(bytes.NewBufferString(bootstrapConfig))

	start := time.Now()
	err := validateCmd.Run()
	logger.Debugf("envoy validation of %d size completed in %s", len(bootstrapConfig), time.Since(start))

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// log a warning and return nil; will allow users to continue to run Gloo locally without
			// relying on the Gloo container with Envoy already published to the expected directory
			logger.Warnf("Unable to validate envoy configuration using envoy at %v; "+
				"skipping additional validation of Gloo config.", envoyExecutable)
			return nil
		}

		// A kill from context cancellation or an outside signal (exit code -1) means validation did not complete.
		// Consumers flatten these errors into report strings, so include the cause in the message.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return fmt.Errorf("%w (ctxErr=%v): %v", ErrValidationInterrupted, ctxErr, err)
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == -1 {
			// envoy crashed. Its output is not logged because it can include secrets from the config.
			if sig, ok := crashSignal(exitErr); ok {
				logger.Errorw(fmt.Sprintf("%v. %s", ErrValidationCrashed, crashRetried),
					"signal", sig.String(), "error", err.Error())
				return fmt.Errorf("%w (%v). %s", ErrValidationCrashed, sig, crashRetried)
			}
			return fmt.Errorf("%w: %v", ErrValidationInterrupted, err)
		}
		return fmt.Errorf("envoy validation mode output: %v, error: %v", err.OutputString(), err.Error())
	}

	return nil
}

// RunEnvoy run Envoy with bootstrap configuration injected from a file
func RunEnvoy(envoyExecutable, inputPath, outputPath string) {
	// 1. Transform the configuration using the Kubernetes Downward API
	bootstrapConfig, err := getAndTransformConfig(inputPath)
	if err != nil {
		log.Fatalf("initializer failed: %v", err)
	}

	// 2. Write to a file for debug purposes
	// since this operation is meant only for debug purposes, we ignore the error
	// this might fail if root fs is read only
	_ = os.WriteFile(outputPath, []byte(bootstrapConfig), 0444)

	// 3. Execute Envoy with the provided configuration
	args := []string{envoyExecutable, "--config-yaml", bootstrapConfig}
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}
	if err = syscall.Exec(args[0], args, os.Environ()); err != nil {
		panic(err)
	}
}

// GetInputConfigPath returns the path to a file containing the Envoy bootstrap configuration
// This configuration may leverage the Kubernetes Downward API
// https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#the-downward-api
func GetInputConfigPath() string {
	return getEnvOrDefault(inputConfigPathEnv, defaultInputConfigPath)
}

// GetOutputConfigPath returns the path to a file where the raw Envoy bootstrap configuration will
// be persisted for debugging purposes
func GetOutputConfigPath() string {
	return getEnvOrDefault(outputConfigPathEnv, defaultOutputConfigPath)
}

// GetEnvoyExecutable returns the Envoy executable
func GetEnvoyExecutable() string {
	return getEnvOrDefault(envoyExecutableEnv, defaultEnvoyExecutable)
}

// getEnvOrDefault returns the value of the environment variable, if one exists, or a default string
func getEnvOrDefault(envName, defaultValue string) string {
	maybeEnvValue := os.Getenv(envName)
	if maybeEnvValue != "" {
		return maybeEnvValue
	}
	return defaultValue
}

// getAndTransformConfig reads a file, transforms it using the Downward API
// and returns the transformed configuration
func getAndTransformConfig(inputFile string) (string, error) {
	inReader, err := os.Open(inputFile)
	if err != nil {
		return "", err
	}
	defer inReader.Close()

	var buffer bytes.Buffer
	err = downward.Transform(inReader, &buffer)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// crashSignal reports whether the process died of a signal that envoy raises itself when it crashes
// (a fault or an abort), as opposed to SIGKILL or SIGTERM, which come from outside the process. A
// crash caused by the config recurs every time that config is validated. A crash can also come from
// the pod's environment, such as an abort after a failed allocation under memory pressure, so it is
// handled as an interruption.
func crashSignal(exitErr *exec.ExitError) (syscall.Signal, bool) {
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok || !status.Signaled() {
		return 0, false
	}
	switch sig := status.Signal(); sig {
	case syscall.SIGSEGV, syscall.SIGABRT, syscall.SIGBUS, syscall.SIGILL, syscall.SIGFPE, syscall.SIGTRAP:
		return sig, true
	default:
		return sig, false
	}
}
