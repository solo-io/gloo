package upstreamsvc_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/aws"
	mock_gloo "github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_settings "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings/mocks"
	. "github.com/solo-io/solo-projects/projects/grpcserver/server/service/internal/testutils"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc"
	mock_converter "github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/converter/mocks"
	"google.golang.org/grpc"
)

var (
	grpcServer     *grpc.Server
	conn           *grpc.ClientConn
	apiserver      v1.UpstreamApiServer
	client         v1.UpstreamApiClient
	mockCtrl       *gomock.Controller
	upstreamClient *mock_gloo.MockUpstreamClient
	inputConverter *mock_converter.MockUpstreamInputConverter
	settingsValues *mock_settings.MockValuesClient
	testErr        = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		upstreamClient = mock_gloo.NewMockUpstreamClient(mockCtrl)
		inputConverter = mock_converter.NewMockUpstreamInputConverter(mockCtrl)
		settingsValues = mock_settings.NewMockValuesClient(mockCtrl)
		apiserver = upstreamsvc.NewUpstreamGrpcService(context.TODO(), upstreamClient, inputConverter, settingsValues)

		grpcServer, conn = MustRunGrpcServer(func(s *grpc.Server) { v1.RegisterUpstreamApiServer(s, apiserver) })
		client = v1.NewUpstreamApiClient(conn)
	})

	AfterEach(func() {
		grpcServer.Stop()
		mockCtrl.Finish()
	})

	Describe("GetUpstream", func() {
		It("works when the upstream client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			upstream := gloov1.Upstream{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: metadata,
			}

			upstreamClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&upstream, nil)

			request := &v1.GetUpstreamRequest{Ref: &ref}
			actual, err := client.GetUpstream(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetUpstreamResponse{Upstream: &upstream}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstreamClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.GetUpstreamRequest{Ref: &ref}
			_, err := client.GetUpstream(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToReadUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListUpstreams", func() {
		It("works when the upstream client works", func() {
			ns1, ns2 := "one", "two"
			upstream1 := gloov1.Upstream{
				Status:   core.Status{State: core.Status_Accepted},
				Metadata: core.Metadata{Namespace: ns1},
			}
			upstream2 := gloov1.Upstream{
				Status:   core.Status{State: core.Status_Pending},
				Metadata: core.Metadata{Namespace: ns2},
			}

			upstreamClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Upstream{&upstream1}, nil)
			upstreamClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gloov1.Upstream{&upstream2}, nil)

			request := &v1.ListUpstreamsRequest{Namespaces: []string{ns1, ns2}}
			actual, err := client.ListUpstreams(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListUpstreamsResponse{Upstreams: []*gloov1.Upstream{&upstream1, &upstream2}}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			ns := "ns"

			upstreamClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.ListUpstreamsRequest{Namespaces: []string{ns}}
			_, err := client.ListUpstreams(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToListUpstreamsError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("StreamUpstreamList", func() {
		It("works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstreamList := []*gloov1.Upstream{{
				Metadata: metadata,
				UpstreamSpec: &gloov1.UpstreamSpec{
					UpstreamType: &gloov1.UpstreamSpec_Aws{
						Aws: &aws.UpstreamSpec{Region: "test"},
					},
				}},
				{
					Metadata: metadata,
					UpstreamSpec: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Aws{
							Aws: &aws.UpstreamSpec{Region: "test2"},
						},
					},
				},
			}

			refreshRate := time.Minute
			request := v1.StreamUpstreamListRequest{
				Namespace: ref.GetNamespace(),
			}
			upstreamChan := make(chan gloov1.UpstreamList, 1)
			upstreamChan <- upstreamList
			errChan := make(chan error)

			defer func() {
				close(upstreamChan)
				close(errChan)
			}()

			settingsValues.EXPECT().GetRefreshRate().Return(refreshRate)
			upstreamClient.EXPECT().
				Watch(ref.GetNamespace(), gomock.Any()).
				Return(upstreamChan, errChan, nil)

			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			streamClient, err := client.StreamUpstreamList(ctx, &request)
			Expect(err).NotTo(HaveOccurred())

			wait := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				defer func() {
					close(wait)
				}()

				actual, err := streamClient.Recv()
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.StreamUpstreamListResponse{Upstreams: upstreamList}
				ExpectEqualProtoMessages(actual, expected)

				errChan <- testErr
				_, err = streamClient.Recv()
				Expect(err).To(HaveOccurred())
				expectedErr := upstreamsvc.ErrorWhileWatchingUpstreams(testErr, metadata.Namespace)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			}()

			select {
			case <-wait:
			case <-time.After(time.Second):
				Fail("expected wait to be closed before 1s")
			}
		})
	})

	Describe("CreateUpstream", func() {
		getInput := func(ref *core.ResourceRef) *v1.UpstreamInput {
			return &v1.UpstreamInput{
				Ref: ref,
				Spec: &v1.UpstreamInput_Aws{
					Aws: &aws.UpstreamSpec{Region: "test"},
				},
			}
		}

		It("works when the upstream client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstream := gloov1.Upstream{
				Metadata: metadata,
				UpstreamSpec: &gloov1.UpstreamSpec{
					UpstreamType: &gloov1.UpstreamSpec_Aws{
						Aws: &aws.UpstreamSpec{Region: "test"},
					},
				},
			}

			inputConverter.EXPECT().
				ConvertInputToUpstreamSpec(getInput(&ref)).
				Return(upstream.UpstreamSpec)
			upstreamClient.EXPECT().
				Write(&upstream, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
				Return(&upstream, nil)

			actual, err := client.CreateUpstream(context.TODO(), &v1.CreateUpstreamRequest{Input: getInput(&ref)})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.CreateUpstreamResponse{Upstream: &upstream}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			upstream := gloov1.Upstream{
				Metadata: metadata,
				UpstreamSpec: &gloov1.UpstreamSpec{
					UpstreamType: &gloov1.UpstreamSpec_Aws{
						Aws: &aws.UpstreamSpec{Region: "test"},
					},
				},
			}

			inputConverter.EXPECT().
				ConvertInputToUpstreamSpec(getInput(&ref)).
				Return(upstream.UpstreamSpec)
			upstreamClient.EXPECT().
				Write(&upstream, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: false}).
				Return(nil, testErr)

			request := &v1.CreateUpstreamRequest{
				Input: getInput(&ref),
			}
			_, err := client.CreateUpstream(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToCreateUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateUpstream", func() {
		getInput := func(ref *core.ResourceRef) *v1.UpstreamInput {
			return &v1.UpstreamInput{
				Ref: ref,
				Spec: &v1.UpstreamInput_Aws{
					Aws: &aws.UpstreamSpec{Region: "test"},
				},
			}
		}

		It("works when the upstream client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstream := gloov1.Upstream{
				Metadata: metadata,
				UpstreamSpec: &gloov1.UpstreamSpec{
					UpstreamType: &gloov1.UpstreamSpec_Aws{
						Aws: &aws.UpstreamSpec{Region: "test"},
					},
				},
			}

			inputConverter.EXPECT().
				ConvertInputToUpstreamSpec(getInput(&ref)).
				Return(upstream.UpstreamSpec)
			upstreamClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&upstream, nil)
			upstreamClient.EXPECT().
				Write(&upstream, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(&upstream, nil)

			actual, err := client.UpdateUpstream(context.TODO(), &v1.UpdateUpstreamRequest{Input: getInput(&ref)})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.UpdateUpstreamResponse{Upstream: &upstream}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors on read", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstreamClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			_, err := client.UpdateUpstream(context.TODO(), &v1.UpdateUpstreamRequest{Input: getInput(&ref)})
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToUpdateUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the upstream client errors on write", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			upstream := gloov1.Upstream{
				Metadata: metadata,
				UpstreamSpec: &gloov1.UpstreamSpec{
					UpstreamType: &gloov1.UpstreamSpec_Aws{
						Aws: &aws.UpstreamSpec{Region: "test"},
					},
				},
			}

			inputConverter.EXPECT().
				ConvertInputToUpstreamSpec(getInput(&ref)).
				Return(upstream.UpstreamSpec)
			upstreamClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(&upstream, nil)
			upstreamClient.EXPECT().
				Write(&upstream, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			_, err := client.UpdateUpstream(context.TODO(), &v1.UpdateUpstreamRequest{Input: getInput(&ref)})
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToUpdateUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("DeleteUpstream", func() {
		It("works when the upstream client works", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			upstreamClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(nil)

			request := &v1.DeleteUpstreamRequest{Ref: &ref}
			actual, err := client.DeleteUpstream(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.DeleteUpstreamResponse{}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the upstream client errors", func() {
			ref := core.ResourceRef{
				Namespace: "ns",
				Name:      "name",
			}

			upstreamClient.EXPECT().
				Delete(ref.Namespace, ref.Name, clients.DeleteOpts{Ctx: context.TODO()}).
				Return(testErr)

			request := &v1.DeleteUpstreamRequest{Ref: &ref}
			_, err := client.DeleteUpstream(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := upstreamsvc.FailedToDeleteUpstreamError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
