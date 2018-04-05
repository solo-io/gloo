package aws

import (
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

const (
	// TODO: execution_type should be moved to a separate package
	RoutePluginKeyExecutionStyle = "execution_style"
)

const (
	ExecutionStyleNone  = "none"
	ExecutionStyleSync  = "sync"
	ExecutionStyleAsync = "async"
)

func GetExecutionStyle(routeExtensions *types.Struct) (string, error) {
	if routeExtensions != nil {
		if style, ok := routeExtensions.Fields[RoutePluginKeyExecutionStyle]; ok {
			executionStyleValue, ok := style.Kind.(*types.Value_StringValue)
			if !ok {
				return "", errors.Errorf("invalid format for %v, expected string", RoutePluginKeyExecutionStyle)
			}
			executionStyle := executionStyleValue.StringValue
			switch executionStyle {
			case ExecutionStyleSync:
			case ExecutionStyleAsync:
			default:
				return "", errors.Errorf("invalid execution style provided. must be one of: %v|%v, you gave: %v", ExecutionStyleSync, ExecutionStyleAsync, executionStyle)
			}
			return executionStyle, nil
		} else {
			return ExecutionStyleNone, nil
		}
	}
	return ExecutionStyleSync, nil
}
