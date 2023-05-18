package groups

import (
	externalapis "github.com/solo-io/external-apis/codegen"
	"github.com/solo-io/skv2/codegen/model"
	"github.com/solo-io/skv2/contrib"
	hubgen "github.com/solo-io/solo-projects/codegen"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	module  = "github.com/solo-io/solo-projects"
	version = "v1"
	apiRoot = "projects/gloo-fed/pkg/api"

	discoveryInputSnapshotCodePath      = "input/snapshot.go"
	discoveryReconcilerSnapshotCodePath = "input/reconciler.go"

	AllGroups, ApiserverGroups []model.Group
)

func init() {
	AllGroups = []model.Group{
		FedGroup,
		GatewayGroup,
		FedGatewayGroup,
		GlooGroup,
		FedGlooGroup,
		EnterpriseGlooGroup,
		FedEnterpriseGroup,
		RateLimitGroup,
		FedRateLimitGroup,
		MultiClusterAdmissionGroup,
	}

	ApiserverGroups = []model.Group{
		FedApiserverGroup,
	}
}

var (
	allFederatedResourceTemplates = append(hubgen.FederatedResourceTemplates, contrib.AllGroupCustomTemplates...)
	allBaseGlooResourceTemplates  = hubgen.BaseGlooResourceTemplates
)

var GatewayGroup = makeGroup(
	"gateway",
	version,
	false,
	allBaseGlooResourceTemplates,
	[]resourceToGenerate{
		{kind: "Gateway"},
		{kind: "MatchableHttpGateway"},
		{kind: "MatchableTcpGateway"},
		{kind: "VirtualService"},
		{kind: "RouteTable"},
	})

var FedGatewayGroup = makeGroup(
	"fed.gateway",
	version,
	true,
	allFederatedResourceTemplates,
	[]resourceToGenerate{
		{kind: "FederatedGateway"},
		{kind: "FederatedMatchableHttpGateway"},
		{kind: "FederatedMatchableTcpGateway"},
		{kind: "FederatedVirtualService"},
		{kind: "FederatedRouteTable"},
	})

var GlooGroup = makeGroup(
	"gloo",
	version,
	false,
	allBaseGlooResourceTemplates,
	[]resourceToGenerate{
		{kind: "Upstream"},
		{kind: "UpstreamGroup"},
		{kind: "Settings"},
		{kind: "Proxy"},
	})

var FedGlooGroup = makeGroup(
	"fed.gloo",
	version,
	true,
	allFederatedResourceTemplates,
	[]resourceToGenerate{
		{kind: "FederatedUpstream"},
		{kind: "FederatedUpstreamGroup"},
		{kind: "FederatedSettings"},
	})

var EnterpriseGlooGroup = makeGroup(
	"enterprise.gloo",
	version,
	false,
	allBaseGlooResourceTemplates,
	[]resourceToGenerate{
		{kind: "AuthConfig"},
	})

var FedEnterpriseGroup = makeGroup(
	"fed.enterprise.gloo",
	version,
	true,
	allFederatedResourceTemplates,
	[]resourceToGenerate{
		{kind: "FederatedAuthConfig"},
	})

var RateLimitGroup = makeGroup(
	"ratelimit.api",
	"v1alpha1",
	false,
	allBaseGlooResourceTemplates,
	[]resourceToGenerate{
		{kind: "RateLimitConfig"},
	})

var FedRateLimitGroup = makeGroup(
	"fed.ratelimit",
	"v1alpha1",
	true,
	allFederatedResourceTemplates,
	[]resourceToGenerate{
		{kind: "FederatedRateLimitConfig"},
	})

var FedGroup = makeGroup(
	"fed",
	version,
	true,
	fedTemplates(),
	[]resourceToGenerate{
		{kind: "GlooInstance"},
		{kind: "FailoverScheme"},
	})

var FedApiserverGroup = makeGroup(
	"fed.rpc",
	version,
	false,
	[]model.CustomTemplates{},
	[]resourceToGenerate{})

