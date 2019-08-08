package nackdetector_test

import (
	"context"
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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

var _ = Describe("NackDetector", func() {

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
			cb := StateChangedCallback(func(id EnvoyStatusId, olds, s State) { stateChan <- s })
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
			envoyStates.CheckIsSync(envoyId, "1", "1")

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
			envoyStates.CheckIsSync(envoyId, "1", "2")

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
			envoyStates.CheckIsSync(envoyId, "1", "1")

			// make time fly:
			newtime := tp.NowTime.Add(envoyStates.WaitTimeForSync).Add(time.Microsecond)
			tp.NowTime = newtime
			tp.AfterChan <- newtime
			Eventually(stateChan).Should(Receive(Equal(OutOfSync)))
		})

		Context("NackDetector", func() {

			It("Should work on happy path", func() {

				tp := newFakeTimeProvider()
				envoyStates.TimeProvider = tp

				nd := NewNackDetectorWithEnvoysState(ctx, envoyStates)

				nd.OnStreamOpen(1, "")
				// envoy saying hello
				req := &v2.DiscoveryRequest{VersionInfo: "", ResponseNonce: "", TypeUrl: "type.googleapis.com/envoy.api.v2.Clusters", Node: &core.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)

				req = &v2.DiscoveryRequest{VersionInfo: "", ResponseNonce: "", TypeUrl: "type.googleapis.com/envoy.api.v2.Listeners", Node: &core.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)

				resp := &v2.DiscoveryResponse{
					Nonce: "1", VersionInfo: "1",
					TypeUrl: "type.googleapis.com/envoy.api.v2.Clusters",
				}
				nd.OnStreamResponse(1, req, resp)

				resp = &v2.DiscoveryResponse{
					Nonce: "2", VersionInfo: "1",
					TypeUrl: "type.googleapis.com/envoy.api.v2.Listeners",
				}
				nd.OnStreamResponse(1, req, resp)

				req = &v2.DiscoveryRequest{VersionInfo: "1", ResponseNonce: "1", TypeUrl: "type.googleapis.com/envoy.api.v2.Clusters", Node: &core.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)

				req = &v2.DiscoveryRequest{VersionInfo: "1", ResponseNonce: "2", TypeUrl: "type.googleapis.com/envoy.api.v2.Listeners", Node: &core.Node{Id: "1"}}
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

				nd := NewNackDetectorWithEnvoysState(ctx, envoyStates)

				nd.OnStreamOpen(1, "")
				// envoy saying hello
				req := &v2.DiscoveryRequest{VersionInfo: "", ResponseNonce: "", TypeUrl: "type.googleapis.com/envoy.api.v2.Clusters", Node: &core.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)

				resp := &v2.DiscoveryResponse{
					Nonce: "1", VersionInfo: "1",
					TypeUrl: "type.googleapis.com/envoy.api.v2.Clusters",
				}
				nd.OnStreamResponse(1, req, resp)

				req = &v2.DiscoveryRequest{VersionInfo: "1", ResponseNonce: "1", TypeUrl: "type.googleapis.com/envoy.api.v2.Clusters", Node: &core.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)

				// now simulate an update
				resp = &v2.DiscoveryResponse{
					Nonce: "2", VersionInfo: "2",
					TypeUrl: "type.googleapis.com/envoy.api.v2.Clusters",
				}
				nd.OnStreamResponse(1, req, resp)

				req = &v2.DiscoveryRequest{VersionInfo: "2", ResponseNonce: "2", TypeUrl: "type.googleapis.com/envoy.api.v2.Clusters", Node: &core.Node{Id: "1"}}
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

				nd := NewNackDetectorWithEnvoysState(ctx, envoyStates)

				nd.OnStreamOpen(1, "ads")
				// envoy saying hello
				req := &v2.DiscoveryRequest{VersionInfo: "", ResponseNonce: "", Node: &core.Node{Id: "1"}}
				nd.OnStreamRequest(1, req)
				resp := &v2.DiscoveryResponse{
					Nonce: "1", VersionInfo: "1",
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
