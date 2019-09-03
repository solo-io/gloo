package main

import (
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/solo-projects/test/extauth/plugins/is_header_present/plugin"
)

//go:generate go build -buildmode=plugin -o ./../IsHeaderPresent.so interface.go

func main() {}

var _ api.ExtAuthPlugin = &plugin.IsHeaderPresentPlugin{}

// This is the exported symbol that the ext auth server will look for.
// In this case we export the struct, not a pointer to it. The resulting symbol will thus be of type *plugin.IsHeaderPresentPlugin,
// which implements the ExtAuthPlugin interface.
//noinspection GoUnusedGlobalVariable
var IsHeaderPresent plugin.IsHeaderPresentPlugin
