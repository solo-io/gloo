# Multicluster Admission Webhook

A k8s [validating admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#validatingadmissionwebhook)
for enforcing RBAC on arbitrary CRDs.

## Usage

For integration with downstream projects, the multicluster admission webhook 
package provides code and Helm chart generation utilities.

The interface between this library and a consuming project consists in a `TypedParser`,
 which provides the mapping logic from arbitrary CRD types to [Placement objects](https://github.com/solo-io/skv2-enterprise/blob/487b73ba21018c1eed719c911630b6fa3ac22d2b/pkg/multicluster-admission-webhook/api/multicluster/v1alpha1/multicluster.proto#L18),
 which represent the clusters and namespaces that a CRD accesses (e.g. by translating
 output resources to those clusters and namespaces). A `TypedParser` interface is generated
 from the downstream projects' CRD definitions, to be implemented by the downstream
 project. The skv2 custom template is exposed in `codegen/placement/template.go`.

Regarding deployment, this package also exports a utility function for generating
the required Helm chart templates and CRD definitions in `codegen/chart/chart.go`.

For an example usage, see `test/internal/generate.go`.

## Note

The multicluster-admission-webhook code was previously located in the skv2-enterprise repo, but it has now been copied into
solo-projects. We may want to consider combining `projects/rbac-validating-webhook` and `projects/multicluster-admission-webhook`
since `rbac-validating-webhook` is basically just importing/calling `multicluster-admission-webhook`.