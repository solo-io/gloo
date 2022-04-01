package directives

import (
	"strconv"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/rotisserie/eris"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
)

const (
	CACHE_CONTROL_DIRECTIVE               = "cacheControl"
	CACHE_CONTROL_MAXAGE_ARGUMENT         = "maxAge"
	CACHE_CONTROL_INHERIT_MAXAGE_ARGUMENT = "inheritMaxAge"
	CACHE_CONTROL_SCOPE_ARGUMENT          = "scope"
)

type cacheControlDirective struct {
	CacheControl *v2.CacheControl
}

func NewCacheControlDirective() *cacheControlDirective {
	return &cacheControlDirective{}
}

// Validate validates that the cacheControl directive usage is syntactically correct.
// At the end of successful validation, `CacheControl` will be populated.
func (d *cacheControlDirective) Validate(directiveVisitorParams DirectiveVisitorParams) (bool, error) {
	arguments := map[string]ast.Value{}
	directive := directiveVisitorParams.Directive
	for _, argument := range directive.Arguments {
		arguments[argument.Name.Value] = argument.Value
	}
	maxAge, maxAgeFound := arguments[CACHE_CONTROL_MAXAGE_ARGUMENT]
	inheritMaxAge, inheritMaxAgeFound := arguments[CACHE_CONTROL_INHERIT_MAXAGE_ARGUMENT]
	scope, scopeFound := arguments[CACHE_CONTROL_SCOPE_ARGUMENT]

	cacheControl := &v2.CacheControl{}
	if maxAgeFound {
		if maxAge.GetKind() != kinds.IntValue {
			return false, NewGraphqlSchemaError(maxAge, `"%s" argument must be an integer value`, CACHE_CONTROL_MAXAGE_ARGUMENT)
		}
		uintMaxAge, err := strconv.ParseUint(maxAge.GetValue().(string), 10, 32)
		if err != nil {
			return false, err
		}
		cacheControl.MaxAge = &wrappers.UInt32Value{Value: uint32(uintMaxAge)}
	}
	if inheritMaxAgeFound {
		if inheritMaxAge.GetKind() != kinds.BooleanValue {
			return false, NewGraphqlSchemaError(maxAge, `"%s" argument must be a boolean value`, CACHE_CONTROL_INHERIT_MAXAGE_ARGUMENT)
		}
		cacheControl.InheritMaxAge = inheritMaxAge.GetValue().(bool)
	}
	if scopeFound {
		if scope.GetKind() != kinds.EnumValue {
			return false, NewGraphqlSchemaError(maxAge, `"%s" argument must be a enum value`, CACHE_CONTROL_SCOPE_ARGUMENT)
		}
		scopeStr := scope.GetValue().(string)
		scope := v2.CacheControl_UNSET
		switch scopeStr {
		case "unset":
			scope = v2.CacheControl_UNSET
		case "public":
			scope = v2.CacheControl_PUBLIC
		case "private":
			scope = v2.CacheControl_PRIVATE
		default:
			return false, eris.Errorf("unimplemented cacheControl scope type %s", scopeStr)
		}
		cacheControl.Scope = scope
	}

	d.CacheControl = cacheControl
	return false, nil
}
