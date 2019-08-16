package main

import (
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/solo-projects/test/extauth/plugins/header_value/plugin"
)

//go:generate go build -buildmode=plugin -o ./../CheckHeaderValue.so interface.go

func main() {}

var _ api.ExtAuthPlugin = &plugin.HeaderValuePlugin{}

// This is the exported symbol that the ext auth server will look for. In this case we export a pointer to the struct,
// so a valid ExtAuthPlugin interface implementation. The resulting symbol loaded by the ext-auth server will thus be
// of type **plugin.RequiredHeaderPlugin. We want to test that the server can handle this case as well.
//noinspection GoUnusedGlobalVariable
var Plugin *plugin.HeaderValuePlugin
