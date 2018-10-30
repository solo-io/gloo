package mocks

//go:generate protoc -I=./ -I=${GOPATH}/src/github.com/gogo/protobuf/  -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ -I=${GOPATH}/src --gogo_out=plugins=grpc:${GOPATH}/src --plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=./project.json:. mock_resources.proto
