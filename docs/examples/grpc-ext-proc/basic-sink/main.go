/*
First pass based on https://github.com/salrashid123/envoy_ext_proc

*/

package main

import (
	"context"
	"flag"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/solo-io/go-utils/log"

	// core "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"

	pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	structpb "google.golang.org/protobuf/types/known/structpb"
	// v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/ext_proc/v3"
	// pb "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/service/ext_proc/v3"
)

var (
	grpcport = flag.String("grpcport", ":18080", "grpcport")
	hs       *health.Server
)

const ()

type server struct{}

type healthServer struct{}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("Handling grpc Check request + %s", in.String())
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func (s *server) Process(srv pb.ExternalProcessor_ProcessServer) error {
	log.Printf("Handling grpc stream open ctx: %v\n", srv.Context())
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			log.Printf("context done\n")
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			// envoy has closed the stream. Dont return anything and close this stream entirely
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive stream request: %v", err)
		}

		// clean the response struct
		resp := &pb.ProcessingResponse{}
		log.Printf("%v", req.Request)
		log.Printf("%v", req.GetMetadataContext().FilterMetadata)
		switch v := req.Request.(type) {
		case *pb.ProcessingRequest_RequestHeaders:
			log.Printf("pb.ProcessingRequest_RequestHeaders %v \n", v)

			// respHeaders := &pb.HeadersResponse{Response: &pb.CommonResponse{HeaderMutation: mutates}}

			resp = &pb.ProcessingResponse{
				DynamicMetadata: &_struct.Struct{Fields: map[string]*structpb.Value{
					"metaKey": &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "hello from meta"}},
				}},
				Response: &pb.ProcessingResponse_RequestHeaders{},
			}

			// instructionString := ""
			// for _, n := range h.RequestHeaders.Headers.Headers {
			// 	if n.Key == "instructions" {
			// 		instructionString = n.Value
			// 	}
			// }

			// if instructionString == "skip" {
			// 	log.Printf("  Skipping request \n")
			// 	return nil
			// }

			// instructionsOnHeaders := map[string]string{}
			// err := json.Unmarshal([]byte(instructionString), &instructionsOnHeaders)
			// if err == nil {

			// 	mutates := &pb.HeaderMutation{}
			// 	log.Printf("  Instructions: %v \n", instructionsOnHeaders)
			// 	for k, v := range instructionsOnHeaders {
			// 		if strings.HasPrefix(k, "scrub-header") {
			// 			mutates.RemoveHeaders = append(mutates.RemoveHeaders, v)
			// 		}
			// 	}
			// 	respHeaders := &pb.HeadersResponse{Response: &pb.CommonResponse{HeaderMutation: mutates}}

			// 	resp = &pb.ProcessingResponse{
			// 		Response: &pb.ProcessingResponse_RequestHeaders{
			// 			RequestHeaders: respHeaders,
			// 		},
			// 		ModeOverride: &v3.ProcessingMode{
			// 			RequestBodyMode:    v3.ProcessingMode_BUFFERED,
			// 			ResponseHeaderMode: v3.ProcessingMode_SKIP,
			// 			ResponseBodyMode:   v3.ProcessingMode_NONE,
			// 		},
			// 	}
			// }
		// case *pb.ProcessingRequest_RequestBody:

		// 	r := req.Request
		// 	b := r.(*pb.ProcessingRequest_RequestBody)
		// 	log.Printf("   RequestBody: %s", string(b.RequestBody.Body))
		// 	log.Printf("   EndOfStream: %T", b.RequestBody.EndOfStream)
		// 	if b.RequestBody.EndOfStream {

		// 		bytesToSend := append(b.RequestBody.Body, []byte(` baaar `)...)
		// 		resp = &pb.ProcessingResponse{
		// 			Response: &pb.ProcessingResponse_RequestBody{
		// 				RequestBody: &pb.BodyResponse{
		// 					Response: &pb.CommonResponse{
		// 						HeaderMutation: &pb.HeaderMutation{
		// 							SetHeaders: []*core.HeaderValueOption{
		// 								{
		// 									Header: &core.HeaderValue{
		// 										Key:   "Content-Length",
		// 										Value: strconv.Itoa(len(bytesToSend)),
		// 									},
		// 								},
		// 							},
		// 						},
		// 						BodyMutation: &pb.BodyMutation{
		// 							Mutation: &pb.BodyMutation_Body{
		// 								Body: bytesToSend,
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 			ModeOverride: &v3.ProcessingMode{
		// 				ResponseHeaderMode: v3.ProcessingMode_SEND,
		// 				ResponseBodyMode:   v3.ProcessingMode_NONE,
		// 			},
		// 		}
		// 	}
		// 	break
		// case *pb.ProcessingRequest_ResponseHeaders:
		// 	log.Printf("pb.ProcessingRequest_ResponseHeaders %v \n", v)
		// 	r := req.Request
		// 	h := r.(*pb.ProcessingRequest_ResponseHeaders)

		// 	responseSize := 0
		// 	for _, n := range h.ResponseHeaders.Headers.Headers {
		// 		if n.Key == "content-length" {
		// 			responseSize, _ = strconv.Atoi(n.Value)
		// 			break
		// 		}
		// 	}

		// 	log.Printf("  Removing access-control-allow-* headers \n")
		// 	rhq := &pb.HeadersResponse{
		// 		Response: &pb.CommonResponse{
		// 			HeaderMutation: &pb.HeaderMutation{
		// 				RemoveHeaders: []string{"access-control-allow-origin", "access-control-allow-credentials"},
		// 				SetHeaders: []*core.HeaderValueOption{
		// 					{
		// 						Header: &core.HeaderValue{
		// 							Key:   "content-type",
		// 							Value: "text/plain",
		// 						},
		// 					},
		// 					{
		// 						Header: &core.HeaderValue{
		// 							Key:   "content-length",
		// 							Value: strconv.Itoa(responseSize + len([]byte(` qux`))),
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	}
		// 	resp = &pb.ProcessingResponse{
		// 		Response: &pb.ProcessingResponse_ResponseHeaders{
		// 			ResponseHeaders: rhq,
		// 		},
		// 		ModeOverride: &v3.ProcessingMode{
		// 			ResponseBodyMode: v3.ProcessingMode_BUFFERED,
		// 		},
		// 	}
		// 	break
		// case *pb.ProcessingRequest_ResponseBody:
		// 	log.Printf("pb.ProcessingRequest_ResponseBody %v \n", v)
		// 	r := req.Request
		// 	b := r.(*pb.ProcessingRequest_ResponseBody)
		// 	if b.ResponseBody.EndOfStream {
		// 		bytesToSend := append(b.ResponseBody.Body, []byte(` qux`)...)
		// 		resp = &pb.ProcessingResponse{
		// 			Response: &pb.ProcessingResponse_ResponseBody{
		// 				ResponseBody: &pb.BodyResponse{
		// 					Response: &pb.CommonResponse{
		// 						BodyMutation: &pb.BodyMutation{
		// 							Mutation: &pb.BodyMutation_Body{
		// 								Body: bytesToSend,
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		}
		// 	}

		// 	break
		default:
			log.Printf("Unknown Request type %v\n", v)
		}

		// At this point we believe we have created a valid response...
		// note that this is sometimes not the case
		// anyways for now just send it
		if err := srv.Send(resp); err != nil {
			log.Printf("send error %v", err)
		}
	}
}

func main() {

	flag.Parse()

	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(1000)}
	s := grpc.NewServer(sopts...)

	pb.RegisterExternalProcessorServer(s, &server{})

	healthpb.RegisterHealthServer(s, &healthServer{})

	log.Printf("Starting gRPC server on port %s\n", *grpcport)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		time.Sleep(time.Second)
		log.Printf("Graceful stop completed")
		os.Exit(0)
	}()
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("killing server with %v", err)
	}
}
