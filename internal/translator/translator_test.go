package translator_test

import (
	"github.com/solo-io/glue/internal/plugins/aws"
	. "github.com/solo-io/glue/internal/translator"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/pkg/plugin"
	"github.com/solo-io/glue/pkg/secretwatcher"
	. "github.com/solo-io/glue/test/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Translator", func() {
	It("works", func() {
		t := NewTranslator([]plugin.TranslatorPlugin{&aws.Plugin{}})
		cfg := NewTestConfig()
		snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
		Expect(err).NotTo(HaveOccurred())
		log.Printf("%v", snap)
		log.Printf("%v", reports)
	})
})
