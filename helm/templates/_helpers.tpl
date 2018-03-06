{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart.
*/}}
{{- define "gloo-chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "gloo-chart.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gloo-chart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* helpers for components */}}

{{/* function discovery */}}
{{- define "function_discovery.name" -}}
{{- default "function-discovery" .Values.function_discovery.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "function_discovery.fullname" -}}
{{- $name := default "function-discovery" .Values.function_discovery.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "function_discovery.serviceAccountName" -}}
{{- if .Values.global.rbacEnabled -}}
{{- template "function_discovery.fullname" . -}}-service-account
{{- else }}
{{- .Values.function_discovery.serviceAccountName | trunc 63 | trimSuffix "-" -}}-service-account
{{- end -}}
{{- end -}}

{{/* ingress */}}
{{- define "ingress.name" -}}
{{- default "ingress" .Values.ingress.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "ingress.fullname" -}}
{{- $name := default "ingress" .Values.ingress.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* gloo */}}
{{- define "gloo.name" -}}
{{- default "gloo" .Values.gloo.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "gloo.fullname" -}}
{{- $name := default "gloo" .Values.gloo.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}