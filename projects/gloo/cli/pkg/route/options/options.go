package options

type Options struct {
	// for storing user input
	GrpcDestinationSpec InputGrpcDestinationSpec
}

// for populating the grpc service spec
type InputGrpcDestinationSpec struct {
	// choose the GrpcService
	Package  string
	Service  string
	Function string
}
