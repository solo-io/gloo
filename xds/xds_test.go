package xds

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"fmt"
	"net"

	"time"

	"github.com/envoyproxy/go-control-plane/api/filter/network"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/glue/config"
	"github.com/solo-io/glue/pkg/log"
	"google.golang.org/grpc"
)

var _ = Describe("Xds", func() {
	envoyBaseDir := os.Getenv("GOPATH") + "/src/github.com/solo-io/glue"
	envoyRunArgs := []string{
		filepath.Join(envoyBaseDir, "envoy"),
		"-c", filepath.Join(envoyBaseDir, "envoy.yaml"),
		"--v2-config-only",
		"--service-cluster", "envoy",
		"--service-node", "envoy",
	}
	var (
		envoyPid    int
		buf         *bytes.Buffer
		srv         *grpc.Server
		envoyConfig envoycache.Cache
	)
	BeforeEach(func() {
		// setup
		envoyCmd := exec.Command("/bin/sh", "-c", strings.Join(envoyRunArgs, " "))
		buf = &bytes.Buffer{}
		envoyCmd.Stdout = buf
		envoyCmd.Stderr = buf
		srv, envoyConfig = createXdsServer()
		go func() {
			if err := envoyCmd.Run(); err != nil {
				log.Fatalf(buf.String() + ": " + err.Error())
			}
		}()
		for envoyCmd.Process == nil {
			time.Sleep(time.Millisecond)
		}
		envoyPid = envoyCmd.Process.Pid
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8081))
		if err != nil {
			log.Fatalf(err.Error())
		}
		go func() {
			if err := srv.Serve(lis); err != nil {
				log.Fatalf(err.Error())
			}
		}()
		resources := mockResources()
		resources[0].Filters = []config.FilterWrapper{
			{
				Filter: network.HttpFilter{
					Name: "envoy.router",
				},
			},
		}
		snapshot, err := createSnapshot(resources)
		if err != nil {
			log.Fatalf(err.Error())
		}
		envoyConfig.SetSnapshot(nodeKey, snapshot)
	})
	AfterEach(func() {
		srv.Stop()
		if err := syscall.Kill(envoyPid, syscall.SIGKILL); err != nil {
			log.Fatalf(err.Error())
		}
	})
	Describe("RunXDS Server", func() {
		FIt("successfully bootstraps the envoy proxy", func() {
			Eventually(func() string {
				str := buf.String()
				return str
			}, time.Second*10).Should(ContainSubstring("lds: add/update listener '" + listener))
		})
	})
})
