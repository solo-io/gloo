package nginx

// TODO(talnordan): Consider making context string be generated with no trailing newline to begin
// with.
const serverContextTemplateText = `server {
{{- if .Locations}}
{{- range .Locations}}
{{.ContextString | indent 4}}
{{- end}}
{{- end}}
}`

func GenerateServerContext(server *Server) ([]byte, error) {
	return generateContext(serverContextTemplateText, server)
}
