package main

import "github.com/solo-io/go-utils/githubutils"

func main() {
	assets := make([]githubutils.ReleaseAssetSpec, 5)
	assets[0] = githubutils.ReleaseAssetSpec{
		Name:       "glooctl-linux-amd64",
		ParentPath: "_output",
		UploadSHA:  true,
	}
	assets[1] = githubutils.ReleaseAssetSpec{
		Name:       "glooctl-darwin-amd64",
		ParentPath: "_output",
		UploadSHA:  true,
	}
	assets[2] = githubutils.ReleaseAssetSpec{
		Name:       "glooctl-windows-amd64.exe",
		ParentPath: "_output",
		UploadSHA:  true,
	}
	assets[3] = githubutils.ReleaseAssetSpec{
		Name:       "glooe-release.yaml",
		ParentPath: "install/manifest",
	}
	assets[4] = githubutils.ReleaseAssetSpec{
		Name:       "gloo-with-read-only-ui-release.yaml",
		ParentPath: "install/manifest",
	}
	spec := githubutils.UploadReleaseAssetSpec{
		Owner:             "solo-io",
		Repo:              "solo-projects",
		Assets:            assets,
		SkipAlreadyExists: true,
	}
	githubutils.UploadReleaseAssetCli(&spec)
}
