package make_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/docker"
)

const (
	StandardGlooImage = "quay.io/solo-io/gloo:1.16.0-beta1"
	StandardSdsImage  = "quay.io/solo-io/sds:1.16.0-beta1"
	FipsSdsImage      = "quay.io/solo-io/sds-fips:1.17.0-beta2-9037"
)

var _ = BeforeSuite(func() {
	for _, image := range []string{StandardGlooImage, StandardSdsImage, FipsSdsImage} {
		_, err := docker.PullIfNotPresent(context.Background(), image, 3)
		Expect(err).NotTo(HaveOccurred(), "can pull image locally")
	}
})

func TestMake(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Make Suite")
}
