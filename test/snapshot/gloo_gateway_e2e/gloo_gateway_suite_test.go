package bugs_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGlooGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gloo Gateway Suite")
}

var (
	ctrl *gomock.Controller
)

var _ = BeforeSuite(func() {
	ctrl, _ = gomock.WithContext(context.TODO(), GinkgoT())
})
