package main

import (
	"flag"
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/solo-io/skv2/codegen/render"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/solo-io/skv2/codegen/writer"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"log"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	codegen_placement "github.com/solo-io/skv2-enterprise/multicluster-admission-webhook/pkg/codegen/placement"
	"github.com/solo-io/skv2/codegen"
	"github.com/solo-io/skv2/codegen/model"
	"github.com/solo-io/skv2/codegen/skv2_anyvendor"
	"github.com/solo-io/skv2/contrib"
	"github.com/solo-io/solo-kit/pkg/code-generator/sk_anyvendor"
	"github.com/solo-io/solo-projects/codegen/groups"
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
	// Pull protos from private repositories
	os.Setenv("GOPRIVATE", "github.com/solo-io")
	anyvendorImports.External["github.com/solo-io/skv2-enterprise"] = []string{
		"**/multicluster-admission-webhook/api/multicluster/v1alpha1/*.proto",
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

	// TECHNICAL DEBT: https://github.com/solo-io/solo-projects/issues/2777
	// Gloo Fed CRDs are generated using skv2
	// To support v1 CRDs we need to update to the latest version of skv2
	// However, recent upgrades to skv2 cause protobuf registration panics in
	// Gloo Edge (https://github.com/solo-io/solo-projects/issues/2673)
	// Until we resolve that issue, we cannot upgrade to the latest version of skv2.
	// This creates 2 pieces of technical debt, that we are intentionally accepting:
	// 1. The skv2MulticlusterCRDs are generated in skv2 and copied to solo-projects.
	//		We now manually maintain that file, because until we upgrade, those CRDs in
	//		skv2 are still v1beta1
	// 2. The Gloo Fed CRDs are generated re-using skv2 code, but are not completed as part
	//		of the normal codegen command. This is because we have not figured out how to create
	//		CRDs from Gloo Fed protos, so we need to generate empty validation schemas

	if !*apiserverProtosOnly {
		// TODO(ilackarms): fix this hack - we copy some skv2 CRDs out of vendor_any into our helm chart.
		// (sam-heilbron) Technical Debt: outlined above why we do not call copySkv2MulticlusterCRDs
		// https://github.com/solo-io/solo-projects/issues/2777
		//copySkv2MulticlusterCRDs()
	}

	// (sam-heilbron) Technical Debt: outlined above why we render CRDs with placeholder schemas
	// https://github.com/solo-io/solo-projects/issues/2777
	if err := renderCrdsWithPlaceholderValidationSchema(skv2Cmd); err != nil {
		log.Fatal(err)
	}

	log.Println("Finished generating code")
}

// (sam-heilbron) Technical Debt: outlined above why we do not call copySkv2MulticlusterCRDs
// https://github.com/solo-io/solo-projects/issues/2777
func copySkv2MulticlusterCRDs() {
	vendoredMultiClusterCRDs := "vendor_any/github.com/solo-io/skv2/crds/multicluster.solo.io_v1alpha1_crds.yaml"
	importedMultiClusterCRDs := "install/helm/gloo-fed/crds/multicluster.solo.io_v1alpha1_imported_crds.yaml"
	if err := os.Rename(vendoredMultiClusterCRDs, importedMultiClusterCRDs); err != nil {
		log.Fatal(err)
	}
}

func renderCrdsWithPlaceholderValidationSchema(skv2Command codegen.Command) error {
	fileWriter := &writer.DefaultFileWriter{
		Root:   util.GetModuleRoot(),
		Header: "Code generated by skv2. DO NOT EDIT.",
	}

	for _, group := range skv2Command.Groups {

		// Render the manifests, using the placeholder validation schema
		defaultManifestsRenderer := render.ManifestsRenderer{
			AppName:     skv2Command.AppName,
			ManifestDir: skv2Command.ManifestRoot,
			ProtoDir:    skv2Command.ProtoDir,
			ResourceFuncs: map[render.OutFile]render.MakeResourceFunc{
				{
					Path: skv2Command.ManifestRoot + "/crds/" + group.Group + "_" + group.Version + "_" + "crds.yaml",
				}: func(group model.Group) ([]metav1.Object, error) {
					return CustomResourceDefinitionsWithPlaceholderValidationSchema(group)
				},
			},
		}
		manifests, err := defaultManifestsRenderer.RenderManifests(group)
		if err != nil {
			return err
		}

		if err := fileWriter.WriteFiles(manifests); err != nil {
			return err
		}

	}
	return nil
}

// Create CRDs for a group
func CustomResourceDefinitionsWithPlaceholderValidationSchema(
	group model.Group,
) (objects []metav1.Object, err error) {
	for _, resource := range group.Resources {

		// This schema is a placeholder because it accepts all configuration, and doesn't prune fields
		placeholderSchema := &v1.CustomResourceValidation{
			OpenAPIV3Schema: &v1.JSONSchemaProps{
				Type: "object",
				Properties: map[string]v1.JSONSchemaProps{
					"spec": {
						Type:                   "object",
						XPreserveUnknownFields: pointer.BoolPtr(true),
					},
					"status": {
						Type:                   "object",
						XPreserveUnknownFields: pointer.BoolPtr(true),
					},
				},
			},
		}

		objects = append(objects, CustomResourceDefinition(resource, placeholderSchema))
	}
	return objects, nil
}

// extracted from skv2@v0.20.0 https://github.com/solo-io/skv2/blob/master/codegen/kuberesource/crd.go
func CustomResourceDefinition(
	resource model.Resource,
	validationSchema *v1.CustomResourceValidation,
) *v1.CustomResourceDefinition {

	group := resource.Group.Group
	version := resource.Group.Version
	kind := resource.Kind
	kindLowerPlural := strings.ToLower(Pluralize(kind))
	kindLower := strings.ToLower(kind)

	var status *v1.CustomResourceSubresourceStatus
	if resource.Status != nil {
		status = &v1.CustomResourceSubresourceStatus{}
	}

	scope := v1.NamespaceScoped
	if resource.ClusterScoped {
		scope = v1.ClusterScoped
	}

	crd := &v1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", kindLowerPlural, group),
		},
		Spec: v1.CustomResourceDefinitionSpec{
			Group: group,
			Scope: scope,
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:    version,
					Served:  true,
					Storage: true,
					Subresources: &v1.CustomResourceSubresources{
						Status: status,
					},
					Schema: validationSchema,
				},
			},
			Names: v1.CustomResourceDefinitionNames{
				Plural:     kindLowerPlural,
				Singular:   kindLower,
				Kind:       kind,
				ShortNames: resource.ShortNames,
				ListKind:   kind + "List",
			},
		},
	}

	if validationSchema != nil {
		// Setting PreserveUnknownFields to false ensures that objects with unknown fields are rejected.
		crd.Spec.PreserveUnknownFields = false
	}
	return crd
}

// Define cases for pluralizing which pluralize library does not handle
var SpecialCases = map[string]string{
	"schema": "schemas",
}

// Pluralize is the canonical pluralization function for SKv2. It should be used to special case
// when we want a different result than the underlying pluralize library
func Pluralize(s string) string {
	c := pluralize.NewClient()
	for singular, plural := range SpecialCases {
		c.AddIrregularRule(singular, plural)
	}
	return c.Plural(s)
}
