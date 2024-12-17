package debug_test

import (
	"fmt"
	cliutil "github.com/solo-io/gloo/pkg/cliutil/install"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDebug(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Debug Suite")

	var currentContext string

	BeforeSuite(func() {
		out, err := cliutil.KubectlOut(nil, "config", "current-context")
		currentContext = string(out)
		Expect(err).NotTo(HaveOccurred(), err.Error()+", "+currentContext)
		fmt.Println("ARIANA BeforeSuite current-context", currentContext)

		Expect(cliutil.Kubectl(nil, "config", "unset", "current-context")).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		Expect(cliutil.Kubectl(nil, "config", "use-context", currentContext)).NotTo(HaveOccurred())
	})

}
