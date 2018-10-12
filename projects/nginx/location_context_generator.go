package nginx

// TODO(talnordan): Can `.Prefix` be optional?
const locationContextTemplateText = `
location {{.Prefix}} {
{{- if .Root}}
    root {{.Root}};
{{- end}}
{{- if .ProxyPass}}
    proxy_pass {{.ProxyPass}};
{{- end}}
}
`

func GenerateLocationContext(location *Location) ([]byte, error) {
	// TODO(talnordan): Consider verifying that either `.Root` is defined or `.ProxyPass` is, but
	// not both.
	return generateContext(locationContextTemplateText, location)
}