var MultiClusterAdmissionGroup = makeMultiClusterAdmissionGroup()

func fedTemplates() []model.CustomTemplates {
	inputDiscoverySnapshot := map[schema.GroupVersion][]string{
		corev1.SchemeGroupVersion: {
			"Service",
			"Pod",
		},
		appsv1.SchemeGroupVersion: {
			"Deployment",
			"DaemonSet",
		},
		schema.GroupVersion{
			Group:   "gloo.solo.io",
			Version: "v1",
		}: {
			"Upstream",
			"UpstreamGroup",
			"Proxy",
			"Settings",
		},
		schema.GroupVersion{
			Group:   "gateway.solo.io",
			Version: "v1",
		}: {
			"Gateway",
			"MatchableHttpGateway",
			"MatchableTcpGateway",
			"VirtualService",
			"RouteTable",
		},
		schema.GroupVersion{
			Group:   "enterprise.gloo.solo.io",
			Version: "v1",
		}: {
			"AuthConfig",
		},
		schema.GroupVersion{
			Group:   "ratelimit.api.solo.io",
			Version: "v1alpha1",
		}: {
			"RateLimitConfig",
		},
	}

	baseGlooGroups := []model.Group{GatewayGroup, GlooGroup, EnterpriseGlooGroup, RateLimitGroup}
	selectFromGroups := map[string][]model.Group{
		"github.com/solo-io/solo-apis":     baseGlooGroups,
		"github.com/solo-io/external-apis": externalapis.Groups,
	}

	crossGroupTemplateParams := contrib.SnapshotTemplateParameters{
		SelectFromGroups: selectFromGroups,
		SnapshotResources: &contrib.HomogenousSnapshotResources{
			ResourcesToSelect: inputDiscoverySnapshot,
		},
	}
	crossGroupTemplateParams.OutputFilename = discoveryInputSnapshotCodePath
	snapshot := contrib.InputSnapshot(crossGroupTemplateParams)

	crossGroupTemplateParams.OutputFilename = discoveryReconcilerSnapshotCodePath
	reconciler := contrib.InputReconciler(crossGroupTemplateParams)
	return append(contrib.AllGroupCustomTemplates, snapshot, reconciler)
}

type resourceToGenerate struct {
	kind     string
	noStatus bool // don't put a status on this resource
}

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
		Module:                  module,
		Resources:               resources,
		RenderManifests:         render,
		RenderTypes:             render,
		RenderClients:           render,
		RenderController:        render,
		MockgenDirective:        render,
		CustomTemplates:         customTemplates,
		ApiRoot:                 apiRoot,
		RenderValidationSchemas: false,
	}
}

func makeMultiClusterAdmissionGroup() model.Group {
	return model.Group{
		GroupVersion: schema.GroupVersion{
			Group:   "multicluster.solo.io",
			Version: "v1alpha1",
		},
		Module: module,
		Resources: []model.Resource{
			{
				Kind: "MultiClusterRole",
				Group: &model.Group{
					GroupVersion: schema.GroupVersion{
						Group:   "multicluster.solo.io",
						Version: "v1alpha1",
					},
				},
				Spec: model.Field{
					Type: model.Type{
						Name: "MultiClusterRoleSpec",
					},
				},
				Status: &model.Field{
					Type: model.Type{
						Name: "MultiClusterRoleStatus",
					},
				},
			},
			{
				Kind: "MultiClusterRoleBinding",
				Group: &model.Group{
					GroupVersion: schema.GroupVersion{
						Group:   "multicluster.solo.io",
						Version: "v1alpha1",
					},
				},
				Spec: model.Field{
					Type: model.Type{
						Name: "MultiClusterRoleBindingSpec",
					},
				},
				Status: &model.Field{
					Type: model.Type{
						Name: "MultiClusterRoleBindingStatus",
					},
				},
			},
		},
		RenderManifests:         true,
		RenderTypes:             true,
		RenderClients:           true,
		MockgenDirective:        true,
		RenderFieldJsonDeepcopy: false,
		ApiRoot:                 apiRoot,
	}
}
