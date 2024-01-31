{{- define "ratelimit.upstreamSpec" -}}
{{- $name := (index . 1) }}
{{- with (first .) }}
{{- $rateLimitName := .Values.global.extensions.rateLimit.service.name }}
{{- if .Values.global.extensions.dataplanePerProxy }}
{{- $rateLimitName = printf "%s-%s" $rateLimitName ($name | kebabcase) }}
{{- end }}{{/* .Values.global.extensions.dataplanePerProxy */}}
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  name: {{ $rateLimitName }}
  namespace: {{ .Release.Namespace }}
  labels:
    app: gloo
    gloo: {{ $rateLimitName }}
spec:
  healthChecks:
  - timeout: 5s
    interval: 10s
    noTrafficInterval: 10s
    unhealthyThreshold: 3
    healthyThreshold: 3
    grpcHealthCheck:
      serviceName: ratelimit
  kube:
    serviceName: {{ $rateLimitName }}
    serviceNamespace: {{ .Release.Namespace }}
    servicePort:  {{ .Values.global.extensions.rateLimit.service.port }}
    serviceSpec:
      grpc: {}
---
{{- end }}{{/* with (first .) */}}
{{- end }}{{/* define "ratelimit.upstreamSpec" */}}

{{/*
As this gets the values context from OSS, $ will refer to the OSS values context.
If additional fields are required, add them at https://github.com/solo-io/gloo/blob/0429470a3f671b1137b36abe105f5df3d583d53f/install/helm/gloo/templates/5-resource-configmap.yaml#L17
*/}}
{{- define "gloo.extraCustomResources.ratelimitUpstreams" -}}
{{- if .Values.global.extensions.rateLimit.enabled }}
{{- include "gloo.dataplaneperproxyhelper" $ }}
{{- $override := dict -}}
{{- if .Values.global.extensions.rateLimit.upstream }}
{{- $override = .Values.global.extensions.rateLimit.upstream.kubeResourceOverride}}
{{- end }}{{/* if .Values.global.extensions.rateLimit.upstream */}}
{{- range $name, $spec := $.ProxiesToCreateDataplaneFor }}
{{- if not $spec.disabled}}
{{- $ctx := (list $ $name $spec)}}
{{ include "gloo.util.merge" (list $ctx $override "ratelimit.upstreamSpec") -}}
{{- end }}{{/* if not $spec.disabled */}}
{{- end }}{{/* range $name, $spec := $.ProxiesToCreateDataplaneFor */}}
{{- end }}{{/* .Values.global.extensions.rateLimit.enabled */}}
{{- end }}{{/* define "gloo.extraCustomResources.ratelimitUpstreams" */}}
