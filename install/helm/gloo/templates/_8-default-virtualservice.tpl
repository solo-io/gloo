{{- define "gloo.customResources.defaultVirtualService" -}}
{{- if .Values.virtualService.enabled }}
---
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: {{ .Values.virtualService.name | kebabcase }}
  namespace: {{ .Release.Namespace }}
spec:
  virtualHost:
    domains:
    {{- if .Values.virtualService.domains }}
    {{- range $domain := .Values.virtualService.domains }}
    - "{{ $domain }}"
    {{- end }}{{/* range $domain := .Values.virtualService.domains */}}
    {{- else}}
    - "*"
    {{- end }}{{/* if .Values.virtualService.domains */}}
    routes:
    - matchers:
    {{- if .Values.virtualService.routes }}
    {{- range $match := .Values.virtualService.routes }}
      - exact: "{{ $match.path }}"
      directResponseAction:
        status: 200
        body: "{{ $match.response }}"
    {{- end }}{{/* range $path, $response := .Values.virtualService.matchers */}}
    {{- else}}
      - exact: "/"
      directResponseAction:
        status: 200
        body: "Welcome to Gloo Edge!"
  {{- end }}{{/* if .Values.virtualService.matchers */}}
{{- end }}{{/* if .Values.virtualService.enabled */}}
{{- end }}{{/* define "gloo.customResources.defaultVirtualService" */}}

