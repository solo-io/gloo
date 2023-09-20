{{- define "gloo.customResources.defaultVirtualService" -}}
{{- if .Values.virtualService.enabled }}
---
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: {{ .Values.virtualService.name | kebabcase }}
  namespace: {{ .Release.Namespace }}
spec:
  virtualHost: {}
{{- end }}{{/* if .Values.virtualService.enabled */}}
{{- end }}{{/* define "gloo.customResources.defaultVirtualService" */}}

