package helpers

import (
	"os"
	"path/filepath"
)

func SoloDirectory() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "solo-io")
}

func HelmDirectory() string {
	return filepath.Join(SoloDirectory(), "gloo-install", "helm", "gloo")
}

func E2eDirectory() string {
	return filepath.Join(SoloDirectory(), "gloo-testing", "e2e")
}

func CertsDirectory() string {
	return filepath.Join(E2eDirectory(), "certs")
}

func ServerCert() string {
	return filepath.Join(CertsDirectory(), "test-ingress.crt")
}
func ServerKey() string {
	return filepath.Join(CertsDirectory(), "test-ingress.key")
}
