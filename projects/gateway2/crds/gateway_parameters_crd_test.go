package crds

import (
	"os"
	"path/filepath"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func TestGatewayParametersReplicasValidationMinimum(t *testing.T) {
	crdPath := filepath.Join(getDirectory(), "..", "..", "..", "install", "helm", "gloo", "crds", "gateway.gloo.solo.io_gatewayparameters.yaml")
	raw, err := os.ReadFile(crdPath)
	if err != nil {
		t.Fatalf("read crd file: %v", err)
	}

	var crd apiextensionsv1.CustomResourceDefinition
	if err := yaml.Unmarshal(raw, &crd); err != nil {
		t.Fatalf("unmarshal crd yaml: %v", err)
	}

	var v1alpha1Schema *apiextensionsv1.JSONSchemaProps
	for _, version := range crd.Spec.Versions {
		if version.Name == "v1alpha1" && version.Schema != nil && version.Schema.OpenAPIV3Schema != nil {
			v1alpha1Schema = version.Schema.OpenAPIV3Schema
			break
		}
	}
	if v1alpha1Schema == nil {
		t.Fatalf("v1alpha1 schema not found in GatewayParameters CRD")
	}

	replicasSchema := v1alpha1Schema.
		Properties["spec"].
		Properties["kube"].
		Properties["deployment"].
		Properties["replicas"]

	if replicasSchema.Minimum == nil {
		t.Fatalf("spec.kube.deployment.replicas must define minimum validation")
	}
	if *replicasSchema.Minimum != 1 {
		t.Fatalf("spec.kube.deployment.replicas minimum must be 1, got %v", *replicasSchema.Minimum)
	}
}
