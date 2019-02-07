package install

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"k8s.io/helm/pkg/manifest"
	"k8s.io/helm/pkg/releaseutil"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	owner    = "solo-io"
	repo     = "solo-projects"
	yamlName = "glooe-release.yaml"
)

func setupGithubClient() (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN env var must be present")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, nil
}

func readManifest(version string) ([]byte, error) {
	var (
		releaseId *int64
		assetId   *int64
		data      []byte

		ctx = context.Background()
	)
	// Get the data
	client, err := setupGithubClient()
	if err != nil {
		return nil, err
	}
	tag := fmt.Sprintf("v%s", version)
	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if release.GetTagName() == tag {
			releaseId = release.ID
		}
	}
	if releaseId == nil {
		return nil, fmt.Errorf("unable to find release with tag: %s", tag)
	}

	assets, _, err := client.Repositories.ListReleaseAssets(ctx, owner, repo, *releaseId, nil)
	for _, asset := range assets {
		if asset.GetName() == yamlName {
			assetId = asset.ID
		}
	}
	if assetId == nil {
		return nil, fmt.Errorf("unable to find %s for given tag: %s", yamlName, tag)
	}

	ra, redirectUrl, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, *assetId)
	if ra != nil {
		data, err = ioutil.ReadAll(ra)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := http.Get(redirectUrl)
		if err != nil {
			return nil, err
		}
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func parseYaml(bigFile []byte) []manifest.Manifest {
	manifests := releaseutil.SplitManifests(string(bigFile))
	return manifest.SplitManifests(manifests)
}
