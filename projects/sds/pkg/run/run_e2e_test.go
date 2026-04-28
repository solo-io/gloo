package run_test

import (
	"context"
	"os"
	"path"
	"time"

	envoy_extensions_transport_sockets_tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	envoy_service_secret_v3 "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/solo-io/gloo/projects/sds/pkg/run"
	"github.com/solo-io/gloo/projects/sds/pkg/server"
	"github.com/solo-io/gloo/projects/sds/pkg/testutils"
	"github.com/spf13/afero"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SDS Server E2E Test", Serial, func() {

	// These tests use the Serial decorator because they rely on a hard-coded port for the SDS server (8236)

	var (
		ctx                                                             context.Context
		cancel                                                          context.CancelFunc
		err                                                             error
		fs                                                              afero.Fs
		dir                                                             string
		keyName, certName, caName, ocspName                             string
		keyNameSymlink, certNameSymlink, caNameSymlink, ocspNameSymlink string
		secret                                                          server.Secret
		testServerAddress                                               = "127.0.0.1:8236"
		sdsClient                                                       = "test-client"
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		fs = afero.NewOsFs()
		dir, err = afero.TempDir(fs, "", "")
		Expect(err).NotTo(HaveOccurred())

		keyPEM, certPEM, caPEM := testutils.MustSelfSignedPEM()
		ocspPEM := []byte("-----BEGIN OCSP RESPONSE-----\nAQIDBAU=\n-----END OCSP RESPONSE-----\n")

		// Kubernetes mounts secrets as a symlink to a ..data directory, so we'll mimic that here
		keyName = path.Join(dir, "/", "tls.key-0")
		certName = path.Join(dir, "/", "tls.crt-0")
		caName = path.Join(dir, "/", "ca.crt-0")
		ocspName = path.Join(dir, "/", "tls.ocsp-staple-0")
		err = afero.WriteFile(fs, keyName, keyPEM, 0644)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(fs, certName, certPEM, 0644)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(fs, caName, caPEM, 0644)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(fs, ocspName, ocspPEM, 0644)
		Expect(err).NotTo(HaveOccurred())
		keyNameSymlink = path.Join(dir, "/", "tls.key")
		certNameSymlink = path.Join(dir, "/", "tls.crt")
		caNameSymlink = path.Join(dir, "/", "ca.crt")
		ocspNameSymlink = path.Join(dir, "/", "tls.ocsp-staple")
		err = os.Symlink(keyName, keyNameSymlink)
		Expect(err).NotTo(HaveOccurred())
		err = os.Symlink(certName, certNameSymlink)
		Expect(err).NotTo(HaveOccurred())
		err = os.Symlink(caName, caNameSymlink)
		Expect(err).NotTo(HaveOccurred())
		err = os.Symlink(ocspName, ocspNameSymlink)
		Expect(err).NotTo(HaveOccurred())

		secret = server.Secret{
			ServerCert:        "test-cert",
			ValidationContext: "test-validation-context",
			SslCaFile:         caName,
			SslCertFile:       certName,
			SslKeyFile:        keyName,
			SslOcspFile:       ocspName,
		}
	})

	AfterEach(func() {
		cancel()

		_ = fs.RemoveAll(dir)
	})

	It("runs and stops correctly", func() {

		go func() {
			defer GinkgoRecover()

			if err := run.Run(ctx, []server.Secret{secret}, sdsClient, testServerAddress); err != nil {
				Expect(err).NotTo(HaveOccurred())
			}
		}()

		// Connect with the server
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(testServerAddress, grpc.WithInsecure())
		Expect(err).NotTo(HaveOccurred())
		defer conn.Close()
		client := envoy_service_secret_v3.NewSecretDiscoveryServiceClient(conn)

		// SDS should serve a valid snapshot once the server has read matching cert material
		Eventually(func() bool {
			_, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			return err == nil
		}, "10s", "200ms").Should(BeTrue())

		// Cancel the context in order to stop the gRPC server
		cancel()

		Eventually(func() bool {
			_, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			return err != nil
		}, "10s", "200ms").Should(BeTrue())

	})

	DescribeTable("correctly picks up cert rotations", func(useOcsp bool) {
		go func() {
			defer GinkgoRecover()
			if !useOcsp {
				secret.SslOcspFile = ""
			}
			_ = run.Run(ctx, []server.Secret{secret}, sdsClient, testServerAddress)
		}()

		// Spin up, initial SDS read, and debounce window after any watcher noise
		time.Sleep(2 * time.Second)

		var conn *grpc.ClientConn
		conn, err = grpc.Dial(testServerAddress, grpc.WithInsecure())
		Expect(err).NotTo(HaveOccurred())
		defer conn.Close()
		client := envoy_service_secret_v3.NewSecretDiscoveryServiceClient(conn)

		paths := []string{keyNameSymlink, certNameSymlink, caNameSymlink}
		if useOcsp {
			paths = append(paths, ocspNameSymlink)
		}
		certs, err := testutils.FilesToBytes(paths...)
		Expect(err).NotTo(HaveOccurred())

		snapshotVersion, err := server.GetSnapshotVersion(certs)
		Expect(err).NotTo(HaveOccurred())

		var resp *envoy_service_discovery_v3.DiscoveryResponse
		assertOCSPState := func(resp *envoy_service_discovery_v3.DiscoveryResponse) {
			var serverSecret *envoy_extensions_transport_sockets_tls_v3.Secret
			for _, resource := range resp.GetResources() {
				parsed := new(envoy_extensions_transport_sockets_tls_v3.Secret)
				Expect(resource.UnmarshalTo(parsed)).To(Succeed())
				if parsed.GetName() == secret.ServerCert {
					serverSecret = parsed
					break
				}
			}
			Expect(serverSecret).NotTo(BeNil())
			Expect(serverSecret.GetTlsCertificate()).NotTo(BeNil())
			if useOcsp {
				ocspBytes, err := os.ReadFile(ocspNameSymlink)
				Expect(err).NotTo(HaveOccurred())
				Expect(serverSecret.GetTlsCertificate().GetOcspStaple().GetInlineBytes()).To(Equal(ocspBytes))
			} else {
				Expect(serverSecret.GetTlsCertificate().GetOcspStaple()).To(BeNil())
			}
		}

		Eventually(func() bool {
			resp, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			return err == nil
		}, "15s", "1s").Should(BeTrue())

		Eventually(func() bool {
			resp, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			return err == nil && resp.VersionInfo == snapshotVersion
		}, "20s", "500ms").Should(BeTrue())
		assertOCSPState(resp)

		// Cert rotation #1 — replace with a new matching key/cert/CA triple
		k1, c1, ca1 := testutils.MustSelfSignedPEMRotation1()
		Expect(os.Remove(keyName)).To(Succeed())
		Expect(afero.WriteFile(fs, keyName, k1, 0644)).To(Succeed())
		Expect(os.Remove(certName)).To(Succeed())
		Expect(afero.WriteFile(fs, certName, c1, 0644)).To(Succeed())
		Expect(os.Remove(caName)).To(Succeed())
		Expect(afero.WriteFile(fs, caName, ca1, 0644)).To(Succeed())

		certs, err = testutils.FilesToBytes(paths...)
		Expect(err).NotTo(HaveOccurred())
		snapshotVersion, err = server.GetSnapshotVersion(certs)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() bool {
			resp, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			if err != nil {
				return false
			}
			return resp.VersionInfo == snapshotVersion
		}, "20s", "500ms").Should(BeTrue())
		assertOCSPState(resp)

		// Cert rotation #2
		k2, c2, ca2 := testutils.MustSelfSignedPEMRotation2()
		Expect(os.Remove(keyName)).To(Succeed())
		Expect(afero.WriteFile(fs, keyName, k2, 0644)).To(Succeed())
		Expect(os.Remove(certName)).To(Succeed())
		Expect(afero.WriteFile(fs, certName, c2, 0644)).To(Succeed())
		Expect(os.Remove(caName)).To(Succeed())
		Expect(afero.WriteFile(fs, caName, ca2, 0644)).To(Succeed())

		certs, err = testutils.FilesToBytes(paths...)
		Expect(err).NotTo(HaveOccurred())
		snapshotVersion, err = server.GetSnapshotVersion(certs)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() bool {
			resp, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			if err != nil {
				return false
			}
			return resp.VersionInfo == snapshotVersion
		}, "20s", "500ms").Should(BeTrue())
		assertOCSPState(resp)
	},
		Entry("with ocsp", true),
		Entry("without ocsp", false),
	)
})
