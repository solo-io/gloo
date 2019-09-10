{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}

{{/* Used to update the values during the render of a template. Useful for taking user-friendly gloo-ee
     values and renaming them to gloo's expected format without leaking implementation details */}}
{{- define "gloo.updatevalues" -}}
{{- if .Values.global.extensions.extAuth.envoySidecar -}}
{{- range $proxyName, $proxy := .Values.gatewayProxies -}}
{{- $_ := set (index $.Values.gatewayProxies $proxyName) "extraContainersHelper" "gloo.extauthcontainer" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/* Container definition for extauth, used in extauth deployment and
     gateway-proxy (envoy) sidecar over unix domain socket */}}
{{- define "gloo.extauthcontainer" -}}
{{- $extAuth := .Values.global.extensions.extAuth -}}
{{- $image := $extAuth.deployment.image -}}
{{- if .Values.global -}}
{{- $image = merge $extAuth.deployment.image .Values.global.image -}}
{{- end -}}
- image: {{template "gloo.image" $image}}
  imagePullPolicy: {{ $image.pullPolicy }}
  name: {{ $extAuth.deployment.name }}
  env:
    - name: POD_NAMESPACE
      valueFrom:
        fieldRef:
          fieldPath: metadata.namespace
    - name: SERVICE_NAME
      value: {{ $extAuth.serviceName  | quote }}
    - name: GLOO_ADDRESS
{{- if $extAuth.deployment.glooAddress }}
      value: {{ $extAuth.deployment.glooAddress }}
{{- else }}
      {{- if .Values.gloo.gloo }}
      value: gloo:{{ .Values.gloo.gloo.deployment.xdsPort }}
      {{- else }}
      value: gloo:{{ .Values.gloo.deployment.xdsPort }}
      {{- end }}
{{- end }}
    - name: SIGNING_KEY
      valueFrom:
        secretKeyRef:
          name: {{ $extAuth.signingKey.name }}
          key: signing-key
    {{- if $extAuth.deployment.debugPort }}
    - name: DEBUG_PORT
      value: {{ $extAuth.deployment.debugPort | quote }}
    {{- end }}
    {{- if $extAuth.deployment.port }}
    - name: SERVER_PORT
      value: {{ $extAuth.deployment.port  | quote }}
    {{- end }}
    {{- if $extAuth.envoySidecar }}
    - name: UDS_ADDR
      value: "/usr/share/shared-data/.sock"
    {{- end }}
    {{- if $extAuth.userIdHeader }}
    - name: USER_ID_HEADER
      value: {{ $extAuth.userIdHeader  | quote }}
    {{- end }}
    {{- if $extAuth.deployment.stats }}
    - name: START_STATS_SERVER
      value: "true"
    {{- end}}
  {{- if or $extAuth.plugins $extAuth.envoySidecar }}
  volumeMounts:
  {{- if $extAuth.envoySidecar }}
  - name: shared-data
    mountPath: /usr/share/shared-data
  {{- end }}
  {{- if $extAuth.plugins }}
  - name: auth-plugins
    mountPath: /auth-plugins
  {{- end }}
  {{- end }}
{{- end }}