package main

import (
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/solo-projects/test/extauth/plugins/required_header/plugin"
)

//go:generate go build -buildmode=plugin -o ./../RequiredHeader.so interface.go

func main() {}

var _ api.ExtAuthPlugin = &plugin.RequiredHeaderPlugin{}

// This is the exported symbol that the ext auth server will look for.
// In this case we export the struct, not a pointer to it. The resulting symbol will thus be of type *plugin.RequiredHeaderPlugin,
// which implements the ExtAuthPlugin interface.
//noinspection GoUnusedGlobalVariable
var RequiredHeader plugin.RequiredHeaderPlugin
