package swagger

import (
	"strings"

	"fmt"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/internal/swagger"
	"github.com/solo-io/gloo-plugins/transformation"
	"github.com/solo-io/gloo/pkg/log"
)

func GetFuncs(us *v1.Upstream) ([]*v1.Function, error) {
	swaggerSpec, err := getSwaggerSpecForUpsrteam(us)
	if err != nil {
		return nil, err
	}
	var consumesJson bool
	for _, contentType := range swaggerSpec.Consumes {
		if contentType == "application/json" {
			consumesJson = true
			break
		}
	}
	if !consumesJson {
		return nil, errors.Errorf("swagger function discovery uses content type application/json; "+
			"available: %v", swaggerSpec.Consumes)
	}
	// TODO: when response transformation is done, look at produces as well
	var funcs []*v1.Function
	for functionPath, pathItem := range swaggerSpec.Paths.Paths {
		funcs = append(funcs, createFunctionsForPath(swaggerSpec.BasePath, functionPath, pathItem.PathItemProps)...)
	}
	return funcs, nil
}

func createFunctionsForPath(basePath, functionPath string, path spec.PathItemProps) []*v1.Function {
	var pathFunctions []*v1.Function
	appendFunction := func(method string, operation *spec.Operation) {
		pathFunctions = append(pathFunctions, createFunctionForOpertaion(method, basePath, functionPath, operation.OperationProps))
	}
	if path.Get != nil {
		appendFunction("GET", path.Get)
	}
	if path.Put != nil {
		appendFunction("PUT", path.Put)
	}
	if path.Post != nil {
		appendFunction("POST", path.Post)
	}
	if path.Delete != nil {
		appendFunction("DELETE", path.Delete)
	}
	if path.Options != nil {
		appendFunction("OPTIONS", path.Options)
	}
	if path.Head != nil {
		appendFunction("HEAD", path.Head)
	}
	if path.Patch != nil {
		appendFunction("PATCH", path.Patch)
	}
	return pathFunctions
}

func createFunctionForOpertaion(method string, basePath, functionPath string, operation spec.OperationProps) *v1.Function {
	var queryParams, headerParams []string
	bodyParams := make(map[string]spec.SchemaProps)
	for _, param := range operation.Parameters {
		// sort parameters by the template they will go into
		switch param.In {
		case "query":
			queryParams = append(queryParams, fmt.Sprintf("%v={%v}", param.Name, param.Name))
		case "header":
			headerParams = append(headerParams, param.Name)
		case "path":
			// nothing to do here, we already get the template
		case "formData":
			log.Warnf("form data params not currently supported; ignoring")
		case "body":
			bodyParams[param.Name] = param.Schema.SchemaProps
		}
	}

	path := swaggerPathToJinjaTemplate(basePath + functionPath)
	if len(queryParams) > 0 {
		path += "?" + strings.Join(queryParams, "&")
	}

	headersTemplate := map[string]string{":method": method}
	for _, name := range headerParams {
		headersTemplate[name] = fmt.Sprintf("{{%v}}", name)
	}

	return &v1.Function{
		Name: operation.ID,
		Spec: transformation.EncodeFunctionSpec(transformation.FunctionSpec{
			Path:   path,
			Header: headersTemplate,
			Body:   constructBodyTemplate(bodyParams),
		}),
	}
}

// TODO: make the body template constructor much more sophistocated
// right now it's only supporting primitive fields (not nested objects)
func constructBodyTemplate(bodyParams map[string]spec.SchemaProps) string {
	bodyTemplate := "{"
	var fields []string
	for name, schema := range bodyParams {
		if schema.Type.Contains("string") {
			// if it's a string, need to do quotes
			fields = append(fields, fmt.Sprintf(`"%v": "{{%v}}"`, name, name))
		} else {
			fields = append(fields, fmt.Sprintf(`"%v": {{%v}}`, name, name))
		}
	}
	bodyTemplate += "}"
	return bodyTemplate
}

func swaggerPathToJinjaTemplate(path string) string {
	path = strings.Replace(path, "{", "{{", -1)
	path = strings.Replace(path, "}", "}}", -1)
	return path
}

func getSwaggerSpecForUpsrteam(us *v1.Upstream) (*spec.Swagger, error) {
	annotations, err := swagger.GetSwaggerAnnotations(us)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid or missing swagger annotations on %v", us.Name)
	}
	switch {
	case annotations.SwaggerURL != "":
		return retrieveSwaggerDocFromUrl(annotations.SwaggerURL)
	case annotations.InlineSwaggerDoc != "":
		return parseSwaggerDoc([]byte(annotations.InlineSwaggerDoc))
	}
	return nil, errors.Errorf("one of %v or %v must be specified on the swagger upstream annotations",
		swagger.AnnotationKeySwaggerDoc,
		swagger.AnnotationKeySwaggerURL)
}

func retrieveSwaggerDocFromUrl(url string) (*spec.Swagger, error) {
	docBytes, err := swag.LoadFromFileOrHTTP(url)
	if err != nil {
		return nil, errors.Wrap(err, "loading swagger doc from url")
	}
	return parseSwaggerDoc(docBytes)
}

func parseSwaggerDoc(docBytes []byte) (*spec.Swagger, error) {
	doc, err := loads.Analyzed(docBytes, "")
	if err != nil {
		log.Warnf("parsing doc as json failed, falling back to yaml")
		jsn, err := swag.YAMLToJSON(docBytes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert yaml to json (after falling back to yaml parsing)")
		}
		doc, err = loads.Analyzed(jsn, "")
		if err != nil {
			return nil, errors.Wrap(err, "invalid swagger doc")
		}
	}
	return doc.Spec(), nil
}
