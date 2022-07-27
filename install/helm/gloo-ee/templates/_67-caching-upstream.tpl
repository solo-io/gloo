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
    interval: 1m
    unhealthyThreshold: 5
    healthyThreshold: 5
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

{{- define "glooe.customResources.cachingUpstreams" -}}
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
{{- end }}{{/* define "glooe.customResources.cachingUpstreams" */}}
