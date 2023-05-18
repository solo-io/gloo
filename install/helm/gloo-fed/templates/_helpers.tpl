{{/* vim: set filetype=mustache: */}}

{{/*
Validate a duration string.
*/}}
{{- define "gloofed.validateDuration" }}
{{- if not (kindIs "invalid" .) }}{{/* only check non-null/non-empty values */}}
{{- $dur := toString . }}{{/* converts input to string; needed to handle "0" which is a valid duration string */}}
{{- $_ := now | mustDateModify $dur }}{{/* try modifying the date `now` with the given duration; returns error on invalid duration */}}
{{- if hasPrefix "-" $dur }}{{/* don't allow negative values */}}
{{- fail (printf "invalid duration %s: must be positive" $dur) }}
{{- else }}
{{- . }}{{/* if everything's good, return the original value */}}
{{- end }}{{/* if hasPrefix "-" $dur */}}
{{- end }}{{/* if not (kindIs "invalid" .) */}}
{{- end }}


{{/* 
  gloofed.podSpecStandardFields is duplicated from the gloo codebase.
  TODO(sheidkamp) will look for a way to reuse/import those charts.
*/}}
{{- define "gloofed.podSpecStandardFields" -}}
{{- with .nodeName -}}
nodeName: {{ . }}
{{ end -}}
{{- with .nodeSelector -}}
nodeSelector: {{ toYaml . | nindent 2 }}
{{ end -}}
{{- with .tolerations -}}
tolerations: {{ toYaml . | nindent 2 }}
{{ end -}}
{{- with .hostAliases -}}
hostAliases: {{ toYaml . | nindent 2 }}
{{ end -}}
{{- with .affinity -}}
affinity: {{ toYaml . | nindent 2 }}
{{ end -}}
{{- with .restartPolicy -}}
restartPolicy: {{ . }}
{{ end -}}
{{- with .priorityClassName -}}
priorityClassName: {{ . }}
{{ end -}}
{{- with .initContainers -}}
initContainers: {{ toYaml . | nindent 2 }}
{{ end -}}
{{- end -}}



{{/* 
  gloofed.securityContext is duplicated from the gloo codebase.
  TODO(sheidkamp) will look for a way to reuse/import those charts.
*/}}
{{- /* 
This template is used to generate the gloo pod or container security context.
It takes 2 values:
  .values - the securityContext passed from the user in values.yaml
  .defaults - the default securityContext for the pod or container

  Depending upon the value of .values.merge, the securityContext will be merged with the defaults or completely replaced.
  In a merge, the values in .values will override the defaults, following the logic of helm's merge function.
Because of this, if a value is "true" in defaults it can not be modified with this method.
*/ -}}
{{- define "gloofed.securityContext" }}
{{- $securityContext := dict -}}
{{- $overwrite := true -}}
{{- if .values -}}
  {{- if .values.mergePolicy }}
    {{- if eq .values.mergePolicy "helm-merge" -}}
      {{- $overwrite = false -}}
    {{- else if ne .values.mergePolicy "no-merge" -}}
      {{- fail printf "value '%s' is not an allowed value for mergePolicy. Allowed values are 'no-merge', 'helm-merge', or an empty string" .values.mergePolicy }}
    {{- end -}}
  {{- end }}
{{- end -}}

{{- if $overwrite -}}
  {{- $securityContext = or .values .defaults (dict) -}}
{{- else -}}
  {{- $securityContext = merge .values .defaults }}
{{- end }}
{{- /* Remove "mergePolicy" if it exists because it is not a part of the kubernetes securityContext definition */ -}}
{{- $securityContext = omit $securityContext "mergePolicy" -}}
{{- with $securityContext -}}
securityContext:{{ toYaml . | nindent 2 }}
{{- end }}
{{- end }}