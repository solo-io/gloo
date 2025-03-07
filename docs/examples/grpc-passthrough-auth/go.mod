module github.com/envoyproxy/envoy/examples/ext_authz/auth/grpc-service

go 1.24

require (
	github.com/envoyproxy/go-control-plane/envoy v1.32.3
	github.com/golang/protobuf v1.5.4
	github.com/solo-io/gloo v1.18.10
	google.golang.org/genproto v0.0.0-20250224174004-546df14abb99
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250219182151-9fdb1cabc7b2
	google.golang.org/grpc v1.70.0
)

require (
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240409071808-615f978279ca // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
