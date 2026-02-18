# HPA Support for Gateway API Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Enable HPA (HorizontalPodAutoscaler) support for auto-provisioned Gateway proxies by making the `replicas` field optional in rendered Deployments, following the approach from kgateway-dev/kgateway#12548.

**Architecture:** When `replicas` is not explicitly set in GatewayParameters, the Deployment spec will omit the `replicas` field entirely, allowing Kubernetes (default: 1 replica) or an externally-managed HPA to control scaling. When `replicas` IS explicitly set, it renders as before. This replaces the previous (commented-out) approach of managing an HPA through the GatewayParameters API.

**Tech Stack:** Go, Helm templates, Kubernetes Gateway API, Ginkgo/Gomega tests

---

### Task 1: Update Helm Deployment Template to Conditionally Render Replicas

**Files:**
- Modify: `projects/gateway2/helm/gloo-gateway/templates/gateway/proxy-deployment.yaml:29-32`
- Modify: `projects/gateway2/helm/gloo-gateway/values.yaml:17-18`

**Step 1: Update the proxy-deployment.yaml template**

In `projects/gateway2/helm/gloo-gateway/templates/gateway/proxy-deployment.yaml`, change lines 29-32 from:

```yaml
  {{- if not $gateway.autoscaling.enabled }}
  replicas: {{ $gateway.replicaCount }}
  {{- end }}
```

to:

```yaml
  {{- if $gateway.replicaCount }}
  replicas: {{ $gateway.replicaCount }}
  {{- end }}
```

This makes replicas render only when the value is explicitly set, rather than gating on the autoscaling flag. When `replicaCount` is nil/0/empty, the field is omitted from the Deployment, allowing K8s (default 1) or an HPA to control replicas.

**Step 2: Update values.yaml to remove default replicaCount**

In `projects/gateway2/helm/gloo-gateway/values.yaml`, change:

```yaml
  # actual default set in default GatewayParam proxyDeployment.replicas
  replicaCount: 1
```

to:

```yaml
  # replicaCount is intentionally unset here. When omitted, the Deployment's replicas
  # field will not be rendered, letting K8s manage it (default: 1) or an HPA control it.
  # Set explicitly in GatewayParameters to override.
  # replicaCount:
```

**Step 3: Remove the autoscaling default values from values.yaml**

Since autoscaling is now handled by omitting replicas (not via a built-in HPA), remove the leftover autoscaling config block:

```yaml
  # leftover autoscaling config; not actually wired up for public use yet
  # see: https://github.com/solo-io/solo-projects/issues/5948
  autoscaling:
    enabled: false
    minReplicas: 1
    maxReplicas: 100
    targetCPUUtilizationPercentage: 80
    # targetMemoryUtilizationPercentage: 80
```

Replace with just enough for the template to not break:

```yaml
  autoscaling:
    enabled: false
```

**Step 4: Verify the template renders correctly by running tests**

Run: `cd /Users/wkrause/development/gloo && go test ./projects/gateway2/deployer/... -v -run "TestSuite" -count=1 2>&1 | head -50`

Expected: Tests may fail at this point since deployer code still sets `ReplicaCount`. That's fine -- we'll fix in subsequent tasks.

**Step 5: Commit**

```bash
git add projects/gateway2/helm/gloo-gateway/templates/gateway/proxy-deployment.yaml projects/gateway2/helm/gloo-gateway/values.yaml
git commit -m "feat: make replicas conditional in gateway proxy deployment template

Only render the replicas field in the Deployment when explicitly set,
allowing K8s or an external HPA to control scaling when omitted."
```

---

### Task 2: Update Default GatewayParameters Helm Template

**Files:**
- Modify: `install/helm/gloo/templates/43-gatewayparameters.yaml:18-25`

**Step 1: Make the replicas field conditional in the default GatewayParameters**

In `install/helm/gloo/templates/43-gatewayparameters.yaml`, change lines 18-25 from:

