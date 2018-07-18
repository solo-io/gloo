package file_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/file"
	"io/ioutil"
	"time"
	"github.com/solo-io/solo-kit/test/mocks"
	"os"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
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
		r, err := client.Write(&mocks.MockResource{
			Data: data,
			Metadata: core.Metadata{
				Name: data,
			},
		}, clients.WriteOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(r).To(BeAssignableToTypeOf(&mocks.MockResource{}))
		Expect(r.GetMetadata().Name).To(Equal(data))
		Expect(r.GetMetadata().Namespace).To(Equal(clients.DefaultNamespace))
		Expect(r.GetMetadata().ResourceVersion).To(Equal("1"))
		Expect(r.(*mocks.MockResource).Data).To(Equal(data))
	})
})
