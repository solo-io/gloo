package envoy_admin_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEnvoyAdminClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EnvoyAdmin Suite")
}
