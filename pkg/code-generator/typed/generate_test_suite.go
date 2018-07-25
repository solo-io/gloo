package typed

import (
	"bytes"
	"text/template"
)

func GenerateTestSuiteCode(params PackageLevelTemplateParams) (string, error) {
	buf := &bytes.Buffer{}
	if err := testSuiteTemplate.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var testSuiteTemplate = template.Must(template.New("typed_client_test_suite").Funcs(funcs).Parse(testSuiteTemplateContents))

const testSuiteTemplateContents = `package {{ .PackageName }}

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test{{ join .ResourceTypes "" }}(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "{{ join .ResourceTypes "" }} Suite")
}
`
