{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gloo-gateway.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Expand the name of the chart.
*/}}
{{- define "gloo-gateway.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "gloo-gateway.fullname" -}}
{{- .Release.Name | printf "glood-%s" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "gloo-gateway.labels" -}}
helm.sh/chart: {{ include "gloo-gateway.chart" . }}
app.kubernetes.io/name: {{ include "gloo-gateway.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Constant labels - labels that are stable across releases
We want this label to remain constant as it is used in glooctl version.
*/}}
{{- define "gloo-gateway.constLabels" -}}
gloo: kube-gateway
{{- end }}

{{/*
Images valid for the GatewayParameters
ref Image api in projects/gateway2/api/v1alpha1/kube/container.proto
*/}}
{{- define "gloo-gateway.gatewayParametersImage" -}}
{{- $image := . -}}
{{- if $image.registry }}
registry: {{ $image.registry }}
{{- end -}}{{/* if $image.registry */}}

{{- /* This has been copied from _helpers.tpl and should be kept in sync */ -}}
{{- if $image.repository }}
{{- $repository := $image.repository -}}
{{- /*
for fips or fips-distroless variants: add -fips to the image repo (name)
*/ -}}
{{- if or $image.fips (has $image.variant (list "fips" "fips-distroless")) -}}
{{- $fipsSupportedImages := list "gloo-ee" "extauth-ee" "gloo-ee-envoy-wrapper" "rate-limit-ee" "discovery-ee" "sds-ee" -}}
{{- if (has $image.repository $fipsSupportedImages) -}}
{{- $repository = printf "%s-fips" $repository -}}
{{- end -}}{{- /* if (has .repository $fipsSupportedImages) */ -}}
{{- end -}}{{- /* if or .fips (has .variant (list "fips" "fips-distroless")) */ -}}
{{ printf "\n" }}
repository: {{ $repository }}
{{- end -}}{{/* if $image.repository */}}

{{- if $image.tag }}
{{- $tag := $image.tag -}}
{{- /*
for distroless or fips-distroless variants: add -distroless to the tag
*/ -}}
{{- if has $image.variant (list "distroless" "fips-distroless") -}}
{{- $distrolessSupportedImages := list "gloo" "gloo-envoy-wrapper" "discovery" "sds" "certgen" "kubectl" "access-logger" "ingress" "gloo-ee" "extauth-ee" "gloo-ee-envoy-wrapper" "rate-limit-ee" "discovery-ee" "sds-ee" "observability-ee" "caching-ee" -}}
{{- if (has $image.repository $distrolessSupportedImages) -}}
{{- $tag = printf "%s-distroless" $tag -}} {{- /* Add distroless suffix to the tag since it contains the same binaries in a different container */ -}}
{{- end -}}{{- /* if (has .repository $distrolessSupportedImages) */ -}}
{{- end }}{{- /* if and .tag (has .variant (list "distroless" "fips-distroless")) */ -}}
{{ printf "\n" }}
tag: {{ $tag }}
{{- end -}}{{/* if $image.tag */}}
{{- if $image.digest }}
digest: {{ $image.digest }}
{{- end -}}{{/* if $image.digest */}}
{{- if $image.pullPolicy }}
pullPolicy: {{ $image.pullPolicy }}
{{- end -}}{{/* if $image.pullPolicy */}}
{{- end }}
