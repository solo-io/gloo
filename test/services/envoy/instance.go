package envoy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"

	"github.com/solo-io/go-utils/threadsafe"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"

	"github.com/solo-io/gloo/test/services"

	"sync"
	"text/template"
	"time"

	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"

	"github.com/onsi/ginkgo/v2"

	errors "github.com/rotisserie/eris"
)

const (
	DefaultProxyName = "default~proxy"
)

type Instance struct {
	defaultBootstrapTemplate *template.Template

	AccessLogAddr string
	AccessLogPort uint32
	// Path to access logs for binary run
	AccessLogs string

	RatelimitAddr string
	RatelimitPort uint32
	ID            string
	Role          string
	envoypath     string
	envoycfg      string
	logs          *SafeBuffer
	LogLevel      string
	cmd           *exec.Cmd
	GlooAddr      string // address for gloo and services
	Port          uint32
	RestXdsPort   uint32

	// Envoy API Version to use, default to V3
	ApiVersion string

	DockerOptions
	UseDocker           bool
	DockerImage         string
	DockerContainerName string

	*RequestPorts

	// adminApiClient represents the client that can be used to access the Envoy Admin API
	adminApiClient *admincli.Client
}

// RequestPorts are the ports that the Instance will listen on for requests
type RequestPorts struct {
	HttpPort   uint32
	HttpsPort  uint32
	HybridPort uint32
	TcpPort    uint32
	AdminPort  uint32
}

// DockerOptions contains extra options for running in docker
type DockerOptions struct {
	// Extra volume arguments
	Volumes []string
	// Extra env arguments.
	// see https://docs.docker.com/engine/reference/run/#env-environment-variables for more info
	Env []string
}

// Deprecated: use RunWith instead
func (ei *Instance) Run(port int) error {
	return ei.RunWith(RunConfig{
		Role:    DefaultProxyName,
		Port:    uint32(port),
		Context: context.TODO(),
	})
}

// Deprecated: use RunWith instead
func (ei *Instance) RunWithRole(role string, port int) error {
	return ei.RunWith(RunConfig{
		Role:    role,
		Port:    uint32(port),
		Context: context.TODO(),
	})
}

// Deprecated: use RunWith instead
func (ei *Instance) RunWithRoleAndRestXds(role string, glooPort, restXdsPort int) error {
	return ei.RunWith(RunConfig{
		Role:        role,
		Port:        uint32(glooPort),
		RestXdsPort: uint32(restXdsPort),
		Context:     context.TODO(),
	})
}

func (ei *Instance) RunWith(runConfig RunConfig) error {
	return ei.runWithAll(runConfig, &templateBootstrapBuilder{
		template: ei.defaultBootstrapTemplate,
	})
}

func (ei *Instance) RunWithConfigFile(ctx context.Context, port int, configFile string) error {
	runConfig := RunConfig{
		Role:    "gloo-system~gateway-proxy",
		Port:    uint32(port),
		Context: ctx,
	}
	boostrapBuilder := &fileBootstrapBuilder{
		file: configFile,
	}
	return ei.runWithAll(runConfig, boostrapBuilder)
}

type RunConfig struct {
	Context context.Context

	Role        string
	Port        uint32
	RestXdsPort uint32
}

func (ei *Instance) runWithAll(runConfig RunConfig, bootstrapBuilder bootstrapBuilder) error {
	go func() {
		<-runConfig.Context.Done()
		ei.Clean()
	}()
	if ei.ID == "" {
		ei.ID = "ingress~for-testing"
	}
	ei.Role = runConfig.Role
	ei.Port = runConfig.Port
	ei.RestXdsPort = runConfig.RestXdsPort
	ei.envoycfg = bootstrapBuilder.Build(ei)

	// construct a client that can be used to access the Admin API
	ei.adminApiClient = admincli.NewClient().
		WithReceiver(ginkgo.GinkgoWriter).
		WithCurlOptions(
			curl.WithPort(int(ei.AdminPort)),
			// We include the verbose output of requests so that we have more information
			// if a test fails
			curl.VerboseOutput(),
			// To reduce potential test flakes, we rely on some basic retries in requests to the Envoy Admin API
			curl.WithRetries(3, 0, 10),
		)

	if ei.UseDocker {
		return ei.runContainer(runConfig.Context)
	}

	args := []string{"--config-yaml", ei.envoycfg, "--disable-hot-restart", "--log-level", ei.LogLevel}

	// run directly
	cmd := exec.CommandContext(runConfig.Context, ei.envoypath, args...)

	safeBuffer := &SafeBuffer{
		buffer: &bytes.Buffer{},
	}
	ei.logs = safeBuffer
	w := io.MultiWriter(ginkgo.GinkgoWriter, safeBuffer)
	cmd.Stdout = w
	cmd.Stderr = w

	err := cmd.Start()
	if err != nil {
		return err
	}
	ei.cmd = cmd

	return ei.waitForEnvoyToBeRunning()
}

