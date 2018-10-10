package nginx

import (
	"bytes"
	"html/template"
)

// TODO(talnordan): Can `.Prefix` or `.Root` be optional?
const locationContextTemplateText = `
location {{.Prefix}} {
    root {{.Root}};
}
`

// TODO(talnordan): Deduplicate common code with `GenerateHttpContext()`
func GenerateLocationContext(location *Location) ([]byte, error) {
	// TODO(talnordan): Does the template name matter?
	tmpl, err := template.New("location context").Parse(locationContextTemplateText)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, location)
	return buffer.Bytes(), nil
}
