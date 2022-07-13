package validation_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/translation"
)

func TestGraphqlValidationUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Graphql Validation Utils Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	err := os.Setenv(translation.GraphqlJsRootEnvVar, "../../../plugins/graphql/js/")
	Expect(err).NotTo(HaveOccurred())
	err = os.Setenv(translation.GraphqlProtoRootEnvVar, "../../../../../ui/src/proto/")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := os.Unsetenv(translation.GraphqlProtoRootEnvVar)
	Expect(err).NotTo(HaveOccurred())
	err = os.Unsetenv(translation.GraphqlJsRootEnvVar)
	Expect(err).NotTo(HaveOccurred())
})
