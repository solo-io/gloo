package main_test

import (
	"context"
	"plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/solo-projects/projects/extauth/plugins/ldap/pkg"
)

// This will be called in CI before tests are run to generate the Ldap.so file needed by this test.
//go:generate go build -buildmode=plugin -o Ldap.so plugin.go

var _ = Describe("Plugin", func() {
	It("can be loaded", func() {

		goPlugin, err := plugin.Open("Ldap.so")
		Expect(err).NotTo(HaveOccurred())

		pluginSymbol, err := goPlugin.Lookup("Plugin")
		Expect(err).NotTo(HaveOccurred())

		extAuthPlugin, ok := pluginSymbol.(api.ExtAuthPlugin)
		Expect(ok).To(BeTrue())

		instance, err := extAuthPlugin.NewConfigInstance(context.TODO())
		Expect(err).NotTo(HaveOccurred())

		typedInstance, ok := instance.(*pkg.Config)
		Expect(ok).To(BeTrue())

		Expect(typedInstance.ServerUrl).To(BeEmpty())
		Expect(typedInstance.UserDnTemplate).To(BeEmpty())
		Expect(typedInstance.AllowedGroups).To(HaveLen(0))
	})
})
