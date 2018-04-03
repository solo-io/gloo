package consul_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/consul/api"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
	. "github.com/solo-io/gloo-storage/consul"
	"github.com/solo-io/gloo-testing/helpers"
)

var _ = Describe("ConsulStorageClient", func() {
	var rootPath string
	var consul *api.Client
	BeforeEach(func() {
		rootPath = helpers.RandString(4)
		c, err := api.NewClient(api.DefaultConfig())
		Expect(err).NotTo(HaveOccurred())
		consul = c
	})
	AfterEach(func() {
		consul.KV().DeleteTree(rootPath, nil)
	})
	Describe("Upstreams", func() {
		Describe("create", func() {
			It("creates the upstream as a consul key", func() {
				client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
				Expect(err).NotTo(HaveOccurred())
				input := &v1.Upstream{
					Name:              "myupstream",
					Type:              "foo",
					ConnectionTimeout: time.Second,
				}
				us, err := client.V1().Upstreams().Create(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(us).NotTo(Equal(input))
				p, _, err := consul.KV().Get(rootPath+"/upstreams/"+input.Name, nil)
				Expect(err).NotTo(HaveOccurred())
				var unmarshalledUpstream v1.Upstream
				err = proto.Unmarshal(p.Value, &unmarshalledUpstream)
				Expect(err).NotTo(HaveOccurred())
				Expect(&unmarshalledUpstream).To(Equal(input))
				resourceVersion := fmt.Sprintf("%v", p.CreateIndex)
				Expect(us.Metadata.ResourceVersion).To(Equal(resourceVersion))
				input.Metadata = us.Metadata
				Expect(us).To(Equal(input))
			})
			It("errors when creating the same upstream twice", func() {
				client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
				Expect(err).NotTo(HaveOccurred())
				input := &v1.Upstream{
					Name:              "myupstream",
					Type:              "foo",
					ConnectionTimeout: time.Second,
				}
				_, err = client.V1().Upstreams().Create(input)
				Expect(err).NotTo(HaveOccurred())
				_, err = client.V1().Upstreams().Create(input)
				Expect(err).To(HaveOccurred())
			})
			Describe("update", func() {
				It("fails if the upstream doesn't exist", func() {
					client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
					input := &v1.Upstream{
						Name:              "myupstream",
						Type:              "foo",
						ConnectionTimeout: time.Second,
					}
					us, err := client.V1().Upstreams().Update(input)
					Expect(err).To(HaveOccurred())
					Expect(us).To(BeNil())
				})
				It("fails if the resourceversion is not up to date", func() {
					client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
					input := &v1.Upstream{
						Name:              "myupstream",
						Type:              "foo",
						ConnectionTimeout: time.Second,
					}
					_, err = client.V1().Upstreams().Create(input)
					Expect(err).NotTo(HaveOccurred())
					v, err := client.V1().Upstreams().Update(input)
					Expect(v).To(BeNil())
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("resource version"))
				})
				It("updates the upstream", func() {
					client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
					input := &v1.Upstream{
						Name:              "myupstream",
						Type:              "foo",
						ConnectionTimeout: time.Second,
					}
					us, err := client.V1().Upstreams().Create(input)
					Expect(err).NotTo(HaveOccurred())
					changed := proto.Clone(input).(*v1.Upstream)
					changed.Type = "bar"
					// match resource version
					changed.Metadata = us.Metadata
					out, err := client.V1().Upstreams().Update(changed)
					Expect(err).NotTo(HaveOccurred())
					Expect(out.Type).To(Equal(changed.Type))
				})
				Describe("get", func() {
					It("fails if the upstream doesn't exist", func() {
						client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
						Expect(err).NotTo(HaveOccurred())
						us, err := client.V1().Upstreams().Get("foo")
						Expect(err).To(HaveOccurred())
						Expect(us).To(BeNil())
					})
					It("returns the upstream", func() {
						client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
						Expect(err).NotTo(HaveOccurred())
						input := &v1.Upstream{
							Name:              "myupstream",
							Type:              "foo",
							ConnectionTimeout: time.Second,
						}
						us, err := client.V1().Upstreams().Create(input)
						Expect(err).NotTo(HaveOccurred())
						out, err := client.V1().Upstreams().Get(input.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(out).To(Equal(us))
						input.Metadata = out.Metadata
						Expect(out).To(Equal(input))
					})
				})
				Describe("list", func() {
					It("returns all existing upstreams", func() {
						client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
						Expect(err).NotTo(HaveOccurred())
						input1 := &v1.Upstream{
							Name:              "myupstream1",
							Type:              "foo",
							ConnectionTimeout: time.Second,
						}
						input2 := &v1.Upstream{
							Name:              "myupstream2",
							Type:              "foo",
							ConnectionTimeout: time.Second,
						}
						input3 := &v1.Upstream{
							Name:              "myupstream3",
							Type:              "foo",
							ConnectionTimeout: time.Second,
						}
						us1, err := client.V1().Upstreams().Create(input1)
						Expect(err).NotTo(HaveOccurred())
						us2, err := client.V1().Upstreams().Create(input2)
						Expect(err).NotTo(HaveOccurred())
						us3, err := client.V1().Upstreams().Create(input3)
						Expect(err).NotTo(HaveOccurred())
						out, err := client.V1().Upstreams().List()
						Expect(err).NotTo(HaveOccurred())
						Expect(out).To(ContainElement(us1))
						Expect(out).To(ContainElement(us2))
						Expect(out).To(ContainElement(us3))
					})
				})
				Describe("watch", func() {
					It("watches", func() {
						client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
						Expect(err).NotTo(HaveOccurred())
						lists := make(chan []*v1.Upstream, 3)
						stop := make(chan struct{})
						defer close(stop)
						errs := make(chan error)
						w, err := client.V1().Upstreams().Watch(&storage.UpstreamEventHandlerFuncs{
							UpdateFunc: func(updatedList []*v1.Upstream, _ *v1.Upstream) {
								lists <- updatedList
							},
						})
						Expect(err).NotTo(HaveOccurred())
						go func() {
							w.Run(stop, errs)
						}()
						input1 := &v1.Upstream{
							Name:              "myupstream1",
							Type:              "foo",
							ConnectionTimeout: time.Second,
						}
						input2 := &v1.Upstream{
							Name:              "myupstream2",
							Type:              "foo",
							ConnectionTimeout: time.Second,
						}
						input3 := &v1.Upstream{
							Name:              "myupstream3",
							Type:              "foo",
							ConnectionTimeout: time.Second,
						}
						us1, err := client.V1().Upstreams().Create(input1)
						Expect(err).NotTo(HaveOccurred())
						us2, err := client.V1().Upstreams().Create(input2)
						Expect(err).NotTo(HaveOccurred())
						us3, err := client.V1().Upstreams().Create(input3)
						Expect(err).NotTo(HaveOccurred())

						Eventually(lists).Should(HaveLen(3))
						list1 := <-lists
						Expect(list1).To(HaveLen(1))
						Expect(list1).To(ContainElement(us1))
						list2 := <-lists
						Expect(list2).To(HaveLen(2))
						Expect(list2).To(ContainElement(us1))
						Expect(list2).To(ContainElement(us2))
						list3 := <-lists
						Expect(list3).To(HaveLen(3))
						Expect(list3).To(ContainElement(us1))
						Expect(list3).To(ContainElement(us2))
						Expect(list3).To(ContainElement(us3))
					})
				})
			})
		})
	})
})
