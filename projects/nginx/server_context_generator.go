package nginx

// TODO(talnordan): Make location prefix and root configurable.
const serverContextTemplateText = `
server {
{{- if .Location}}
    location / {
        root /data/www;
    }
{{- end}}
}
`

func GenerateServerContext(server *Server) ([]byte, error) {
	return generateContext(serverContextTemplateText, server)
}
