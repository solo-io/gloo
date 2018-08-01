package rbac

import (
	"fmt"

	"github.com/ory/ladon"
	"github.com/solo-io/solo-kit/pkg/api/v1/rbac/policy"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

func GeneratePolicy(pol *policy.Policy) ladon.Policy {
	return &ladon.DefaultPolicy{
		ID:          id(pol.Metadata),
		Description: pol.Description,
		Subjects:    pol.Subjects,
		Effect:      ladon.AllowAccess,
		Resources:   pol.Resources,
		Actions:     actions(pol.Capabilities),
		Conditions:  conditions(pol),
		Meta:        pol.Meta,
	}
}

func id(meta core.Metadata) string {
	return fmt.Sprintf("%v.%v", meta.Namespace, meta.Name)
}

func actions(caps []policy.Capability) []string {
	var actions []string
	for _, cap := range caps {
		actions = append(actions, cap.String())
	}
	return actions
}

func conditions(pol *policy.Policy) ladon.Conditions {
	return ladon.Conditions{
		"": &ladon.EqualsSubjectCondition{},
	}
}
