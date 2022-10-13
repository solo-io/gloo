package placement

import (
	"io/ioutil"

	"github.com/solo-io/skv2/codegen/model"
	"github.com/solo-io/skv2/codegen/util"
	"github.com/solo-io/skv2/contrib"
)

const (
	typedParserCustomTemplatePath = "parser.gotmpl"
)

// Returns the skv2 code template for generating typed placement parsers.
func TypedParser(params contrib.SnapshotTemplateParameters) model.CustomTemplates {
	templateContentsBytes, err := ioutil.ReadFile(util.MustGetThisDir() + "/" + typedParserCustomTemplatePath)
	if err != nil {
		panic(err)
	}
	templateContents := string(templateContentsBytes)
	return params.ConstructTemplate(params, templateContents, false)
}
