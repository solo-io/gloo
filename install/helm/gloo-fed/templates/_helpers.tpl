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