```yaml
{{- $replicas := 1 -}}
{{- if $gg.proxyDeployment -}}
{{- if $gg.proxyDeployment.replicas -}}
{{- $replicas = $gg.proxyDeployment.replicas -}}
{{- end -}}{{/* if $gg.proxyDeployment.replicas */}}
{{- end }}{{/* if $gg.proxyDeployment */}}
    deployment:
      replicas: {{ $replicas }}
```

to:

```yaml
{{- if and $gg.proxyDeployment $gg.proxyDeployment.replicas }}
    deployment:
      replicas: {{ $gg.proxyDeployment.replicas }}
{{- end }}{{/* if proxyDeployment.replicas */}}
```

This means the default GatewayParameters will only include `deployment.replicas` when explicitly configured in the helm values. When not set, the field is omitted entirely, which means the deployer won't set `ReplicaCount` and the Deployment will omit `replicas`, enabling HPA support by default.

**Step 2: Commit**

```bash
git add install/helm/gloo/templates/43-gatewayparameters.yaml
git commit -m "feat: make replicas optional in default GatewayParameters

Only include deployment.replicas in the default GatewayParameters when
explicitly configured, enabling HPA support by default."
```

---

### Task 3: Update Deployer Code

**Files:**
- Modify: `projects/gateway2/deployer/deployer.go:110-122,332-343`
- Modify: `projects/gateway2/deployer/values_helpers.go:92-111`

**Step 1: Clean up the commented-out autoscaling code in deployer.go**

In `projects/gateway2/deployer/deployer.go`, replace lines 332-343:

```go
	// deployment values
	gateway.ReplicaCount = deployConfig.GetReplicas()

	// TODO: The follow stanza has been commented out as autoscaling support has been removed.
	// see https://github.com/solo-io/solo-projects/issues/5948 for more info.
	//
	// autoscalingVals := getAutoscalingValues(kubeProxyConfig.GetAutoscaling())
	// vals.Gateway.Autoscaling = autoscalingVals
	// if autoscalingVals == nil && deployConfig.GetReplicas() != nil {
	// 	replicas := deployConfig.GetReplicas().GetValue()
	// 	vals.Gateway.ReplicaCount = &replicas
	// }
```

with:

```go
	// deployment values
	// Only set ReplicaCount when explicitly configured. When nil, the Deployment
	// template will omit the replicas field, letting K8s or an HPA control scaling.
	gateway.ReplicaCount = deployConfig.GetReplicas()
```

**Step 2: Update the GetGvksToWatch comment in deployer.go**

In `projects/gateway2/deployer/deployer.go`, update the comment block at lines 110-122. Change:

```go
	// In order to get the GVKs for the resources to watch, we need:
	// - a placeholder Gateway (only the name and namespace are used, but the actual values don't matter,
	//   as we only care about the GVKs of the rendered resources)
	// - the minimal values that render all the proxy resources (HPA is not included because it's not
	//   fully integrated/working at the moment)
	// - a flag to indicate whether mtls is enabled, so we can render the secret if needed
```

to:

```go
	// In order to get the GVKs for the resources to watch, we need:
	// - a placeholder Gateway (only the name and namespace are used, but the actual values don't matter,
	//   as we only care about the GVKs of the rendered resources)
	// - the minimal values that render all the proxy resources
	// - a flag to indicate whether mtls is enabled, so we can render the secret if needed
```

**Step 3: Clean up the commented-out autoscaling code in values_helpers.go**

In `projects/gateway2/deployer/values_helpers.go`, replace lines 92-111:

