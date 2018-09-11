{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "gloo-ee.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "gloo-ee.fullname" -}}
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
{{- define "gloo-ee.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/* gloo */}}
{{- define "control_plane.name" -}}
{{- default "control-plane" .Values.control_plane.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "control_plane.fullname" -}}
{{- $name := default "control-plane" .Values.control_plane.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/* data-plane */}}
{{- define "data_plane.name" -}}
{{- default "data-plane" .Values.data_plane.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "data_plane.fullname" -}}
{{- $name := default "data-plane" .Values.data_plane.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/* discovery */}}
{{- define "discovery.name" -}}
{{- default "discovery" .Values.discovery.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "discovery.fullname" -}}
{{- $name := default "discovery" .Values.discovery.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{/* gloo_i */}}
{{- define "gloo_i.name" -}}
{{- default "gloo-i" .Values.gloo_i.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "gloo_i.fullname" -}}
{{- $name := default "gloo-i" .Values.gloo_i.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
