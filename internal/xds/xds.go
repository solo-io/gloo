package xds

import (
	"fmt"
	"net"

	envoyapi "github.com/envoyproxy/go-control-plane/api"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/solo-io/glue/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// For now we're only running one envoy instance
const NodeKey = envoycache.Key("glue-envoy")

type hasher struct{}

func (h hasher) Hash(node *envoyapi.Node) (envoycache.Key, error) {
	return NodeKey, nil
}

func RunXDS(port int) (envoycache.Cache, *grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %v", err)
	}
	envoyCache := envoycache.NewSimpleCache(hasher{}, func(key envoycache.Key) {
		log.Debugf("CACHE: Key Updated: %s", key)
	})
	opts := []grpc_zap.Option{
		grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
			// by default everything will be logged
			return true
		}),
	}
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.NewNop(), opts...),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				log.Debugf("STREAM: info: %v\n", info)
				log.Debugf("STREAM: serverStream: %v\n", ss)
				return handler(srv, ss)
			},
		)),
	)
	xdsServer := xds.NewServer(envoyCache)
	envoyapi.RegisterAggregatedDiscoveryServiceServer(grpcServer, xdsServer)
	envoyapi.RegisterEndpointDiscoveryServiceServer(grpcServer, xdsServer)
	envoyapi.RegisterClusterDiscoveryServiceServer(grpcServer, xdsServer)
	envoyapi.RegisterRouteDiscoveryServiceServer(grpcServer, xdsServer)
	envoyapi.RegisterListenerDiscoveryServiceServer(grpcServer, xdsServer)

	go func() {
		log.Debugf("xDS server listening on %d", port)
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
		}
	}()
	return envoyCache, grpcServer, nil
}
