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
	"github.com/solo-io/solo-kit/pkg/errors"
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
		name := "foo"
		input := &mocks.MockResource{
			Data: name,
			Metadata: core.Metadata{
				Name: name,
			},
		}
		r, err := client.Write(input, clients.WriteOptions{})
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Write(input, clients.WriteOptions{})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsExist(err)).To(BeTrue())

		Expect(r).To(BeAssignableToTypeOf(&mocks.MockResource{}))
		Expect(r.GetMetadata().Name).To(Equal(name))
		Expect(r.GetMetadata().Namespace).To(Equal(clients.DefaultNamespace))
		Expect(r.GetMetadata().ResourceVersion).To(Equal("1"))
		Expect(r.(*mocks.MockResource).Data).To(Equal(name))


		_, err = client.Write(input, clients.WriteOptions{
			OverwriteExisting: true,
		})
		Expect(err).NotTo(HaveOccurred())

		var read mocks.MockResource
		err = client.Read(name, &read, clients.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(r).To(Equal(&read))

		err = client.Read(name, &read, clients.GetOptions{Namespace: "doesntexist"})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())
	})
})
