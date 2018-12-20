package v1

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/gateway/pkg/api/v1"
	v12 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"
)

var _ = Describe("using the changeset resource client factory", func() {

	const changesetName = "cs-name"

	var (
		err              error
		root             string
		csClient         ChangeSetClient
		initialChangeset *ChangeSet
		initialVs        *v1.VirtualService
	)

	// Create a changeset with some initial data
	BeforeEach(func() {

		root, err = ioutil.TempDir("", "changeset_resource_client_test")
		Expect(err).To(BeNil())

		csClient, err = NewChangeSetClient(&factory.FileResourceClientFactory{RootDir: root})
		Expect(err).To(BeNil())

		initialVs = &v1.VirtualService{
			Metadata: core.Metadata{
				Name:            "test-vs-1",
				Namespace:       defaults.GlooSystem,
				ResourceVersion: "1",
			},
			VirtualHost: &v12.VirtualHost{
				Name:    "my-virtual-host-1",
				Domains: []string{"solo.io", "is.awesome"},
				Routes: []*v12.Route{
					{
						Matcher: &v12.Matcher{
							Methods: []string{"GET", "POST"},
						},
					},
				},
			},
		}

		initialChangeset = NewChangeSet(defaults.GlooSystem, changesetName)
		initialChangeset.Branch.Value = "my-branch"
		initialChangeset.PendingAction = Action_NONE
		initialChangeset.Description.Value = "my-description"
		initialChangeset.Data.VirtualServices = []*v1.VirtualService{initialVs}

		_, err = csClient.Write(initialChangeset, clients.WriteOpts{})
		Expect(err).To(BeNil())

		_, err = os.Stat(filepath.Join(root, defaults.GlooSystem, "cs-name.yaml"))
		Expect(err).To(BeNil())
	})

	// Clean up
	AfterEach(func() {
		os.RemoveAll(root)
	})

	Describe("using the factory to create a client for an unsupported resource type", func() {
		It("causes an error", func() {
			_, err = v12.NewSecretClientWithToken(&ChangesetResourceClientFactory{ChangesetClient: csClient}, changesetName)
			Expect(err).To(Not(BeNil()))
		})
	})

	Describe("creating a resource client with a factory that points to a non-existing changeset", func() {
		It("causes an error when we try to use the client", func() {
			invalidClient, err := v1.NewVirtualServiceClientWithToken(
				&ChangesetResourceClientFactory{ChangesetClient: csClient},
				"non-existing-changeset",
			)
			Expect(err).To(BeNil())

			_, err = invalidClient.Read(defaults.GlooSystem, initialVs.Metadata.Name, clients.ReadOpts{})

			Expect(err).To(Not(BeNil()))
			Expect(errors.IsNotExist(err)).To(BeTrue())
			// Verify that the error is actually due to the changeset (and not the virtual service) not existing
			Expect(strings.Contains(err.Error(), "non-existing-changeset")).To(BeTrue())
		})
	})

	Describe("retrieving the resource contained in the initial changeset", func() {

		var (
			vs       *v1.VirtualService
			vsClient v1.VirtualServiceClient
		)

		BeforeEach(func() {
			vsClient, err = v1.NewVirtualServiceClientWithToken(
				&ChangesetResourceClientFactory{ChangesetClient: csClient},
				changesetName,
			)
			Expect(err).To(BeNil())
			vs, err = vsClient.Read(defaults.GlooSystem, initialVs.Metadata.Name, clients.ReadOpts{})
		})

		It("does not cause an error", func() {
			Expect(err).To(BeNil())
		})

		It("retrieves the correct resource", func() {
			Expect(vs.Equal(initialChangeset.Data.VirtualServices[0])).To(BeTrue())
		})
	})

	Describe("adding a resource to the changeset", func() {

		var (
			err      error
			vsClient v1.VirtualServiceClient
		)

		BeforeEach(func() {
			vsClient, err = v1.NewVirtualServiceClientWithToken(
				&ChangesetResourceClientFactory{ChangesetClient: csClient},
				changesetName,
			)
			Expect(err).To(BeNil())
		})

		Context("adding a new resource", func() {

			var newVs *v1.VirtualService

			BeforeEach(func() {
				newVs = &v1.VirtualService{
					Metadata: core.Metadata{Name: "test-vs-2", Namespace: defaults.GlooSystem},
					VirtualHost: &v12.VirtualHost{
						Name:    "my-virtual-host-2",
						Domains: []string{"some.domain"},
						Routes:  []*v12.Route{{Matcher: &v12.Matcher{Headers: []*v12.HeaderMatcher{{Name: "some-header", Value: "header-value"}}}}},
					},
				}
				newVs, err = vsClient.Write(newVs, clients.WriteOpts{})
			})

			It("does not cause an error", func() {
				Expect(err).To(BeNil())
			})

			It("correctly adds the resource to the changeset", func() {
				vs2, err := vsClient.Read(defaults.GlooSystem, newVs.Metadata.Name, clients.ReadOpts{})
				Expect(err).To(BeNil())
				Expect(vs2.Equal(newVs)).To(BeTrue())
			})

			It("increments the changeset edit count", func() {
				cs, err := csClient.Read(defaults.GlooSystem, initialChangeset.Metadata.Name, clients.ReadOpts{})
				Expect(err).To(BeNil())
				Expect(cs.EditCount.Value).To(BeEquivalentTo(1))
			})

			It("results in the changeset containing two resources", func() {
				vsList, err := vsClient.List(defaults.GlooSystem, clients.ListOpts{})
				Expect(err).To(BeNil())
				Expect(len(vsList)).To(BeEquivalentTo(2))
			})

			Describe("deleting a resource from the changeset", func() {

				var initEditCount uint32

				BeforeEach(func() {
					cs, err := csClient.Read(defaults.GlooSystem, initialChangeset.Metadata.Name, clients.ReadOpts{})
					Expect(err).To(BeNil())

					initEditCount = cs.EditCount.Value
				})

				Context("deleting an existing resource", func() {

					BeforeEach(func() {
						err = vsClient.Delete(defaults.GlooSystem, initialVs.Metadata.Name, clients.DeleteOpts{})
					})

					It("does not cause an error", func() {
						Expect(err).To(BeNil())
					})

					It("causes the correct resource to be deleted", func() {
						vsList, err := vsClient.List(defaults.GlooSystem, clients.ListOpts{})
						Expect(err).To(BeNil())
						Expect(len(vsList)).To(BeEquivalentTo(1))

						vs1, err := vsList.Find(defaults.GlooSystem, newVs.Metadata.Name)
						Expect(err).To(BeNil())
						Expect(vs1.Metadata.Name).To(BeEquivalentTo(newVs.Metadata.Name))
					})

					It("increments the changeset edit count", func() {
						cs, err := csClient.Read(defaults.GlooSystem, initialChangeset.Metadata.Name, clients.ReadOpts{})
						Expect(err).To(BeNil())
						Expect(cs.EditCount.Value).To(BeEquivalentTo(initEditCount + 1))
					})
				})

				Context("deleting a non-existing resource", func() {

					Context("ignoring attempts to delete a non-existing resource", func() {

						BeforeEach(func() {
							err = vsClient.Delete(defaults.GlooSystem, "test-vs-X", clients.DeleteOpts{IgnoreNotExist: true})
						})

						It("does not cause an error", func() {
							Expect(err).To(BeNil())
						})

						It("does not increase the edit count", func() {
							cs, err := csClient.Read(defaults.GlooSystem, initialChangeset.Metadata.Name, clients.ReadOpts{})
							Expect(err).To(BeNil())
							Expect(cs.EditCount.Value).To(BeEquivalentTo(initEditCount))
						})
					})

					Context("failing on attempts to delete a non-existing resource", func() {
						It("does cause an error", func() {
							err = vsClient.Delete(defaults.GlooSystem, "test-vs-X", clients.DeleteOpts{IgnoreNotExist: false})
							Expect(err).To(Not(BeNil()))
							Expect(errors.IsNotExist(err)).To(BeTrue())
						})
					})
				})
			})
		})

		Context("trying to add a resource that already exists", func() {
			It("does cause an error", func() {
				_, err = vsClient.Write(initialVs, clients.WriteOpts{})
				Expect(err).To(Not(BeNil()))
				Expect(errors.IsExist(err)).To(BeTrue())
			})

		})
	})

	Describe("updating a resource inside the changeset", func() {

		var (
			initiallyReadVs, updatedVs *v1.VirtualService
			vsClient                   v1.VirtualServiceClient
			newVirtualHostName         = "new-virtual-host-name"
		)

		BeforeEach(func() {
			vsClient, err = v1.NewVirtualServiceClientWithToken(
				&ChangesetResourceClientFactory{ChangesetClient: csClient},
				changesetName,
			)
			Expect(err).To(BeNil())

			initiallyReadVs, err = vsClient.Read(defaults.GlooSystem, initialVs.Metadata.Name, clients.ReadOpts{})
			Expect(err).To(BeNil())

			Expect(err).To(BeNil())

			Expect(initiallyReadVs.VirtualHost.Name).To(BeEquivalentTo("my-virtual-host-1"))
		})

		Context("the resource is not stale", func() {

			var initEditCount uint32

			BeforeEach(func() {
				cs, err := csClient.Read(defaults.GlooSystem, initialChangeset.Metadata.Name, clients.ReadOpts{})
				Expect(err).To(BeNil())
				initEditCount = cs.EditCount.Value

				initiallyReadVs.VirtualHost.Name = newVirtualHostName
				updatedVs, err = vsClient.Write(initiallyReadVs, clients.WriteOpts{OverwriteExisting: true})
			})

			It("does not cause an error", func() {
				Expect(err).To(BeNil())
			})

			It("correctly updates the resource", func() {
				Expect(updatedVs.VirtualHost.Name).To(BeEquivalentTo(newVirtualHostName))
			})

			It("increments the changeset edit count", func() {
				cs, err := csClient.Read(defaults.GlooSystem, initialChangeset.Metadata.Name, clients.ReadOpts{})
				Expect(err).To(BeNil())
				Expect(cs.EditCount.Value).To(BeEquivalentTo(initEditCount + 1))
			})
		})
	})

	Describe("the changeset correctly keeps track of the number of edits", func() {

		It("updates the edit count correctly after multiple mutations", func() {
			cs, err := csClient.Read(defaults.GlooSystem, initialChangeset.Metadata.Name, clients.ReadOpts{})
			Expect(err).To(BeNil())
			initEditCount := cs.EditCount.Value

			vsClient, err := v1.NewVirtualServiceClientWithToken(
				&ChangesetResourceClientFactory{ChangesetClient: csClient},
				changesetName,
			)
			Expect(err).To(BeNil())

			// Read and up[date 3 times
			vs, err := vsClient.Read(defaults.GlooSystem, initialVs.Metadata.Name, clients.ReadOpts{})
			Expect(err).To(BeNil())
			_, err = vsClient.Write(vs, clients.WriteOpts{OverwriteExisting: true})
			Expect(err).To(BeNil())

			vs, err = vsClient.Read(defaults.GlooSystem, initialVs.Metadata.Name, clients.ReadOpts{})
			Expect(err).To(BeNil())
			_, err = vsClient.Write(vs, clients.WriteOpts{OverwriteExisting: true})
			Expect(err).To(BeNil())

			vs, err = vsClient.Read(defaults.GlooSystem, initialVs.Metadata.Name, clients.ReadOpts{})
			Expect(err).To(BeNil())
			_, err = vsClient.Write(vs, clients.WriteOpts{OverwriteExisting: true})
			Expect(err).To(BeNil())

			// Get the changeset again
			cs, err = csClient.Read(defaults.GlooSystem, initialChangeset.Metadata.Name, clients.ReadOpts{})
			Expect(err).To(BeNil())
			finalEditCount := cs.EditCount.Value

			Expect(finalEditCount).To(BeEquivalentTo(initEditCount + 3))
		})
	})

	Describe("watching a resource inside the changeset", func() {

		It("notifies us when a change has occurred", func() {

			vsClient, err := v1.NewVirtualServiceClientWithToken(
				&ChangesetResourceClientFactory{ChangesetClient: csClient},
				changesetName,
			)
			Expect(err).To(BeNil())

			vsListChan, errChan, err := vsClient.Watch(defaults.GlooSystem, clients.WatchOpts{RefreshRate: 10 * time.Millisecond})
			Expect(err).To(BeNil())

			// Update the virtual service in another goroutine
			go func(vsClient v1.VirtualServiceClient) {

				vs, err := vsClient.Read(defaults.GlooSystem, initialVs.Metadata.Name, clients.ReadOpts{})
				Expect(err).To(BeNil())

				vs.VirtualHost.Name = "a-new-name"
				_, err = vsClient.Write(vs, clients.WriteOpts{OverwriteExisting: true})
				Expect(err).To(BeNil())

			}(vsClient)

			// Listen for the message notifying us about the update
			var success, skippedInitialRead bool

		DONE:
			for {
				select {
				case vsList := <-vsListChan:
					vs, err := vsList.Find(defaults.GlooSystem, initialVs.Metadata.Name)
					Expect(err).To(BeNil())

					// The first message on the channel should be an initial read. If that's the case, skip it
					if !skippedInitialRead && vs.Metadata.ResourceVersion == "1" {
						skippedInitialRead = true
						continue
					}

					Expect(vs.VirtualHost.Name).To(BeEquivalentTo("a-new-name"))
					success = true
					break DONE

				case chanErr := <-errChan:
					err = chanErr
					break DONE

				case <-time.After(1 * time.Second):
					fmt.Println("timeout")
					break DONE
				}
			}

			Expect(err).To(BeNil())
			Expect(success).To(BeTrue())
		})
	})
})
