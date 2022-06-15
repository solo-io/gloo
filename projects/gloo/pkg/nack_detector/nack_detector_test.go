package nackdetector_test

import (
	"context"
	"time"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/types"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/nack_detector"
)

type fakeTimeProvider struct {
	NowTime   time.Time
	AfterChan chan time.Time
}

func (f *fakeTimeProvider) Now() time.Time {
	return f.NowTime
}

func (f *fakeTimeProvider) After(t time.Duration) <-chan time.Time {
	return f.AfterChan
}

func newFakeTimeProvider() *fakeTimeProvider {
	return &fakeTimeProvider{
		// use a constant data for easy understanding in test logs
		NowTime:   time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC),
		AfterChan: make(chan time.Time, 10),
	}
}

var _ = Describe("nackDetector", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
	})
	AfterEach(func() {
		cancel()
	})

	Context("EnvoysState", func() {

		var (
			envoyStates *EnvoysState
			stateChan   chan State
			es          *EnvoyState
			envoyId     DiscoveryServiceId
		)

		BeforeEach(func() {
			stateChan = make(chan State, 10)
			cb := StateChangedCallback(func(ctx context.Context, id EnvoyStatusId, olds, s State) { stateChan <- s })
			es = &EnvoyState{
				ServerVersion: "1",
				ServerNonce:   "1",
				EnvoyStatus: EnvoyStatus{
					EnvoyStatusId: EnvoyStatusId{
						NodeId: "node1",
						StreamId: DiscoveryServiceId{
							GrpcStreamId: 1,
							TypeUrl:      "url",
						},
					},
					State: InSync,
				},
			}

			envoyId.GrpcStreamId = 1
			envoyId.TypeUrl = "clusters"
			envoyStates = NewEnvoysState(ctx, cb)
		})

		initialRequest := func() {
			// set initial call
			envoyStates.Set(envoyId, es)
			// envoy responds
			envoyStates.CheckIsSync(ctx, envoyId, "1", "1")

			Eventually(stateChan).Should(Receive(Equal(InSync)))

		}

		It("should detect in sync", func() {
			initialRequest()
		})

		It("should detect in sync", func() {
			initialRequest()

			// a new config
			es.ServerVersion = "2"
			es.ServerNonce = "2"
			envoyStates.Set(envoyId, es)

			// envoy nacks: correct nonce and incorrect version
			envoyStates.CheckIsSync(ctx, envoyId, "1", "2")

			Eventually(stateChan).Should(Receive(Equal(OutOfSyncNack)))
		})

		It("should detect in out of sync no nack", func() {
			tp := newFakeTimeProvider()
			envoyStates.TimeProvider = tp

			initialRequest()

			// a new config
			es.ServerVersion = "2"
			es.ServerNonce = "2"
			envoyStates.Set(envoyId, es)

			// envoy just confused:  incorrect nonce and incorrect version
			envoyStates.CheckIsSync(ctx, envoyId, "1", "1")

			// make time fly:
			newtime := tp.NowTime.Add(envoyStates.WaitTimeForSync).Add(time.Microsecond)
			tp.NowTime = newtime
			tp.AfterChan <- newtime
			Eventually(stateChan).Should(Receive(Equal(OutOfSync)))
		})

		Context("nackDetector", func() {

			It("Should work on happy path", func() {

				tp := newFakeTimeProvider()
				envoyStates.TimeProvider = tp

				nd := NewNackDetectorWithEnvoysState(envoyStates)

				nd.OnStreamOpen(ctx, 1, "")
				// envoy saying hello
				req := &envoy_service_discovery_v3.DiscoveryRequest{
					VersionInfo:   "",
					ResponseNonce: "",
					TypeUrl:       types.ClusterTypeV3,
					Node:          &envoy_config_core_v3.Node{Id: "1"},
				}
				nd.OnStreamRequest(1, req)

				req = &envoy_service_discovery_v3.DiscoveryRequest{
					VersionInfo:   "",
					ResponseNonce: "",
					TypeUrl:       types.ListenerTypeV3,
					Node:          &envoy_config_core_v3.Node{Id: "1"},
				}
				nd.OnStreamRequest(1, req)

				resp := &envoy_service_discovery_v3.DiscoveryResponse{
					Nonce:       "1",
					VersionInfo: "1",
					TypeUrl:     types.ClusterTypeV3,
				}
				nd.OnStreamResponse(1, req, resp)

				resp = &envoy_service_discovery_v3.DiscoveryResponse{
					Nonce:       "2",
					VersionInfo: "1",
					TypeUrl:     types.ListenerTypeV3,
				}
				nd.OnStreamResponse(1, req, resp)

				req = &envoy_service_discovery_v3.DiscoveryRequest{VersionInfo: "1", ResponseNonce: "1", TypeUrl: types.ClusterTypeV3, Node: &envoy_config_core_v3.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)

				req = &envoy_service_discovery_v3.DiscoveryRequest{VersionInfo: "1", ResponseNonce: "2", TypeUrl: "type.googleapis.com/envoy.api.v2.Listeners", Node: &envoy_config_core_v3.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)

				var lastchanvalue State
				timeout := time.After(time.Second)

				var firstvalue State

				select {
				case firstvalue = <-stateChan:
				case <-timeout:
				}
			Loop:
				for {

					select {
					case lastchanvalue = <-stateChan:
					case <-timeout:
						break Loop
					}
				}

				Expect(firstvalue).To(Equal(New))
				Expect(lastchanvalue).To(Equal(InSync))
			})

			It("Should update correctly", func() {

				tp := newFakeTimeProvider()
				envoyStates.TimeProvider = tp

				nd := NewNackDetectorWithEnvoysState(envoyStates)

				nd.OnStreamOpen(ctx, 1, "")
				// envoy saying hello
				req := &envoy_service_discovery_v3.DiscoveryRequest{
					VersionInfo:   "",
					ResponseNonce: "",
					TypeUrl:       types.ClusterTypeV3,
					Node:          &envoy_config_core_v3.Node{Id: "1"},
				}
				nd.OnStreamRequest(1, req)

				resp := &envoy_service_discovery_v3.DiscoveryResponse{
					Nonce:       "1",
					VersionInfo: "1",
					TypeUrl:     types.ClusterTypeV3,
				}
				nd.OnStreamResponse(1, req, resp)

				req = &envoy_service_discovery_v3.DiscoveryRequest{
					VersionInfo:   "1",
					ResponseNonce: "1",
					TypeUrl:       types.ClusterTypeV3,
					Node:          &envoy_config_core_v3.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)

				// now simulate an update
				resp = &envoy_service_discovery_v3.DiscoveryResponse{
					Nonce:       "2",
					VersionInfo: "2",
					TypeUrl:     types.ClusterTypeV3,
				}
				nd.OnStreamResponse(1, req, resp)

				req = &envoy_service_discovery_v3.DiscoveryRequest{
					VersionInfo:   "2",
					ResponseNonce: "2",
					TypeUrl:       types.ClusterTypeV3,
					Node:          &envoy_config_core_v3.Node{Id: "1"},
				}
				nd.OnStreamRequest(1, req)

				// make time fly:
				newtime := tp.NowTime.Add(envoyStates.WaitTimeForSync).Add(time.Microsecond)
				tp.NowTime = newtime
				tp.AfterChan <- newtime

				var lastchanvalue State
				timeout := time.After(time.Second)
			Loop:
				for {

					select {
					case lastchanvalue = <-stateChan:
					case <-timeout:
						break Loop
					}
				}

				Expect(lastchanvalue).To(Equal(InSync))
			})

			It("Should catch envoys that are not responding", func() {

				tp := newFakeTimeProvider()
				envoyStates.TimeProvider = tp

				nd := NewNackDetectorWithEnvoysState(envoyStates)

				nd.OnStreamOpen(ctx, 1, "ads")
				// envoy saying hello
				req := &envoy_service_discovery_v3.DiscoveryRequest{
					VersionInfo:   "",
					ResponseNonce: "",
					Node:          &envoy_config_core_v3.Node{Id: "1"},
				}
				nd.OnStreamRequest(1, req)
				resp := &envoy_service_discovery_v3.DiscoveryResponse{
					Nonce:       "1",
					VersionInfo: "1",
				}
				nd.OnStreamResponse(1, req, resp)

				// sleep a bit as the processing is in a go-routine
				time.Sleep(time.Second)

				// make time fly:
				newtime := tp.NowTime.Add(envoyStates.WaitTimeForSync).Add(time.Microsecond)
				tp.NowTime = newtime
				tp.AfterChan <- newtime

				Eventually(stateChan).Should(Receive(Equal(New)))
				Eventually(stateChan, "5s").Should(Receive(Equal(OutOfSync)))

			})

		})

	})

})
