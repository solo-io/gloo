package nginx

import (
	"bytes"
	"html/template"

	"github.com/Masterminds/sprig"
)

// TODO(talnordan): Consider changing the return type to `(string, error)`
func generateContext(templateText string, data interface{}) ([]byte, error) {
	// TODO(talnordan): Does the template name matter?
	tmpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(templateText)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	return buffer.Bytes(), nil
}
