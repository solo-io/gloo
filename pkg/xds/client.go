package xds

import (
	"context"
	"math"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

// GetXdsClientConnection returns a gRPC connection to an xDS server
// or an error if a connection could not be established
func GetXdsClientConnection(ctx context.Context, xdsServerAddress string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		xdsServerAddress,
		// We are using non secure gRPC to the xDS server (Gloo) with the assumption that it will be secured by envoy.
		// If this assumption is not correct this needs to change.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// We block to ensure that a connection is established before progressing
		grpc.WithBlock(),
		// We only expect the received messages to be large, as those contain the set of configuration for the service
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
}
