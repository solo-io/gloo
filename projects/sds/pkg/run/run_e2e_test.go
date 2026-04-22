package run_test

import (
	"context"
	"os"
	"path"
	"time"

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
		ctx                                            context.Context
		cancel                                         context.CancelFunc
		err                                            error
		fs                                             afero.Fs
		dir                                            string
		keyName, certName, caName                      string
		keyNameSymlink, certNameSymlink, caNameSymlink string
		secret                                         server.Secret
		testServerAddress                              = "127.0.0.1:8236"
		sdsClient                                      = "test-client"
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		fs = afero.NewOsFs()
		dir, err = afero.TempDir(fs, "", "")
		Expect(err).NotTo(HaveOccurred())

		keyPEM, certPEM, caPEM := testutils.MustSelfSignedPEM()

		// Kubernetes mounts secrets as a symlink to a ..data directory, so we'll mimic that here
		keyName = path.Join(dir, "/", "tls.key-0")
		certName = path.Join(dir, "/", "tls.crt-0")
		caName = path.Join(dir, "/", "ca.crt-0")
		err = afero.WriteFile(fs, keyName, keyPEM, 0644)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(fs, certName, certPEM, 0644)
		Expect(err).NotTo(HaveOccurred())
		err = afero.WriteFile(fs, caName, caPEM, 0644)
		Expect(err).NotTo(HaveOccurred())
		keyNameSymlink = path.Join(dir, "/", "tls.key")
		certNameSymlink = path.Join(dir, "/", "tls.crt")
		caNameSymlink = path.Join(dir, "/", "ca.crt")
		err = os.Symlink(keyName, keyNameSymlink)
		Expect(err).NotTo(HaveOccurred())
		err = os.Symlink(certName, certNameSymlink)
		Expect(err).NotTo(HaveOccurred())
		err = os.Symlink(caName, caNameSymlink)
		Expect(err).NotTo(HaveOccurred())

		secret = server.Secret{
			ServerCert:        "test-cert",
			ValidationContext: "test-validation-context",
			SslCaFile:         caName,
			SslCertFile:       certName,
			SslKeyFile:        keyName,
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

	It("correctly picks up cert rotations", func() {
		go func() {
			defer GinkgoRecover()
			_ = run.Run(ctx, []server.Secret{secret}, sdsClient, testServerAddress)
		}()

		// Spin up, initial SDS read, and debounce window after any watcher noise
		time.Sleep(2 * time.Second)

		var conn *grpc.ClientConn
		conn, err = grpc.Dial(testServerAddress, grpc.WithInsecure())
		Expect(err).NotTo(HaveOccurred())
		defer conn.Close()
		client := envoy_service_secret_v3.NewSecretDiscoveryServiceClient(conn)

		certs, err := testutils.FilesToBytes(keyNameSymlink, certNameSymlink, caNameSymlink)
		Expect(err).NotTo(HaveOccurred())

		snapshotVersion, err := server.GetSnapshotVersion(certs)
		Expect(err).NotTo(HaveOccurred())

		var resp *envoy_service_discovery_v3.DiscoveryResponse

		Eventually(func() bool {
			resp, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			return err == nil
		}, "15s", "1s").Should(BeTrue())

		Eventually(func() bool {
			resp, err = client.FetchSecrets(ctx, &envoy_service_discovery_v3.DiscoveryRequest{})
			return err == nil && resp.VersionInfo == snapshotVersion
		}, "20s", "500ms").Should(BeTrue())

		// Cert rotation #1 — replace with a new matching key/cert/CA triple
		k1, c1, ca1 := testutils.MustSelfSignedPEMRotation1()
		Expect(os.Remove(keyName)).To(Succeed())
		Expect(afero.WriteFile(fs, keyName, k1, 0644)).To(Succeed())
		Expect(os.Remove(certName)).To(Succeed())
		Expect(afero.WriteFile(fs, certName, c1, 0644)).To(Succeed())
		Expect(os.Remove(caName)).To(Succeed())
		Expect(afero.WriteFile(fs, caName, ca1, 0644)).To(Succeed())

		certs, err = testutils.FilesToBytes(keyNameSymlink, certNameSymlink, caNameSymlink)
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

		// Cert rotation #2
		k2, c2, ca2 := testutils.MustSelfSignedPEMRotation2()
		Expect(os.Remove(keyName)).To(Succeed())
		Expect(afero.WriteFile(fs, keyName, k2, 0644)).To(Succeed())
		Expect(os.Remove(certName)).To(Succeed())
		Expect(afero.WriteFile(fs, certName, c2, 0644)).To(Succeed())
		Expect(os.Remove(caName)).To(Succeed())
		Expect(afero.WriteFile(fs, caName, ca2, 0644)).To(Succeed())

		certs, err = testutils.FilesToBytes(keyNameSymlink, certNameSymlink, caNameSymlink)
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
	})
})
