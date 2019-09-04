package main

import (
	"github.com/solo-io/go-utils/docsutils"
)

func main() {
	spec := docsutils.DocsPRSpec{
		Owner:           "solo-io",
		Repo:            "solo-projects",
		Product:         "gloo",
		ChangelogPrefix: "glooe",
		DocsParentPath:  "",
		CliPrefix:       "glooctl",
	}
	docsutils.PushDocsCli(&spec)
}