```go
// TODO: Removing until autoscaling is re-added.
// See: https://github.com/solo-io/solo-projects/issues/5948
// Convert autoscaling values from GatewayParameters into helm values to be used by the deployer.
// func getAutoscalingValues(autoscaling *v1.Autoscaling) *helmAutoscaling {
// 	hpaConfig := autoscaling.HorizontalPodAutoscaler
// 	if hpaConfig == nil {
// 		return nil
// 	}

// 	trueVal := true
// 	autoscalingVals := &helmAutoscaling{
// 		Enabled: &trueVal,
// 	}
// 	autoscalingVals.MinReplicas = hpaConfig.MinReplicas
// 	autoscalingVals.MaxReplicas = hpaConfig.MaxReplicas
// 	autoscalingVals.TargetCPUUtilizationPercentage = hpaConfig.TargetCpuUtilizationPercentage
// 	autoscalingVals.TargetMemoryUtilizationPercentage = hpaConfig.TargetMemoryUtilizationPercentage

// 	return autoscalingVals
// }
```

with nothing (delete the entire block). The autoscaling approach is replaced by simply omitting replicas.

**Step 4: Commit**

```bash
git add projects/gateway2/deployer/deployer.go projects/gateway2/deployer/values_helpers.go
git commit -m "refactor: remove commented-out autoscaling code from deployer

HPA support is now enabled by omitting the replicas field from the
Deployment when not explicitly set, rather than via a built-in HPA."
```

---

### Task 4: Update GatewayParameters API Type Comment

**Files:**
- Modify: `projects/gateway2/api/v1alpha1/gateway_parameters_types.go:194-199`

**Step 1: Update the Replicas field comment**

In `projects/gateway2/api/v1alpha1/gateway_parameters_types.go`, change the `ProxyDeployment` struct:

```go
// Configuration for the Proxy deployment in Kubernetes.
type ProxyDeployment struct {
	// The number of desired pods. Defaults to 1.
	//
	// +kubebuilder:validation:Optional
	Replicas *uint32 `json:"replicas,omitempty"`
}
```

to:

```go
// Configuration for the Proxy deployment in Kubernetes.
type ProxyDeployment struct {
	// The number of desired pods.
	// If omitted, the Deployment's replicas field will not be set, letting the
	// Kubernetes control plane manage it (default: 1). This allows an external
	// HPA to control scaling without conflict.
	//
	// +kubebuilder:validation:Optional
	Replicas *uint32 `json:"replicas,omitempty"`
}
```

**Step 2: Commit**

```bash
git add projects/gateway2/api/v1alpha1/gateway_parameters_types.go
git commit -m "docs: update ProxyDeployment.Replicas comment for HPA support"
```

---

### Task 5: Update Tests - Deployer Tests

**Files:**
- Modify: `projects/gateway2/deployer/deployer_test.go`

**Step 1: Update the `validateGatewayParametersPropagation` function**

The validation function at line ~902 checks `Expect(dep.Spec.Replicas).ToNot(BeNil())` and assumes replicas is always set. We need to update it to handle the case where Replicas may be nil in the expected GatewayParameters.

Change the replicas assertion in `validateGatewayParametersPropagation` (around line 900-903) from:

```go
			dep := objs.findDeployment(defaultNamespace, defaultDeploymentName)
			Expect(dep).ToNot(BeNil())
			Expect(dep.Spec.Replicas).ToNot(BeNil())
			Expect(*dep.Spec.Replicas).To(Equal(int32(*expectedGwp.Deployment.Replicas)))
```

to:

```go
			dep := objs.findDeployment(defaultNamespace, defaultDeploymentName)
			Expect(dep).ToNot(BeNil())
			if expectedGwp.Deployment != nil && expectedGwp.Deployment.Replicas != nil {
				Expect(dep.Spec.Replicas).ToNot(BeNil())
				Expect(*dep.Spec.Replicas).To(Equal(int32(*expectedGwp.Deployment.Replicas)))
			} else {
				Expect(dep.Spec.Replicas).To(BeNil())
			}
```

**Step 2: Add a test entry for nil replicas (HPA-compatible)**

Add a new test entry in the `DescribeTable("create and validate objs"` block (after the "no listeners on gateway" entry, around line 1375):

