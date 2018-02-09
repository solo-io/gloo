package aws

import (
	"log"
	"os"
	"testing"
)

func TestFetchingAllVersions(t *testing.T) {
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if id == "" || secret == "" {
		return // skip the test
	}
	region := "us-east-1"
	lambdas, err := AWSFetcher(region, AccessToken{ID: id, Secret: secret})
	if err != nil {
		t.Errorf("unable to get aws lambdas %q", err)
	}
	if len(lambdas) == 0 {
		t.Errorf("expected at least one lambda; verify in aws for region %s", region)
	}
	for _, l := range lambdas {
		log.Println(l)
	}
}
