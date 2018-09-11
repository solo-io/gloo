package propagator_test

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	. "github.com/solo-io/solo-kit/pkg/api/v1/propagator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/test/mocks"
)

var (
	badStatus = core.Status{
		State:  core.Status_Rejected,
		Reason: "it gave me gas",
	}
	goodStatus = core.Status{
		State: core.Status_Accepted,
	}
	pendingStatus = core.Status{
		State: core.Status_Pending,
	}
)

var _ = Describe("Propagator", func() {
	var tmpdir string
	BeforeEach(func() {
		var err error
		tmpdir, err = ioutil.TempDir("", "propagator-test")
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})
	It("propagates errors from a set of child resources to a set of parent resources", func() {
		parent1 := mocks.NewMockResource("namespace1", "parent1")
		parent2 := mocks.NewFakeResource("namespace2", "parent2")
		parents := resources.InputResourceList{
			parent1,
			parent2,
		}
		child1 := mocks.NewFakeResource("namespace1", "child1")
		child2 := mocks.NewMockResource("namespace2", "child2")
		children := resources.InputResourceList{
			child1,
			child2,
		}
		mockRc, err := mocks.NewMockResourceClient(&factory.FileResourceClientFactory{
			RootDir: tmpdir,
		})
		Expect(err).NotTo(HaveOccurred())
		fakeRc, err := mocks.NewFakeResourceClient(&factory.FileResourceClientFactory{
			RootDir: tmpdir,
		})
		Expect(err).NotTo(HaveOccurred())

		resourceClients := make(clients.ResourceClients)
		resourceClients.Add(mockRc.BaseClient())
		resourceClients.Add(fakeRc.BaseClient())
		prop := NewPropagator("luffy", parents, children, resourceClients)
		ctx, cancel := context.WithCancel(context.Background())
		errs := make(chan error)

		// get em in there
		parent1, err = mockRc.Write(parent1, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		parent2, err = fakeRc.Write(parent2, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		child1, err = fakeRc.Write(child1, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		child2, err = mockRc.Write(child2, clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		// start the guy
		err = prop.PropagateStatuses(errs, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Millisecond,
		})
		Expect(err).NotTo(HaveOccurred())

		// now update some statuses
		child1.SetStatus(badStatus)
		child2.SetStatus(goodStatus)

		child1, err = fakeRc.Write(child1, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())
		child2, err = mockRc.Write(child2, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())

	l:
		for {
			select {
			case <-time.After(time.Second * 3):
				break l
			case err := <-errs:
				log.Print(err)
				Expect(err.Error()).NotTo(ContainSubstring("resource version error"))
			}
		}

		// parents should (eventually) have a bad status
		Eventually(func() (core.Status, error) {
			parent1, err = mockRc.Read(parent1.Metadata.Namespace, parent1.Metadata.Name, clients.ReadOpts{Ctx: ctx})
			if err != nil {
				return core.Status{}, err
			}
			return parent1.GetStatus(), err
		}, time.Second*5).Should(Equal(core.Status{
			State:      0,
			Reason:     "",
			ReportedBy: "",
			SubresourceStatuses: map[string]*core.Status{
				"*mocks.FakeResource.namespace1.child1": {
					State:               2,
					Reason:              "it gave me gas",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
				"*mocks.MockResource.namespace2.child2": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
			},
		}))
		Eventually(func() (core.Status, error) {
			parent2, err = fakeRc.Read(parent2.Metadata.Namespace, parent2.Metadata.Name, clients.ReadOpts{Ctx: ctx})
			if err != nil {
				return core.Status{}, err
			}
			return parent2.GetStatus(), err
		}, time.Second*5).Should(Equal(core.Status{
			State:      0,
			Reason:     "",
			ReportedBy: "",
			SubresourceStatuses: map[string]*core.Status{
				"*mocks.MockResource.namespace2.child2": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
				"*mocks.FakeResource.namespace1.child1": {
					State:               2,
					Reason:              "it gave me gas",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
			},
		}))

		//
		// try to see accepted after both children accepted
		//
		child1.SetStatus(goodStatus)
		child2.SetStatus(goodStatus)

		child1, err = fakeRc.Write(child1, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())
		child2, err = mockRc.Write(child2, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())

	le:
		for {
			select {
			case <-time.After(time.Second * 3):
				break le
			case err := <-errs:
				log.Print(err)
				Expect(err.Error()).NotTo(ContainSubstring("resource version error"))
			}
		}

		// parents should (eventually) have a good status
		Eventually(func() (core.Status, error) {
			parent1, err = mockRc.Read(parent1.Metadata.Namespace, parent1.Metadata.Name, clients.ReadOpts{Ctx: ctx})
			if err != nil {
				return core.Status{}, err
			}
			return parent1.GetStatus(), err
		}, time.Second*5).Should(Equal(core.Status{
			State:      0,
			Reason:     "",
			ReportedBy: "",
			SubresourceStatuses: map[string]*core.Status{
				"*mocks.FakeResource.namespace1.child1": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
				"*mocks.MockResource.namespace2.child2": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
			},
		}))
		Eventually(func() (core.Status, error) {
			parent2, err = fakeRc.Read(parent2.Metadata.Namespace, parent2.Metadata.Name, clients.ReadOpts{Ctx: ctx})
			if err != nil {
				return core.Status{}, err
			}
			return parent2.GetStatus(), err
		}, time.Second*5).Should(Equal(core.Status{
			State:      0,
			Reason:     "",
			ReportedBy: "",
			SubresourceStatuses: map[string]*core.Status{
				"*mocks.MockResource.namespace2.child2": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
				"*mocks.FakeResource.namespace1.child1": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
			},
		}))

		//
		// try it again after cancel
		//
		cancel()

		// want to see old status
		child1.SetStatus(badStatus)
		child2.SetStatus(badStatus)

		child1, err = fakeRc.Write(child1, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())
		child2, err = mockRc.Write(child2, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())

		// wait on all resources to be read in
	lep:
		for {
			select {
			case <-time.After(time.Second * 3):
				break lep
			case err := <-errs:
				log.Print(err)
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).NotTo(ContainSubstring("resource version error"))
			}
		}

		// parents should (eventually) have a good status
		Eventually(func() (core.Status, error) {
			parent1, err = mockRc.Read(parent1.Metadata.Namespace, parent1.Metadata.Name, clients.ReadOpts{Ctx: ctx})
			if err != nil {
				return core.Status{}, err
			}
			return parent1.GetStatus(), err
		}, time.Second*5).Should(Equal(core.Status{
			State:      0,
			Reason:     "",
			ReportedBy: "",
			SubresourceStatuses: map[string]*core.Status{
				"*mocks.FakeResource.namespace1.child1": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
				"*mocks.MockResource.namespace2.child2": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
			},
		}))
		Eventually(func() (core.Status, error) {
			parent2, err = fakeRc.Read(parent2.Metadata.Namespace, parent2.Metadata.Name, clients.ReadOpts{Ctx: ctx})
			if err != nil {
				return core.Status{}, err
			}
			return parent2.GetStatus(), err
		}, time.Second*5).Should(Equal(core.Status{
			State:      0,
			Reason:     "",
			ReportedBy: "",
			SubresourceStatuses: map[string]*core.Status{
				"*mocks.MockResource.namespace2.child2": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
				"*mocks.FakeResource.namespace1.child1": {
					State:               1,
					Reason:              "",
					ReportedBy:          "",
					SubresourceStatuses: nil,
				},
			},
		}))
	})
})
