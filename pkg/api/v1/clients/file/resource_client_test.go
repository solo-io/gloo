package file_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"

	"time"

	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/file"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/test/mocks"
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
		client = NewResourceClient(tmpDir, &mocks.MockResource{})
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
		r1, err := client.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Write(input, clients.WriteOpts{})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsExist(err)).To(BeTrue())

		Expect(r1).To(BeAssignableToTypeOf(&mocks.MockResource{}))
		Expect(r1.GetMetadata().Name).To(Equal(name))
		Expect(r1.GetMetadata().Namespace).To(Equal(clients.DefaultNamespace))
		Expect(r1.GetMetadata().ResourceVersion).To(Equal("1"))
		Expect(r1.(*mocks.MockResource).Data).To(Equal(name))

		_, err = client.Write(input, clients.WriteOpts{
			OverwriteExisting: true,
		})
		Expect(err).NotTo(HaveOccurred())

		read, err := client.Read(name, clients.GetOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(read).To(Equal(r1))

		_, err = client.Read(name, clients.GetOpts{Namespace: "doesntexist"})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		name = "boo"
		input = &mocks.MockResource{
			Data: name,
			Metadata: core.Metadata{
				Name: name,
			},
		}
		r2, err := client.Write(input, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		list, err := client.List(clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).To(ContainElement(r2))

		err = client.Delete("adsfw", clients.DeleteOpts{})
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotExist(err)).To(BeTrue())

		err = client.Delete("adsfw", clients.DeleteOpts{
			IgnoreNotExist: true,
		})
		Expect(err).NotTo(HaveOccurred())

		err = client.Delete(r2.GetMetadata().Name, clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())
		list, err = client.List(clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(list).To(ContainElement(r1))
		Expect(list).NotTo(ContainElement(r2))

		w, errs, err := client.Watch(clients.WatchOpts{RefreshRate: time.Millisecond})
		Expect(err).NotTo(HaveOccurred())
		go func() {
			// event 1
			r2, err = client.Write(r2, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			time.Sleep(time.Millisecond * 2)
			// event 2
			err = client.Delete(r2.GetMetadata().Name, clients.DeleteOpts{})
			Expect(err).NotTo(HaveOccurred())
			time.Sleep(time.Millisecond * 5)
			// event 3
			r2, err = client.Write(r2, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			log.Printf("reached")
		}()
		Eventually(w, time.Second*3).Should(HaveLen(3))
		list1 := <-w
		Expect(list1).To(ContainElement(r1))
		Expect(list1).To(ContainElement(r2))
		list2 := <-w
		Expect(list2).To(ContainElement(r1))
		Expect(list2).NotTo(ContainElement(r2))
		list3 := <-w
		Expect(list3).To(ContainElement(r1))
		Expect(list3).To(ContainElement(r2))
		Expect(errs).To(HaveLen(0))
	})
})
