# CRD Reference Documentation

### Background
We are trying to migrate our APIs and CRDs to be built using Kubebuilder.

Historically, our docs are generated using Protocol Buffers as the source of truth, and code generation (https://github.com/solo-io/solo-kit/tree/main/pkg/code-generator/docgen) to generate the Markdown for these APIs.

With the migration to Kubebuilder, protos are no longer the source of truth, and we must rely on other means to
generate documentation. As https://github.com/kubernetes-sigs/controller-tools/issues/240 calls out, there isn't an agreed
upon approach yet. We chose to use https://github.com/fybrik/crdoc as it allows us to generate the Markdown for our docs,
which is used by Hugo. Other tools produced HTML, which wouldn't have worked for our current setup.

### Which CRDs are documented using this mechanism?
We opted to only add this to our newer APIs, which do not rely on protos as the source of truth. This means the following CRDs are affected:
- GatewayParameters
- DirectReponseAction

### How are the reference docs for our proto-based APIs generated?
We rely on our [solo-kit docgen](https://github.com/solo-io/solo-kit/tree/main/pkg/code-generator/docgen) code.

### How do I regenerate the documentation locally?
1. Install crdoc locally
```bash
make install-go-gools
```

2. Run the generation script.

This can be done by either using the go script directly:
```bash
 go run docs/content/crds/generate.go
```
or relying on the Make target:
```bash
 make generate-crd-reference-docs
```

### When is this regenerated in CI?
When running `make generated-code`.
