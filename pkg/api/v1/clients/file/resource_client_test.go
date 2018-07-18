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
		client = NewResourceClient(tmpDir, time.Millisecond, &mocks.MockResource{})
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
		r1, err := client.Write(input, clients.WriteOptions{})
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Write(input, clients.WriteOptions{})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsExist(err)).To(BeTrue())

		Expect(r1).To(BeAssignableToTypeOf(&mocks.MockResource{}))
		Expect(r1.GetMetadata().Name).To(Equal(name))
		Expect(r1.GetMetadata().Namespace).To(Equal(clients.DefaultNamespace))
		Expect(r1.GetMetadata().ResourceVersion).To(Equal("1"))
		Expect(r1.(*mocks.MockResource).Data).To(Equal(name))


		_, err = client.Write(input, clients.WriteOptions{
			OverwriteExisting: true,
		})
		Expect(err).NotTo(HaveOccurred())

		read, err := client.Read(name, clients.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(read).To(Equal(r1))

		_, err = client.Read(name, clients.GetOptions{Namespace: "doesntexist"})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		name = "boo"
		input = &mocks.MockResource{
			Data: name,
			Metadata: core.Metadata{
				Name: name,
			},
		}
		r2, err := client.Write(input, clients.WriteOptions{})
		Expect(err).NotTo(HaveOccurred())

		list, err := client.List(clients.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).To(ContainElement(r2))
	})
})
