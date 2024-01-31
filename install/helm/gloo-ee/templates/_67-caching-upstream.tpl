{{- define "caching.upstreamSpec" -}}
{{- $name := (index . 1) }}
{{- with (first .) }}
{{- $cachingName := .Values.global.extensions.caching.name }}
{{- if .Values.global.extensions.dataplanePerProxy }}
{{- $cachingName = printf "%s-%s" $cachingName ($name | kebabcase) }}
{{- end }}{{/* .Values.global.extensions.dataplanePerProxy */}}
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  name: {{ $cachingName }}
  namespace: {{ .Release.Namespace }}
  labels:
    app: gloo
    gloo: {{ $cachingName }}
spec:
  healthChecks:
  - timeout: 5s
    interval: 10s
    noTrafficInterval: 10s
    unhealthyThreshold: 3
    healthyThreshold: 3
    grpcHealthCheck:
      serviceName: caching-service
  kube:
    serviceName: {{ $cachingName }}
    serviceNamespace: {{ .Release.Namespace }}
    servicePort:  {{ .Values.global.extensions.caching.service.httpPort }}
    serviceSpec:
      grpc: {}
---
{{- end }}{{/* with (first .) */}}
{{- end }}{{/* define "caching.upstreamSpec" */}}

{{/*
As this gets the values context from OSS, $ will refer to the OSS values context.
If additional fields are required, add them at https://github.com/solo-io/gloo/blob/0429470a3f671b1137b36abe105f5df3d583d53f/install/helm/gloo/templates/5-resource-configmap.yaml#L17
*/}}
{{- define "gloo.extraCustomResources.cachingUpstreams" -}}
{{- if .Values.global.extensions.caching.enabled }}
{{- include "gloo.dataplaneperproxyhelper" $ }}
{{- $override := dict -}}
{{- if .Values.global.extensions.caching.upstream }}
{{- $override = .Values.global.extensions.caching.upstream.kubeResourceOverride}}
{{- end }}{{/* if .Values.global.extensions.caching.upstream */}}
{{- range $name, $spec := $.ProxiesToCreateDataplaneFor }}
{{- if not $spec.disabled}}
{{- $ctx := (list $ $name $spec)}}
{{- include "gloo.util.merge" (list $ctx $override "caching.upstreamSpec") -}}
{{- end }}{{/* if not $spec.disabled */}}
{{- end }}{{/* range $name, $spec := $.ProxiesToCreateDataplaneFor */}}
{{- end }}{{/* .Values.global.extensions.caching.enabled */}}
{{- end }}{{/* define "gloo.extraCustomResources.cachingUpstreams" */}}
