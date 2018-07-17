package file_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/file"
	"io/ioutil"
	"time"
	"github.com/solo-io/solo-kit/test/mocks"
	"os"
)

var _ = Describe("Base", func() {
	var (
		client *ResourceClient
		tmpDir string
	)
	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "base_test")
		Expect(err).NotTo(HaveOccurred())
		client = NewResourceClient(tmpDir, time.Millisecond)
	})
	AfterEach(func() {
		os.RemoveAll(tmpDir)
	})
	It("CRUDs resources", func() {
		data := "foo"
		r, err := client.Create(&mocks.MockResource{
			Data: data,
		}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(r).To(BeAssignableToTypeOf(&mocks.MockResource{}))
		Expect(r.(*mocks.MockResource).Data).To(Equal(data))
	})
})
