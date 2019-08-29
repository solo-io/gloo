package pkg

import "github.com/go-ldap/ldap"

//go:generate mockgen -destination mocks/client_builder_mock.go -package mocks github.com/solo-io/solo-projects/projects/extauth/plugins/ldap/pkg ClientBuilder
//go:generate mockgen -destination mocks/client_mock.go -package mocks github.com/go-ldap/ldap Client

// Used to be able to generate mocks of the `Dial` client constructor function
type ClientBuilder interface {
	Dial(network, addr string) (ldap.Client, error)
}

func NewLdapClientBuilder() *ldapClientBuilder {
	return &ldapClientBuilder{}
}

type ldapClientBuilder struct{}

// Just delegates
func (b *ldapClientBuilder) Dial(network, addr string) (ldap.Client, error) {
	return ldap.Dial(network, addr)
}
