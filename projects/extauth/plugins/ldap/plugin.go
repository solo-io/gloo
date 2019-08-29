package main

import (
	"github.com/solo-io/ext-auth-plugins/api"
	impl "github.com/solo-io/solo-projects/projects/extauth/plugins/ldap/pkg"
)

func main() {}

// Compile-time assertion
var _ api.ExtAuthPlugin = impl.LdapFactory{}

// This is the exported symbol that GlooE will look for.
//noinspection GoUnusedGlobalVariable
var Plugin impl.LdapFactory
