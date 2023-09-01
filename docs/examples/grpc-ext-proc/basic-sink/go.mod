module github.com/solo-io/ext-proc-examples/basic-sink

go 1.20

require (
	github.com/envoyproxy/go-control-plane v0.11.1
	github.com/solo-io/go-utils v0.24.6
	google.golang.org/grpc v1.56.0-dev
)

require (
	github.com/cncf/xds/go v0.0.0-20230428030218-4003588d1b74 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/k0kubun/pp v2.3.0+incompatible // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230526203410-71b5a4ffd15e // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/solo-io/gloo => ./../../../..

replace github.com/envoyproxy/go-control-plane => github.com/solo-io/go-control-plane-fork-v2 v0.0.0-20230820180611-2cf7e0da27c8
