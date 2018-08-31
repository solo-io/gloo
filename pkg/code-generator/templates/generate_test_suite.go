package templates

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
