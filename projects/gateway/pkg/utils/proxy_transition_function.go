package utils

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/hashutils"
)

func TransitionFunction(original, desired *v1.Proxy) (bool, error) {
	equal, ok := hashutils.HashableEqual(original, desired)
	if ok && equal {
		desired.Status = original.Status
		return true, nil
	}
	return false, nil
}