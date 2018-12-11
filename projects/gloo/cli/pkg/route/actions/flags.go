package actions

import (
	ro "github.com/solo-io/gloo/projects/gloo/cli/pkg/route/options"
	"github.com/spf13/cobra"
)

// NOTE: currently unused

func addGrpcDestinationFlags(cmd *cobra.Command, o *ro.Options) {
	flags := cmd.Flags()

	flags.StringVar(&o.GrpcDestinationSpec.Package,
		"grpc.package",
		"",
		"gRPC package name")

	flags.StringVar(&o.GrpcDestinationSpec.Service,
		"grpc.service",
		"",
		"gRPC service name")

	flags.StringVar(&o.GrpcDestinationSpec.Function,
		"grpc.function",
		"",
		"gRPC function name")

}
