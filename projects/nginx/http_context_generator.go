package nginx

const httpContextTemplateText = `http {
{{- if .Server}}
	server {
	}
{{- end}}
}`

func GenerateHttpContext(http *Http) ([]byte, error) {
	return generateContext(httpContextTemplateText, http)
}
