package nginx

import (
	"bytes"
	"html/template"
)

func generateContext(templateText string, data interface{}) ([]byte, error) {
	// TODO(talnordan): Does the template name matter?
	tmpl, err := template.New("").Parse(templateText)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	return buffer.Bytes(), nil
}