func (ei *Instance) Binary() string {
	return ei.envoypath
}

func (ei *Instance) LocalAddr() string {
	return ei.GlooAddr
}

func (ei *Instance) EnablePanicMode() error {
	return ei.setRuntimeConfiguration(
		map[string]string{
			"upstream.healthy_panic_threshold": "100",
		})
}

func (ei *Instance) DisablePanicMode() error {
	return ei.setRuntimeConfiguration(
		map[string]string{
			"upstream.healthy_panic_threshold": "0",
		})
}

func (ei *Instance) setRuntimeConfiguration(queryParameters map[string]string) error {
	return ei.AdminClient().ModifyRuntimeConfiguration(context.Background(), queryParameters)
}

func (ei *Instance) Clean() {
	if ei == nil {
		return
	}
	_ = ei.AdminClient().ShutdownServer(context.Background())

	if ei.cmd != nil {
		ei.cmd.Process.Kill()
		ei.cmd.Wait()
	}

	if ei.UseDocker {
		// An earlier call to quitquitquit should kill and exit the container
		// This is just a backup to make sure it really gets deleted
		services.MustStopAndRemoveContainer(ei.DockerContainerName)
	}
}

func (ei *Instance) runContainer(ctx context.Context) error {
	args := []string{"run", "--rm", "--name", ei.DockerContainerName,
		"-p", fmt.Sprintf("%d:%d", ei.HttpPort, ei.HttpPort),
		"-p", fmt.Sprintf("%d:%d", ei.HttpsPort, ei.HttpsPort),
		"-p", fmt.Sprintf("%d:%d", ei.TcpPort, ei.TcpPort),
		"-p", fmt.Sprintf("%d:%d", ei.HybridPort, ei.HybridPort),
		"-p", fmt.Sprintf("%d:%d", ei.AdminPort, ei.AdminPort),
	}

	for _, volume := range ei.DockerOptions.Volumes {
		args = append(args, "-v", volume)
	}

	for _, env := range ei.DockerOptions.Env {
		args = append(args, "-e", env)
	}

	args = append(args,
		"--entrypoint=envoy",
		ei.DockerImage,
		"--disable-hot-restart",
		"--log-level", ei.LogLevel,
		"--config-yaml", ei.envoycfg,
	)

	fmt.Fprintln(ginkgo.GinkgoWriter, args)
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "Unable to start envoy container")
	}

	// cmd.Run() is entering an infinite loop here (not sure why).
	// This is a temporary workaround to poll the container until the admin port is ready for traffic
	return ei.waitForEnvoyToBeRunning()
}

func (ei *Instance) waitForEnvoyToBeRunning() error {
	pingInterval := time.Tick(time.Millisecond * 100)
	pingDuration := time.Second * 10
	pingEndpoint := fmt.Sprintf("localhost:%d", ei.AdminPort)

	ctx, cancel := context.WithTimeout(context.Background(), pingDuration)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return errors.Errorf("timed out waiting for envoy on %s", pingEndpoint)

		case <-pingInterval:
			conn, _ := net.Dial("tcp", pingEndpoint)
			if conn != nil {
				conn.Close()
				return nil
			}
			continue
		}
	}
}

func (ei *Instance) Logs() (string, error) {
	if ei.UseDocker {
		logsArgs := []string{"logs", ei.DockerContainerName}
		cmd := exec.Command("docker", logsArgs...)
		byt, err := cmd.CombinedOutput()
		if err != nil {
			return "", errors.Wrap(err, "Unable to fetch logs from envoy container")
		}
		return string(byt), nil
	}

	return ei.logs.String(), nil
}

// Deprecated: Prefer StructuredConfigDump
func (ei *Instance) ConfigDump() (string, error) {
	var outLocation threadsafe.Buffer

	err := ei.AdminClient().ConfigDumpCmd(context.Background(), nil).WithStdout(&outLocation).Run().Cause()
	if err != nil {
		return "", err
	}

	return outLocation.String(), nil
}

func (ei *Instance) StructuredConfigDump() (*adminv3.ConfigDump, error) {
	return ei.AdminClient().GetConfigDump(context.Background(), nil)
}

func (ei *Instance) Statistics() (string, error) {
	return ei.AdminClient().GetStats(context.Background(), nil)
}

func (ei *Instance) AdminClient() *admincli.Client {
	return ei.adminApiClient
}

// SafeBuffer is a goroutine safe bytes.Buffer
type SafeBuffer struct {
	buffer *bytes.Buffer
	mutex  sync.Mutex
}

// Write appends the contents of p to the buffer, growing the buffer as needed. It returns
// the number of bytes written.
func (s *SafeBuffer) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.Write(p)
}

// String returns the contents of the unread portion of the buffer
// as a string.  If the Buffer is a nil pointer, it returns "<nil>".
func (s *SafeBuffer) String() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.String()
}
