package run_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/solo-io/gloo/projects/sds/pkg/run"
	"github.com/solo-io/gloo/projects/sds/pkg/server"
	"github.com/solo-io/gloo/projects/sds/pkg/testutils"
)

func TestRun_returnsWhenContextCanceled(t *testing.T) {
	keyPEM, certPEM, caPEM := testutils.MustSelfSignedPEM()

	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.pem")
	certPath := filepath.Join(dir, "cert.pem")
	caPath := filepath.Join(dir, "ca.pem")
	writeFile(t, keyPath, keyPEM)
	writeFile(t, certPath, certPEM)
	writeFile(t, caPath, caPEM)

	secret := server.Secret{
		ServerCert:        "test-cert",
		ValidationContext: "test-validation-context",
		SslCaFile:         caPath,
		SslCertFile:       certPath,
		SslKeyFile:        keyPath,
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- run.Run(ctx, []server.Secret{secret}, "test-client", "127.0.0.1:0")
	}()

	time.Sleep(250 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned unexpected error after context cancel: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}

func writeFile(t *testing.T, path string, contents []byte) {
	t.Helper()
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}
}
