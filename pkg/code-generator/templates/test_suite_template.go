package templates

import (
	"text/template"
)

var ProjectTestSuiteTemplate = template.Must(template.New("project_template").Funcs(funcs).Parse(`package {{ .Version }}

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test{{ upper_camel .Name }}(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "{{ upper_camel .Name }} Suite")
}





`))
