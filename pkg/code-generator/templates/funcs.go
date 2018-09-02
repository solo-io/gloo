package templates

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var funcs = template.FuncMap{
	"join":        strings.Join,
	"lowercase":   strings.ToLower,
	"lower_camel": strcase.ToLowerCamel,
	"upper_camel": strcase.ToCamel,
	"snake":       strcase.ToSnake,
}
