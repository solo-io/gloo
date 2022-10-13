package main

import (
	"flag"
	"log"
	"os"

	"github.com/solo-io/skv2/codegen"
	"github.com/solo-io/skv2/codegen/model"
	"github.com/solo-io/skv2/codegen/skv2_anyvendor"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/solo-io/skv2/contrib"
	"github.com/solo-io/solo-projects/codegen/groups"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/codegen/chart"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/codegen/placement"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// This is intended to be an example of how to import the multicluster admission webhook into a project.
// It was copied over (along with the rest of the multicluster-admission-webhook) from skv2-enterprise into solo-projects
// to remove our direct dependency on the skv2-enterprise repo.
// It is not currently run as part of CI, but is left here for reference.

var (
	rootDir          = "projects/multicluster-admission-webhook/test/internal/"
	skv2Imports      = skv2_anyvendor.CreateDefaultMatchOptions([]string{rootDir + "api/*.proto"})
	apiRoot          = rootDir + "pkg"
	testManifestRoot = rootDir + "helm/charts/test"
	testGroups       = []model.Group{
		{
			GroupVersion: schema.GroupVersion{
				Group:   "test.multicluster.solo.io",
				Version: "v1alpha1",
			},
			Module: "github.com/solo-io/solo-projects",
			Resources: []model.Resource{
				{
					Kind: "Test",
					Spec: model.Field{
						Type: model.Type{
							Name:      "TestSpec",
							GoPackage: "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/test/internal/api/types",
						},
					},
				},
			},
			RenderManifests: true,
			RenderTypes:     true,
			RenderClients:   true,
			ApiRoot:         rootDir + "api",
		},
	}
)

func main() {
	log.Println("starting generate")
	// change wd to codegen dir
	codegenDir := util.MustGetThisDir()
	if err := os.Chdir(codegenDir); err != nil {
		log.Fatal(err)
	}

	chartOnly := flag.Bool("chart", false, "only generate the helm chart")
	flag.Parse()

	if err := generateRbacValidatingWebhook().Execute(); err != nil {
		log.Fatal(err)
	}
	if !*chartOnly {
		skv2Cmd := codegen.Command{
			AppName:         "multicluster-admission-webhook",
			Groups:          testGroups,
			AnyVendorConfig: skv2Imports,
			RenderProtos:    true,
			ManifestRoot:    testManifestRoot,
		}
		if err := skv2Cmd.Execute(); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Finished generating code")
}

func generateRbacValidatingWebhook() codegen.Command {
	topLevelTemplates := []model.CustomTemplates{
		placement.TypedParser(contrib.SnapshotTemplateParameters{
			OutputFilename: apiRoot + "/placement/generated_parser.go",
			SelectFromGroups: map[string][]model.Group{
				"github.com/solo-io/solo-projects": testGroups,
			},
			SnapshotResources: contrib.HomogenousSnapshotResources{
				ResourcesToSelect: map[schema.GroupVersion][]string{
					schema.GroupVersion{
						Group:   "test.multicluster.solo.io",
						Version: "v1alpha1",
					}: {
						"Test",
					},
				},
			},
		}),
	}
	// change the apiRoot used for these tests so we don't modify the gloo-fed copy
	mcGroup := groups.MultiClusterAdmissionGroup
	mcGroup.ApiRoot = "projects/multicluster-admission-webhook/pkg/api"
	groups := []model.Group{
		mcGroup, // Generate CRD definitions for MultiClusterRole/RoleBinding
	}
	chart := chart.GenerateChart(
		model.Data{
			ApiVersion:  "v1",
			Description: "multicluster-admission-webhook test chart.",
			Name:        "multicluster-admission-webhook-test",
			Version:     "0.0.1",
			Home:        "https://solo.io",
		},
		testGroups,
		nil,
	)
	return codegen.Command{
		AppName:           "test",
		AnyVendorConfig:   skv2Imports,
		ManifestRoot:      testManifestRoot,
		Groups:            groups,
		TopLevelTemplates: topLevelTemplates,
		Chart:             chart,
	}
}
