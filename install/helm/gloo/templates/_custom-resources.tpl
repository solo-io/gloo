{{/* vim: set filetype=mustache: */}}

{{/*
Render Gloo Edge custom resources. This is done in a helper function, rather than included as a template,
to ensure that we don't try to apply the custom resources before the service backing the validation webhook
is ready.
When adding a new CR to the helm chart, the template should be included here.
*/}}
{{- define "gloo.customResources" -}}
{{- include "gloo.customResources.defaultGateways" . }}
---
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: default
  namespace: gloo-system
spec:
  virtualHost:
    routes:
      - matchers:
        - prefix: /
        directResponseAction:
          status: 200
          body: "Hello, world!"

{{ end }}
