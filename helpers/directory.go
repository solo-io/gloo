package helpers

import (
	"os"
	"path/filepath"
)

func SoloDirectory() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "solo-io")
}

func E2eDirectory() string {
	return filepath.Join(SoloDirectory(), "gloo-testing", "e2e")
}
