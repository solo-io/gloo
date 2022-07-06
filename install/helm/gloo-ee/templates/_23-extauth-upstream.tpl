{{- define "glooe.customResources.extauthUpstreams" -}}
{{- if .Values.global.extensions.extAuth.enabled }}
{{- $extAuth := .Values.global.extensions.extAuth }}
{{- $extAuthName := $extAuth.service.name }}

{{- if $extAuth.envoySidecar }}
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  name: {{ $extAuthName }}-sidecar
  namespace: {{ $.Release.Namespace }}
  labels:
    app: gloo
    gloo: {{ $extAuthName }}
spec:
  useHttp2: true
  pipe:
    path: "/usr/share/shared-data/.sock"
{{- end }}
---
{{- if $extAuth.standaloneDeployment }}
{{- include "gloo.dataplaneperproxyhelper" $ }}
{{- range $name, $spec := $.ProxiesToCreateDataplaneFor }}
{{- if not $spec.disabled }}
{{- if $.Values.global.extensions.dataplanePerProxy }}
{{- $extAuthName = printf "%s-%s" $extAuth.service.name ($name | kebabcase) }}
{{- end }}
apiVersion: gloo.solo.io/v1
kind: Upstream
metadata:
  name: {{ $extAuthName }}
  namespace: {{ $.Release.Namespace }}
  labels:
    app: gloo
    gloo: {{ $extAuthName }}
spec:
  useHttp2: true
  healthChecks:
  - timeout: 5s
    interval: 10s
    unhealthyThreshold: 3
    healthyThreshold: 3
    grpcHealthCheck:
      serviceName: {{ $extAuth.serviceName }}
  kube:
    serviceName: {{ $extAuthName }}
    serviceNamespace: {{ $.Release.Namespace }}
    servicePort:  {{ $extAuth.service.port }}
    serviceSpec:
      grpc: {}
  {{- if $.Values.global.glooMtls.enabled }}
  sslConfig:
    secretRef:
      name: gloo-mtls-certs
      namespace: {{ $.Release.Namespace }}
  {{- end }}
---
{{- end }}{{/* if not $spec.disabled */}}
{{- end }}{{/* range $name, $spec := $.ProxiesToCreateDataplaneFor */}}
{{- end }}{{/* $extAuth.standaloneDeployment */}}
{{- end }}{{/* .Values.global.extensions.extAuth.enabled */}}
{{- end }}{{/* define "glooe.customResources.extauthUpstreams" */}}
