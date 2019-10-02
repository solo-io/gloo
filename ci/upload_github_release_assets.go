package main

import "github.com/solo-io/go-utils/githubutils"

func main() {
	assets := []githubutils.ReleaseAssetSpec{
		{
			Name:       "glooe-release.yaml",
			ParentPath: "install/manifest",
		},
		{
			Name:       "gloo-with-read-only-ui-release.yaml",
			ParentPath: "install/manifest",
		},
	}
	spec := githubutils.UploadReleaseAssetSpec{
		Owner:             "solo-io",
		Repo:              "solo-projects",
		Assets:            assets,
		SkipAlreadyExists: true,
	}
	githubutils.UploadReleaseAssetCli(&spec)
}
