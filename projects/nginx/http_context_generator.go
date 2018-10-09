package nginx

import (
	"bytes"
	"html/template"
)

const httpContextTemplateText = `
http {
{{- if .Server}}
	server {
	}
{{- end}}
}
`

func GenerateHttpContext(http *Http) ([]byte, error) {
	tmpl, err := template.New("HTTP context").Parse(httpContextTemplateText)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, http)
	return buffer.Bytes(), nil
}
