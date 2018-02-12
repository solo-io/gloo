package aws

import (
	"github.com/pkg/errors"
	"github.com/solo-io/glue/pkg/api/types/v1"
)

const (
	// TODO: execution_type should be moved to a separate package
	RoutePluginKeyExecutionStyle = "execution_style"
)

type ExecutionStyle string

const (
	ExecutionStyleSync  ExecutionStyle = "sync"
	ExecutionStyleAsync ExecutionStyle = "async"
)

func GetExecutionStyle(routeSpec v1.RoutePluginSpec) (ExecutionStyle, error) {
	if style, ok := routeSpec[RoutePluginKeyExecutionStyle]; ok {
		executionStyle, ok := style.(ExecutionStyle)
		if !ok {
			return "", errors.Errorf("invalid format for %v, expected string", RoutePluginKeyExecutionStyle)
		}
		switch executionStyle {
		case ExecutionStyleSync:
		case ExecutionStyleAsync:
		default:
			return "", errors.Errorf("invalid execution style provided. must be one of: %v|%v, you gave: %v", ExecutionStyleSync, ExecutionStyleAsync, executionStyle)
		}
		return executionStyle, nil
	}
	return ExecutionStyleSync, nil
}