```go
		Entry("nil replicas allows HPA to control scaling", &input{
			dInputs: defaultDeployerInputs(),
			gw:      defaultGateway(),
			defaultGwp: &gw2_v1alpha1.GatewayParameters{
				TypeMeta: metav1.TypeMeta{
					Kind:       gw2_v1alpha1.GatewayParametersKind,
					APIVersion: gw2_v1alpha1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      wellknown.DefaultGatewayParametersName,
					Namespace: defaultNamespace,
					UID:       "1237",
				},
				Spec: gw2_v1alpha1.GatewayParametersSpec{
					Kube: &gw2_v1alpha1.KubernetesProxyConfig{
						Service: &gw2_v1alpha1.Service{
							Type: ptr.To(corev1.ServiceTypeLoadBalancer),
						},
					},
				},
			},
		}, &expectedOutput{
			validationFunc: func(objs clientObjects, inp *input) error {
				dep := objs.findDeployment(defaultNamespace, defaultDeploymentName)
				Expect(dep).ToNot(BeNil())
				Expect(dep.Spec.Replicas).To(BeNil(), "replicas should be nil when not explicitly set, allowing HPA control")
				return nil
			},
		}),
```

**Step 3: Run the deployer tests**

Run: `cd /Users/wkrause/development/gloo && go test ./projects/gateway2/deployer/... -v -count=1 2>&1 | tail -30`

Expected: All tests pass.

**Step 4: Commit**

```bash
git add projects/gateway2/deployer/deployer_test.go
git commit -m "test: add HPA support test and update replicas assertions"
```

---

### Task 6: Update Tests - Merge Tests

**Files:**
- Modify: `projects/gateway2/deployer/merge_test.go`

**Step 1: Add merge test for nil replicas preservation**

Add a new test case in `merge_test.go` to verify that merging a GatewayParameters with no replicas set in `src` preserves the `dst` replicas (or keeps nil):

```go
	It("preserves nil replicas when src has no deployment", func() {
		dst := &gw2_v1alpha1.GatewayParameters{
			Spec: gw2_v1alpha1.GatewayParametersSpec{
				Kube: &gw2_v1alpha1.KubernetesProxyConfig{},
			},
		}
		src := &gw2_v1alpha1.GatewayParameters{
			Spec: gw2_v1alpha1.GatewayParametersSpec{
				Kube: &gw2_v1alpha1.KubernetesProxyConfig{},
			},
		}
		out := deepMergeGatewayParameters(dst, src)
		Expect(out.Spec.Kube.GetDeployment().GetReplicas()).To(BeNil())
	})
```

**Step 2: Run the merge tests**

Run: `cd /Users/wkrause/development/gloo && go test ./projects/gateway2/deployer/... -v -run "deepMerge" -count=1 2>&1 | tail -20`

Expected: All tests pass.

**Step 3: Commit**

```bash
git add projects/gateway2/deployer/merge_test.go
git commit -m "test: add merge test for nil replicas preservation"
```

---

### Task 7: Final Verification

**Step 1: Run all deployer tests**

Run: `cd /Users/wkrause/development/gloo && go test ./projects/gateway2/deployer/... -v -count=1 2>&1 | tail -40`

Expected: All tests pass.

**Step 2: Run the full gateway2 test suite**

Run: `cd /Users/wkrause/development/gloo && go test ./projects/gateway2/... -count=1 2>&1 | tail -20`

Expected: All tests pass.

**Step 3: Verify the helm chart renders correctly without replicas**

Run: `cd /Users/wkrause/development/gloo && go test ./projects/gateway2/deployer/... -v -run "should work with empty params" -count=1`

Expected: PASS - the empty params case should now produce a Deployment without a replicas field.

**Step 4: Compile check**

Run: `cd /Users/wkrause/development/gloo && go build ./projects/gateway2/...`

Expected: Build succeeds with no errors.
