{{/* vim: set filetype=mustache: */}}

{{/*
Render Gloo Edge Enterprise custom resources. This is done in a helper function, rather than included as a template,
to ensure that we don't try to apply the custom resources before the service backing the validation webhook
is ready.
When adding a new CR to the helm chart, the template should be included here.
*/}}
{{- define "glooe.customResources" -}}
{{- include "glooe.customResources.ratelimitUpstreams" . }}
{{- include "glooe.customResources.extauthUpstreams" . }}
{{- include "glooe.customResources.cachingUpstreams" . }}
{{ end }}