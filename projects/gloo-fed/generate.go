package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/solo-io/skv2/codegen"
	"github.com/solo-io/skv2/codegen/model"
	"github.com/solo-io/skv2/codegen/skv2_anyvendor"
	"github.com/solo-io/skv2/contrib"
	"github.com/solo-io/solo-kit/pkg/code-generator/sk_anyvendor"
	"github.com/solo-io/solo-projects/codegen/groups"
	codegen_placement "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/codegen/placement"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	rbacEnabledKinds = map[schema.GroupVersion][]string{
		schema.GroupVersion{
			Group:   "fed.gloo.solo.io",
			Version: "v1",
		}: {
			"FederatedUpstream",
			"FederatedSettings",
			"FederatedUpstreamGroup",
		},
		schema.GroupVersion{
			Group:   "fed.gateway.solo.io",
			Version: "v1",
		}: {
			"FederatedGateway",
			"FederatedMatchableHttpGateway",
			"FederatedRouteTable",
			"FederatedVirtualService",
		},
		schema.GroupVersion{
			Group:   "fed.enterprise.gloo.solo.io",
			Version: "v1",
		}: {
			"FederatedAuthConfig",
		},
		schema.GroupVersion{
			Group:   "fed.ratelimit.solo.io",
			Version: "v1alpha1",
		}: {
			"FederatedRateLimitConfig",
		},
		schema.GroupVersion{
			Group:   "fed.solo.io",
			Version: "v1",
		}: {
			"FailoverScheme",
		},
	}
)

func main() {
	os.RemoveAll("vendor_any")
	log.Println("starting generate")

	anyvendorImports := sk_anyvendor.CreateDefaultMatchOptions(
		[]string{
			"projects/apiserver/**/*.proto",
			"projects/gloo-fed/**/*.proto",
			"projects/gloo/api/enterprise/**/*.proto",
		},
	)
	anyvendorImports.External["github.com/solo-io/skv2"] = []string{
		"api/**/*.proto",
		"crds/*multicluster.solo.io_v1alpha1_crds.yaml",
	}
	anyvendorImports.External["github.com/solo-io/solo-apis"] = []string{
		"api/rate-limiter/**/*.proto",
		"api/gloo/**/*.proto",
	}
	// Need to create copy of fed group to pass just failover schemes into chart
	fedGroup := groups.FedGroup
	fedGroup.Resources = nil
	for _, resource := range groups.FedGroup.Resources {
		if strings.Contains(resource.Spec.Type.Name, "Failover") {
			fedGroup.Resources = append(fedGroup.Resources, resource)
		}
	}
	apiserverProtosOnly := flag.Bool("apiserver", false, "only generate the apiserver protos")
	flag.Parse()

	skv2Cmd := codegen.Command{
		AppName:      "gloo-fed",
		RenderProtos: true,
		AnyVendorConfig: &skv2_anyvendor.Imports{
			Local:    anyvendorImports.Local,
			External: anyvendorImports.External,
		},
		ManifestRoot: "install/helm/gloo-fed",
	}

	if *apiserverProtosOnly {
		skv2Cmd.Groups = groups.ApiserverGroups
	} else {
		skv2Cmd.Groups = append(groups.AllGroups,
			model.Group{
				CustomTemplates: []model.CustomTemplates{
					codegen_placement.TypedParser(contrib.SnapshotTemplateParameters{
						OutputFilename: "placement/generated_parser.go",
						SelectFromGroups: map[string][]model.Group{
							"": {
								groups.FedEnterpriseGroup,
								groups.FedRateLimitGroup,
								groups.FedGlooGroup,
								groups.FedGatewayGroup,
								groups.FedGroup,
							},
						},
						SnapshotResources: contrib.HomogenousSnapshotResources{
							ResourcesToSelect: rbacEnabledKinds,
						},
					}),
				},
				ApiRoot: "projects/rbac-validating-webhook/pkg",
			},
		)
	}
	if err := skv2Cmd.Execute(); err != nil {
		log.Fatal(err)
	}

	if !*apiserverProtosOnly {
		// TODO(ilackarms): fix this hack - we copy some skv2 CRDs out of vendor_any into our helm chart.
		copySkv2MulticlusterCRDs()
	}

	log.Println("Finished generating code")
}

func copySkv2MulticlusterCRDs() {
	vendoredMultiClusterCRDs := "vendor_any/github.com/solo-io/skv2/crds/multicluster.solo.io_v1alpha1_crds.yaml"
	importedMultiClusterCRDs := "install/helm/gloo-fed/crds/multicluster.solo.io_v1alpha1_imported_crds.yaml"
	if err := os.Rename(vendoredMultiClusterCRDs, importedMultiClusterCRDs); err != nil {
		log.Fatal(err)
	}
}
