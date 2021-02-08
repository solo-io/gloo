package options

import "context"

type Options struct {
	Ctx           context.Context
	Namespace     string
	ApiserverPort string
}
