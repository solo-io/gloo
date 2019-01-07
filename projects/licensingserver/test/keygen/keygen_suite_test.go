package keygen_test

import (
	"os"
	"testing"

	"github.com/solo-io/solo-kit/pkg/utils/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKeygen(t *testing.T) {
	RegisterFailHandler(Fail)
	if os.Getenv(RUN_KEYGEN_TESTS) == "1" {
		RunSpecs(t, "keygen test Suite")
	} else {
		log.Printf("Skipping keygen test suite, to run set RUN_KEYGEN_TESTS=1")
	}
}
