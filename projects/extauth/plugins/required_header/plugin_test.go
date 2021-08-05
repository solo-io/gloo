package main

import (
	"context"
	"plugin"

	structpb "github.com/golang/protobuf/ptypes/struct"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-plugins/api"
)

// This will be called in CI before tests are run to generate the RequiredHeader.so file needed by this test.
//go:generate go build -buildmode=plugin -trimpath -o RequiredHeader.so plugin.go

var _ = Describe("Plugin", func() {

	It("can be loaded", func() {

		goPlugin, err := plugin.Open("RequiredHeader.so")
		Expect(err).NotTo(HaveOccurred())

		pluginStructPtr, err := goPlugin.Lookup("Plugin")
		Expect(err).NotTo(HaveOccurred())

		extAuthPlugin, ok := pluginStructPtr.(api.ExtAuthPlugin)
		Expect(ok).To(BeTrue())

		instance, err := extAuthPlugin.NewConfigInstance(context.TODO())
		Expect(err).NotTo(HaveOccurred())

		typedInstance, ok := instance.(*structpb.Struct)
		Expect(ok).To(BeTrue())

		Expect(typedInstance).To(Equal(&structpb.Struct{}))
	})
})
