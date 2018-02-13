package openapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"unicode"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"

	"github.com/solo-io/glue-discovery/pkg/source"
)

func GetOpenAPIFetcher() *openapiFetcher {
	return &openapiFetcher{}
}

type openapiFetcher struct{}

func (f *openapiFetcher) Fetch(u *source.Upstream) ([]source.Function, error) {
	oURL := openapiURL(u)
	if oURL == nil {
		return nil, fmt.Errorf("unable to find Open API Docs URL for %s", u.ID)
	}

	spec, err := loadSwaggerSpec(oURL)
	if err != nil {
		return nil, err
	}

	return parse(spec)
}

func (f *openapiFetcher) CanFetch(u *source.Upstream) bool {
	return openapiURL(u) != nil
}

func openapiURL(u *source.Upstream) *url.URL {
	v, exists := u.Spec["openapi"]
	if !exists {
		return nil
	}
	oURL, err := url.Parse(v.(string))
	if err != nil {
		log.Printf("Not a valid URL for Open API Docs %v", v)
		return nil
	}
	return oURL
}

func loadSwaggerSpec(u *url.URL) (*spec.Swagger, error) {
	var data json.RawMessage
	urlPath := strings.ToLower(u.Path)
	if strings.HasSuffix(urlPath, ".yaml") || strings.HasSuffix(urlPath, ".yml") {
		d, err := swag.YAMLDoc(u.String())
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load Open API YAML from %v", u)
		}
		data = d
	} else {
		d, err := swag.LoadFromFileOrHTTP(u.String())
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load Open API document frm %v ", u)
		}
		data = json.RawMessage(d)
	}

	doc, err := loads.Analyzed(data, "")
	if err != nil {
		return nil, errors.Wrapf(err, "unable to analyze Open API document for %v", u)
	}

	return doc.Spec(), nil
}

func parse(s *spec.Swagger) ([]source.Function, error) {
	basePath := s.BasePath
	version := s.Info.Version
	protocol := ""
	if len(s.Schemes) > 0 {
		protocol = s.Schemes[0]
	}

	var functions []source.Function
	for path, v := range s.Paths.Paths {
		ops := map[*spec.Operation]string{
			v.Delete:  "delete",
			v.Get:     "get",
			v.Head:    "head",
			v.Options: "options",
			v.Patch:   "patch",
			v.Post:    "post",
			v.Put:     "put",
		}
		for op, lbl := range ops {
			if op == nil {
				continue
			}
			f := source.Function{
				Name: functionName(op, path, lbl),
				Spec: map[string]interface{}{
					"HTTPMethod":  strings.ToUpper(lbl),
					"Protocol":    protocol,
					"Description": op.Description,
					"Version":     version,
					"URITemplate": basePath + path,
				},
			}
			functions = append(functions, f)
		}
	}

	return functions, nil
}

func functionName(op *spec.Operation, path, label string) string {
	if op.ID != "" {
		return op.ID
	}

	var resource []rune
	upperNext := false
	for _, c := range path {
		if c == '/' || c == '{' || c == '}' {
			upperNext = true
			continue
		}
		if upperNext {
			upperNext = false
			resource = append(resource, unicode.ToUpper(c))
		} else {
			resource = append(resource, c)
		}

	}
	return label + string(resource)
}
