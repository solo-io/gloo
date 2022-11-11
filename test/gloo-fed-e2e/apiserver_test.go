package gloo_fed_e2e_test

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"google.golang.org/grpc"
)

var _ = Describe("Remote Envoy Config Getter", func() {
	var (
		portFwd       *exec.Cmd
		apiserverPort = 10101
	)

	BeforeEach(func() {
		apiserverPort, err = cliutil.GetFreePort()
		Expect(err).NotTo(HaveOccurred())
		portFwd, err = cliutil.PortForward(defaults.GlooSystem, "svc/gloo-fed-console", strconv.Itoa(apiserverPort), "10101", false)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		// kill port forward
		if portFwd.Process != nil {
			portFwd.Process.Kill()
			portFwd.Process.Release()
		}
	})

	It("works for getting remote envoy config", func() {
		serverAddr := "localhost:" + strconv.Itoa(apiserverPort)

		glooInstanceRefName := fmt.Sprintf("%s-%s", managementClusterConfig.KubeContext, defaults.GlooSystem)

		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(serverAddr, opts...)
		Expect(err).NotTo(HaveOccurred())
		defer conn.Close()
		client := rpc_edge_v1.NewGlooInstanceApiClient(conn)
		resp, err := client.GetConfigDumps(context.TODO(), &rpc_edge_v1.GetConfigDumpsRequest{
			GlooInstanceRef: &v1.ObjectRef{
				Name:      glooInstanceRefName,
				Namespace: defaults.GlooSystem,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(resp.GetConfigDumps())).To(Equal(1))
		Expect(resp.GetConfigDumps()[0].Error).To(BeEmpty())
	})
})
