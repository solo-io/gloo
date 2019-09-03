package main

import (
	"github.com/solo-io/ext-auth-plugin-examples/plugins/required_header/pkg"
	"github.com/solo-io/ext-auth-plugins/api"
)

//go:generate go build -buildmode=plugin -o ./../RequiredHeaderValue.so interface.go

func main() {}

var _ api.ExtAuthPlugin = &pkg.RequiredHeaderPlugin{}

// This is the exported symbol that the ext auth server will look for. In this case we export a pointer to the struct,
// so a valid ExtAuthPlugin interface implementation. The resulting symbol loaded by the ext-auth server will thus be
// of type **plugin.HeaderValuePlugin. We want to test that the server can handle this case as well.
//noinspection GoUnusedGlobalVariable
var Plugin *pkg.RequiredHeaderPlugin
