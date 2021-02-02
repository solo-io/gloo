package main

import (
	"log"
	"os"

	"github.com/solo-io/skv2/codegen"
	"github.com/solo-io/skv2/codegen/model"
	"github.com/solo-io/skv2/codegen/skv2_anyvendor"
	"github.com/solo-io/solo-kit/pkg/code-generator/sk_anyvendor"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	log.Println("starting generate")

	anyvendorImports := sk_anyvendor.CreateDefaultMatchOptions([]string{
		"api/**/*.proto",
	})
	anyvendorImports.External["github.com/solo-io/skv2"] = []string{
		"api/**/*.proto",
		"crds/*multicluster.solo.io_v1alpha1_crds.yaml",
	}
	anyvendorImports.External["github.com/solo-io/solo-apis"] = []string{
		"api/rate-limiter/**/*.proto",
		"api/gloo/**/*.proto",
	}
	// Pull protos from private repositories
	os.Setenv("GOPRIVATE", "github.com/solo-io")
	anyvendorImports.External["github.com/solo-io/skv2-enterprise"] = []string{
		"**/multicluster-admission-webhook/api/multicluster/v1alpha1/*.proto",
	}

	skv2Cmd := codegen.Command{
		AppName:      "gloo-fed",
		RenderProtos: true,
		Groups:       AllGroups,
		AnyVendorConfig: &skv2_anyvendor.Imports{
			Local:    anyvendorImports.Local,
			External: anyvendorImports.External,
		},
		ManifestRoot: "install/helm/charts/gloo-fed",
	}
	if err := skv2Cmd.Execute(); err != nil {
		log.Fatal(err)
	}

	log.Println("Finished generating code")
}

var (
	module    = "github.com/solo-io/solo-projects/projects/gloo-fed"
	apiRoot   = "pkg/api"
	AllGroups []model.Group
)

type resourceToGenerate struct {
	kind     string
	noStatus bool // don't put a status on this resource
}

func init() {
	AllGroups = []model.Group{
		FedGroup,
		FedGatewayGroup,
		FedGlooGroup,
		FedEnterpriseGroup,
	}
}

var FedGroup = makeGroup(
	"fed",
	"v1",
	true,
	nil,
	[]resourceToGenerate{
		{kind: "GlooInstance"},
		{kind: "FailoverScheme"},
	})

var FedGlooGroup = makeGroup(
	"fed.gloo",
	"v1",
	true,
	nil,
	[]resourceToGenerate{
		{kind: "FederatedUpstream"},
		{kind: "FederatedUpstreamGroup"},
		{kind: "FederatedSettings"},
	})

var FedGatewayGroup = makeGroup(
	"fed.gateway",
	"v1",
	true,
	nil,
	[]resourceToGenerate{
		{kind: "FederatedGateway"},
		{kind: "FederatedVirtualService"},
		{kind: "FederatedRouteTable"},
	})

var FedEnterpriseGroup = makeGroup(
	"fed.enterprise.gloo",
	"v1",
	true,
	nil,
	[]resourceToGenerate{
		{kind: "FederatedAuthConfig"},
	})

func makeGroup(
	groupPrefix, version string,
	render bool,
	customTemplates []model.CustomTemplates,
	resourcesToGenerate []resourceToGenerate,
) model.Group {
	var resources []model.Resource
	for _, resource := range resourcesToGenerate {
		res := model.Resource{
			Kind: resource.kind,
			Spec: model.Field{
				Type: model.Type{
					Name: resource.kind + "Spec",
				},
			},
		}
		if !resource.noStatus {
			res.Status = &model.Field{Type: model.Type{
				Name: resource.kind + "Status",
			}}
		}
		resources = append(resources, res)
	}

	return model.Group{
		GroupVersion: schema.GroupVersion{
			Group:   groupPrefix + "." + "solo.io",
			Version: version,
		},
		Module:           module,
		Resources:        resources,
		RenderManifests:  render,
		RenderTypes:      render,
		RenderClients:    render,
		RenderController: render,
		MockgenDirective: render,
		CustomTemplates:  customTemplates,
		ApiRoot:          apiRoot,
	}
}
