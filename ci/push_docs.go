package main

import (
	"github.com/solo-io/go-utils/docsutils"
)

func main() {
	spec := docsutils.DocsPRSpec{
		Owner: "solo-io",
		Repo: "solo-projects",
		Product: "gloo",
		ChangelogPrefix: "glooe",
		DocsParentPath: "projects/gloo/doc",
		ApiPaths: []string {
			"docs/v1/github.com/solo-io/solo-projects",
		},
		CliPrefix: "glooctl",
	}
	docsutils.PushDocsCli(&spec)
}