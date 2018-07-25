package apiserver

//go:generate protoc -I=./ -I=${GOPATH}/src/github.com/gogo/protobuf/  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ -I=${GOPATH}/src --gogo_out=plugins=grpc,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:${GOPATH}/src --plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=. api_server.proto
