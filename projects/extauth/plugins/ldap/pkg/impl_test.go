package pkg_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/extauth/plugins/ldap/pkg"
)

var _ = Describe("LDAP plugin", func() {

	var factory = &pkg.LdapFactory{}

	It("returns expected config instance", func() {
		config, err := factory.NewConfigInstance(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(config).To(BeAssignableToTypeOf(&pkg.Config{}))
	})

	Describe("config validation", func() {

		It("fails if required attributes are missing", func() {
			_, err := factory.GetAuthService(context.Background(), &pkg.Config{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(pkg.MissingRequiredConfigError([]string{"ServerUrl", "UserDnTemplate", "AllowedGroups"})))

			_, err = factory.GetAuthService(context.Background(), &pkg.Config{
				UserDnTemplate: "uid=%s,ou=people,dc=solo,dc=io",
				AllowedGroups:  []string{"cn=developers,ou=groups,dc=solo,dc=io"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(pkg.MissingRequiredConfigError([]string{"ServerUrl"})))

			_, err = factory.GetAuthService(context.Background(), &pkg.Config{
				ServerUrl: "localhost",
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(pkg.MissingRequiredConfigError([]string{"UserDnTemplate", "AllowedGroups"})))
		})

		It("fails if user DN template is malformed", func() {
			_, err := factory.GetAuthService(context.Background(), &pkg.Config{
				ServerUrl:      "localhost",
				UserDnTemplate: "ou=people,dc=solo,dc=io",
				AllowedGroups:  []string{"cn=developers,ou=groups,dc=solo,dc=io"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(pkg.MalformedTemplate(0)))

			_, err = factory.GetAuthService(context.Background(), &pkg.Config{
				ServerUrl:      "localhosta",
				UserDnTemplate: "uid=%s,cn=%s,ou=groups,dc=solo,dc=io",
				AllowedGroups:  []string{"cn=developers,ou=groups,dc=solo,dc=io"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(pkg.MalformedTemplate(2)))
		})

	})

	It("correctly returns auth service", func() {
		svc, err := factory.GetAuthService(context.Background(), &pkg.Config{
			ServerUrl:      "localhost",
			UserDnTemplate: "uid=%s,ou=people,dc=solo,dc=io",
			AllowedGroups:  []string{"cn=developers,ou=groups,dc=solo,dc=io"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())
	})
})
