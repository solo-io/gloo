package setup_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

import (
	fdssetup "github.com/solo-io/solo-kit/projects/discovery/pkg/fds/setup"
	gloosetup "github.com/solo-io/solo-kit/projects/gloo/pkg/setup"
)

var _ = Describe("Setup", func() {
	It("works", func() {
		Expect(runFds()).To(BeNil())

	})
})

func runFds() error {
	opts, err := gloosetup.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	return fdssetup.Setup(opts)
}
