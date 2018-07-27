package apiclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"fmt"

	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/apiclient"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/mocks"
	"google.golang.org/grpc"
)

var _ = Describe("Base", func() {
	var (
		client *ResourceClient
		cc     *grpc.ClientConn
		err    error
	)
	BeforeEach(func() {
		// give grpc server time to start
		time.Sleep(time.Second)
		cc, err = grpc.Dial(fmt.Sprintf("127.0.0.1:%v", port), grpc.WithInsecure())
		Expect(err).NotTo(HaveOccurred())
		client = NewResourceClient(cc, &mocks.MockResource{})
	})
	AfterEach(func() {
		cc.Close()
	})
	It("CRUDs resources", func() {
		helpers.TestCrudClient(client)
	})
})
