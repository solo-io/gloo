package typed

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

type ResourceLevelTemplateParams struct {
	PackageName           string
	ResourceType          string
	IsDataType            bool
	IsInputType           bool
	ResourceTypeLowerCase string
	PluralName            string
	GroupName             string
	Version               string
	ShortName             string
	Fields                []string
}

type Field struct {
	Name   string
	Type   string
	Fields []Field
}

type PackageLevelTemplateParams struct {
	PackageName         string
	ResourceTypes       []string
	ResourceLevelParams map[string]*ResourceLevelTemplateParams
}

var funcs = template.FuncMap{
	"join":      strings.Join,
	"lowercase": strcase.ToLowerCamel,
	"uppercase": strcase.ToCamel,
	"clients": func(params PackageLevelTemplateParams, withTypes bool) string {
		var clientParams []string
		for _, resource := range params.ResourceTypes {
			paramString := strcase.ToLowerCamel(resource) + "Client"
			if withTypes {
				paramString += " " + resource + "Client"
			}
			clientParams = append(clientParams, paramString)
		}
		return strings.Join(clientParams, ", ")
	},
	"need_clientset": func(params PackageLevelTemplateParams) bool {
		for _, p := range params.ResourceLevelParams {
			if !p.IsInputType && p.IsDataType {
				return true
			}
		}
		return false
	},
	"need_memory_client": func(params PackageLevelTemplateParams) bool {
		for _, p := range params.ResourceLevelParams {
			if !p.IsInputType && !p.IsDataType {
				return true
			}
		}
		return false
	},
}
