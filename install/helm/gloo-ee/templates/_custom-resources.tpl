{{/* vim: set filetype=mustache: */}}

{{/*
Render Gloo Edge Enterprise custom resources. This is done in a helper function, rather than included as a template,
to ensure that we don't try to apply the custom resources before the service backing the validation webhook
is ready.
When adding a new CR to the helm chart, the template should be included here.
As this gets the values context from OSS, $ will refer to the OSS values context.
If additional fields are required, add them at https://github.com/solo-io/gloo/blob/0429470a3f671b1137b36abe105f5df3d583d53f/install/helm/gloo/templates/5-resource-configmap.yaml#L17
*/}}
{{- define "gloo.extraCustomResources" -}}
{{- include "gloo.extraCustomResources.ratelimitUpstreams" . }}
{{- include "gloo.extraCustomResources.extauthUpstreams" . }}
{{- include "gloo.extraCustomResources.cachingUpstreams" . }}
{{ end }}