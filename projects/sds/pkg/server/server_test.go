package server_test

import (
	"context"
	"strings"
	"time"

	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	envoy_service_secret_v3 "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/sds/pkg/server"
	"github.com/solo-io/gloo/projects/sds/pkg/testutils"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
)

var _ = Describe("SDS Server", func() {

	var (
		fs                        afero.Fs
		dir                       string
		keyFile, certFile, caFile afero.File
		err                       error
		serverAddr                = "127.0.0.1:8888"
		sdsClient                 = "test-client"
		srv                       *server.Server
	)

	BeforeEach(func() {
		fs = afero.NewOsFs()
		dir, err = afero.TempDir(fs, "", "")
		Expect(err).NotTo(HaveOccurred())

		keyPEM, certPEM, caPEM := testutils.MustSelfSignedPEM()

		keyFile, err = afero.TempFile(fs, dir, "")
		Expect(err).NotTo(HaveOccurred())
		_, err = keyFile.Write(keyPEM)
		Expect(err).NotTo(HaveOccurred())

		certFile, err = afero.TempFile(fs, dir, "")
		Expect(err).NotTo(HaveOccurred())
		_, err = certFile.Write(certPEM)
		Expect(err).NotTo(HaveOccurred())

		caFile, err = afero.TempFile(fs, dir, "")
		Expect(err).NotTo(HaveOccurred())
		_, err = caFile.Write(caPEM)
		Expect(err).NotTo(HaveOccurred())

		secrets := []server.Secret{
			{
				ServerCert:        "test-server",
				SslCaFile:         caFile.Name(),
				SslCertFile:       certFile.Name(),
				SslKeyFile:        keyFile.Name(),
				ValidationContext: "test-validation",
			},
		}
		srv = server.SetupEnvoySDS(secrets, sdsClient, serverAddr)
	})

	AfterEach(func() {
		_ = fs.RemoveAll(dir)
	})

	It("correctly reads tls secrets from files to generate snapshot version", func() {
		certs, err := testutils.FilesToBytes(keyFile.Name(), certFile.Name(), caFile.Name())
		Expect(err).NotTo(HaveOccurred())

		snapshotVersionBefore, err := server.GetSnapshotVersion(certs)
		Expect(err).NotTo(HaveOccurred())

		_, err = keyFile.WriteString(`newFileString`)
		Expect(err).NotTo(HaveOccurred())
		certs, err = testutils.FilesToBytes(keyFile.Name(), certFile.Name(), caFile.Name())
		Expect(err).NotTo(HaveOccurred())

		snapshotVersionAfter, err := server.GetSnapshotVersion(certs)
		Expect(err).NotTo(HaveOccurred())
		Expect(snapshotVersionAfter).NotTo(Equal(snapshotVersionBefore))
	})

	Context("Test gRPC Server", func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)

		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())
			_, err = srv.Run(ctx)
			// Give it a second to come up + read the certs
			time.Sleep(time.Second * 1)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			cancel()
		})

		It("accepts client connections & updates secrets", func() {
			// Check that it's answering
			var conn *grpc.ClientConn

			// Initiate a connection with the server
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			Expect(err).NotTo(HaveOccurred())
			defer conn.Close()

			client := envoy_service_secret_v3.NewSecretDiscoveryServiceClient(conn)

			// Before any snapshot is set, expect an error when fetching secrets
			resp, err := client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			Expect(err).To(HaveOccurred())

			// After snapshot is set, expect to see the secrets
			srv.UpdateSDSConfig(ctx)
			resp, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetResources()).To(HaveLen(2))
			Expect(resp.Validate()).To(Succeed())

			// Check that the resources contain the expected data
			for _, resource := range resp.GetResources() {
				if resource.GetTypeUrl() == "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret" {
					resourceData := resource.String()
					if strings.Contains(resourceData, "test-server") {
						Expect(resource.String()).To(ContainSubstring("certificate_chain"))
						Expect(resource.String()).To(ContainSubstring("private_key"))
					} else if strings.Contains(resourceData, "test-validation") {
						Expect(resource.String()).To(ContainSubstring("trusted_ca"))
					}
				}
			}
		})
	})
})
