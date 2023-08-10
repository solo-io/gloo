module github.com/solo-io/ext-proc-examples/basic-sink

go 1.20

require (
	github.com/envoyproxy/go-control-plane v0.11.1
	github.com/solo-io/gloo v1.15.0-rc2
	github.com/solo-io/go-utils v0.24.6
	google.golang.org/grpc v1.57.0
)

replace github.com/solo-io/gloo => ./../../../..
