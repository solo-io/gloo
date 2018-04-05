package grpc

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

//go:generate 2gobytes -p grpc -a annotationsDescriptorBytes -i google/api/annotations.proto.descriptor  -o annotations.google.descriptor.go
//go:generate 2gobytes -p grpc -a httpDescriptorBytes -i google/api/http.proto.descriptor -o http.google.descriptor.go
//go:generate 2gobytes -p grpc -a descriptorsDescriptorBytes -i google/api/descriptors.proto.descriptor -o descriptors.google.descriptor.go

var annotationsDescriptor, httpDescriptor, descriptorsDescriptor descriptor.FileDescriptorProto

func init() {
	err := proto.Unmarshal(annotationsDescriptorBytes, &annotationsDescriptor)
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(httpDescriptorBytes, &httpDescriptor)
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(descriptorsDescriptorBytes, &descriptorsDescriptor)
	if err != nil {
		panic(err)
	}
}

func addGoogleApisDescriptor(set *descriptor.FileDescriptorSet) {
	set.File = append([]*descriptor.FileDescriptorProto{&descriptorsDescriptor}, set.File...)
}

func addGoogleApisHttp(set *descriptor.FileDescriptorSet) {
	set.File = append([]*descriptor.FileDescriptorProto{&httpDescriptor}, set.File...)
}

func addGoogleApisAnnotations(packageName string, set *descriptor.FileDescriptorSet) {
	for _, file := range set.File {
		if *file.Package == packageName {
			file.Dependency = append(file.Dependency, "google/api/annotations.proto")
		}
	}
	set.File = append([]*descriptor.FileDescriptorProto{&annotationsDescriptor}, set.File...)
}
