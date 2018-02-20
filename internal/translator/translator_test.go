package translator_test

import (
	"github.com/solo-io/gloo-plugins/aws"
	. "github.com/solo-io/gloo-testing/helpers"
	. "github.com/solo-io/gloo/internal/translator"
	"github.com/solo-io/gloo/pkg/endpointdiscovery"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/plugin"
	"github.com/solo-io/gloo/pkg/secretwatcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Translator", func() {
	It("works", func() {
		t := NewTranslator([]plugin.TranslatorPlugin{&aws.Plugin{}})
		cfg := NewTestConfig()
		snap, reports, err := t.Translate(cfg, secretwatcher.SecretMap{}, endpointdiscovery.EndpointGroups{})
		Expect(err).NotTo(HaveOccurred())
		log.Debugf("%v", snap)
		log.Debugf("%v", reports)
	})
})
