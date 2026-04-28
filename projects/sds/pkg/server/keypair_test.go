package server

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	rsrc "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/solo-io/gloo/projects/sds/pkg/testutils"
)

func TestReadAndValidateSecret_mismatchedKeyAndCert(t *testing.T) {
	keyA, _, caA := testutils.MustSelfSignedPEM()
	_, certB, _ := testutils.MustSelfSignedPEMRotation1()

	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.pem")
	certPath := filepath.Join(dir, "cert.pem")
	caPath := filepath.Join(dir, "ca.pem")
	if err := os.WriteFile(keyPath, keyA, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(certPath, certB, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caPath, caA, 0o600); err != nil {
		t.Fatal(err)
	}

	sec := Secret{
		SslKeyFile:        keyPath,
		SslCertFile:       certPath,
		SslCaFile:         caPath,
		ServerCert:        "server",
		ValidationContext: "vc",
	}
	_, _, err := readAndValidateSecret(context.Background(), sec)
	if err == nil {
		t.Fatal("expected error for mismatched key and certificate")
	}
}

func TestReadAndValidateSecret_eventuallyConsistentAfterRotation(t *testing.T) {
	keyA, _, caA := testutils.MustSelfSignedPEM()
	keyB, certB, _ := testutils.MustSelfSignedPEMRotation1()

	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.pem")
	certPath := filepath.Join(dir, "cert.pem")
	caPath := filepath.Join(dir, "ca.pem")
	// Start with A's key and B's cert (mismatch).
	if err := os.WriteFile(keyPath, keyA, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(certPath, certB, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caPath, caA, 0o600); err != nil {
		t.Fatal(err)
	}

	sec := Secret{
		SslKeyFile:        keyPath,
		SslCertFile:       certPath,
		SslCaFile:         caPath,
		ServerCert:        "server",
		ValidationContext: "vc",
	}

	go func() {
		time.Sleep(150 * time.Millisecond)
		if err := os.WriteFile(keyPath, keyB, 0o600); err != nil {
			panic(err)
		}
	}()

	_, _, err := readAndValidateSecret(context.Background(), sec)
	if err != nil {
		t.Fatalf("expected success after key rotated to match cert: %v", err)
	}
}

func TestUpdateSDSConfig_matchingPair(t *testing.T) {
	keyPEM, certPEM, caPEM := testutils.MustSelfSignedPEM()
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.pem")
	certPath := filepath.Join(dir, "cert.pem")
	caPath := filepath.Join(dir, "ca.pem")
	if err := os.WriteFile(keyPath, keyPEM, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(certPath, certPEM, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caPath, caPEM, 0o600); err != nil {
		t.Fatal(err)
	}

	srv := SetupEnvoySDS([]Secret{{
		SslKeyFile:        keyPath,
		SslCertFile:       certPath,
		SslCaFile:         caPath,
		ServerCert:        "server",
		ValidationContext: "vc",
	}}, "client", "127.0.0.1:0")
	ctx := context.Background()
	if err := srv.UpdateSDSConfig(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateSDSConfig_mismatchedPairDoesNotReplaceExistingSnapshot(t *testing.T) {
	keyA, certA, caA := testutils.MustSelfSignedPEM()
	_, certB, _ := testutils.MustSelfSignedPEMRotation1()

	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key.pem")
	certPath := filepath.Join(dir, "cert.pem")
	caPath := filepath.Join(dir, "ca.pem")
	if err := os.WriteFile(keyPath, keyA, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(certPath, certA, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caPath, caA, 0o600); err != nil {
		t.Fatal(err)
	}

	srv := SetupEnvoySDS([]Secret{{
		SslKeyFile:        keyPath,
		SslCertFile:       certPath,
		SslCaFile:         caPath,
		ServerCert:        "server",
		ValidationContext: "vc",
	}}, "client", "127.0.0.1:0")
	ctx := context.Background()
	if err := srv.UpdateSDSConfig(ctx); err != nil {
		t.Fatal(err)
	}

	before, err := srv.snapshotCache.GetSnapshot("client")
	if err != nil {
		t.Fatal(err)
	}
	versionBefore := before.GetVersion(rsrc.SecretType)
	if versionBefore == "" {
		t.Fatal("expected non-empty secret snapshot version")
	}

	if err := os.WriteFile(certPath, certB, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := srv.UpdateSDSConfig(ctx); err == nil {
		t.Fatal("expected update to fail for mismatched key and certificate")
	}

	after, err := srv.snapshotCache.GetSnapshot("client")
	if err != nil {
		t.Fatal(err)
	}
	versionAfter := after.GetVersion(rsrc.SecretType)
	if versionAfter != versionBefore {
		t.Fatalf("expected snapshot version to remain %q after failed update, got %q", versionBefore, versionAfter)
	}
}
