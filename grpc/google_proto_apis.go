package grpc

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

//go:generate 2gobytes -p grpc -a annotationsDescriptorBytes -i google/api/annotations.proto.descriptor  -o annotations.google.descriptor.go
//go:generate 2gobytes -p grpc -a httpDescriptorBytes -i google/api/http.proto.descriptor -o http.google.descriptor.go

var annotationsDescriptor, httpDescriptor descriptor.FileDescriptorProto

func init() {
	err := proto.Unmarshal(annotationsDescriptorBytes, &annotationsDescriptor)
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(httpDescriptorBytes, &httpDescriptor)
	if err != nil {
		panic(err)
	}
}

func addGoogleApisHttp(set *descriptor.FileDescriptorSet) {
	set.File = append([]*descriptor.FileDescriptorProto{&httpDescriptor}, set.File...)
}

func addGoogleApisAnnotations(set *descriptor.FileDescriptorSet) {
	set.File = append([]*descriptor.FileDescriptorProto{&annotationsDescriptor}, set.File...)
}
