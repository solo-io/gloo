package main

import (
	"github.com/solo-io/ext-auth-plugin-examples/plugins/required_header/pkg"
	"github.com/solo-io/ext-auth-plugins/api"
)

func main() {}

// Compile-time assertion
var _ api.ExtAuthPlugin = &pkg.RequiredHeaderPlugin{}

// This is the exported symbol that GlooE will look for.
//noinspection GoUnusedGlobalVariable
var Plugin pkg.RequiredHeaderPlugin
