package validation

import (
	"context"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/stringutils"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	internal_placement "github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/internal/placement"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/rbac"
	"go.uber.org/zap"
	admission_v1 "k8s.io/api/admission/v1"
	authv1 "k8s.io/api/authentication/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	PlacementParsingError = func(err error, req *admission.Request) error {
		return eris.Wrapf(err, "failed to parse placement for resource %s.%s, of group/version/kind %s/%s/%s",
			req.Name, req.Namespace, req.Kind.Group, req.Kind.Version, req.Kind.Kind)
	}

	UnsupportedOperationError = func(op admission_v1.Operation) error {
		return eris.Errorf("unsupported operation type %v found", op)
	}
)

func NewMultiClusterAdmissionValidator(
	clientset multicluster_v1alpha1.Clientset,
	matcher internal_placement.Matcher,
	parser rbac.Parser,
) MultiClusterAdmissionValidator {
	return &multiClusterAdmissionValidator{
		roleBindingClient: clientset.MultiClusterRoleBindings(),
		roleClient:        clientset.MultiClusterRoles(),
		matcher:           matcher,
		parser:            parser,
	}
}

type multiClusterAdmissionValidator struct {
	roleBindingClient multicluster_v1alpha1.MultiClusterRoleBindingClient
	roleClient        multicluster_v1alpha1.MultiClusterRoleClient
	matcher           internal_placement.Matcher
	parser            rbac.Parser
}

/*
	An action is mapped to a list of placements.

	The action is allowed if and only if for each of its placements, there exists a rule whose placement
	namespaces and clusters are either a wildcard or a superset of the resource's placement.

	Note that all applicable placement rules must be defined within a single MultiClusterBinding, i.e.
	a resource placement cannot be allowed if its placements are permitted by Rules that exist across
	multiple MultiClusterRoleBindings.
*/
func (m *multiClusterAdmissionValidator) ActionIsAllowed(
	ctx context.Context,
	roleBinding *multicluster_v1alpha1.MultiClusterRoleBinding,
	req *admission.Request,
) (allowed bool, err error) {
	logger := contextutils.LoggerFrom(ctx)
	role, err := m.roleClient.GetMultiClusterRole(ctx, client.ObjectKey{
		Name:      roleBinding.Spec.GetRoleRef().GetName(),
		Namespace: roleBinding.Spec.GetRoleRef().GetNamespace(),
	})
	if err != nil {
		return false, err
	}

	var rawObj []byte
	// attempt to unmarshal request object early, so it doesn't need to be done for each rule
	switch req.Operation {
	case admission_v1.Create:
		logger.Debugf("create event: %+v", req)
		rawObj = req.Object.Raw
	case admission_v1.Update:
		logger.Debugf("update event: %+v", req)
		rawObj = req.Object.Raw
	case admission_v1.Delete:
		// Passing old object based on this comment:
		// https://github.com/kubernetes-sigs/controller-runtime/blob/525a2d7dce5661cb90395b40a31575589ae504a2/pkg/webhook/admission/validator.go#L94
		logger.Debugf("delete event: %+v", req)
		rawObj = req.OldObject.Raw
	case admission_v1.Connect:
		return false, UnsupportedOperationError(req.Operation)
	}

	// Grab placement data from request object
	placements, err := m.parser.Parse(ctx, rawObj)
	if err != nil {
		return false, PlacementParsingError(err, req)
	}
	logger.Debugf("Parsed placements: %+v", placements)

	gvk := req.Kind
	// enumerate over rules to look for a match
	for _, ruleIter := range role.Spec.GetRules() {
		rule := ruleIter
		ruleCtx := contextutils.WithLoggerValues(ctx,
			zap.String("rule_type", rule.GetAction().String()),
			zap.String("rule_group", rule.GetApiGroup()),
			zap.String("rule_kind", rule.GetKind().GetValue()),
		)
		logger := contextutils.LoggerFrom(ruleCtx)
		// if groups do not match, continue immediately
		if rule.GetApiGroup() != gvk.Group {
			logger.Debug("api groups do not match")
			continue
		}
		// if kind is unspecified, apply to all kinds in the group
		if rule.GetKind() != nil && rule.GetKind().GetValue() != gvk.Kind {
			logger.Debug("kind is non-nil, and does not match")
			continue
		}

		switch rule.GetAction() {
		case multicluster_types.MultiClusterRoleSpec_Rule_CREATE:
			if req.Operation != admission_v1.Create {
				logger.Debugf("Skipping, action type is CREATE and req Operation is %s", req.Operation)
				continue
			}
		case multicluster_types.MultiClusterRoleSpec_Rule_UPDATE:
			if req.Operation != admission_v1.Update {
				logger.Debugf("Skipping, action type is UPDATE and req Operation is %s", req.Operation)
				continue
			}
		case multicluster_types.MultiClusterRoleSpec_Rule_DELETE:
			if req.Operation != admission_v1.Delete {
				logger.Debugf("Skipping, action type is DELETE and req Operation is %s", req.Operation)
				continue
			}
		}

		// Action is permitted only if all resource placements are allowed by some matching placement within the rule.
		if allowed := m.placementsAllowed(ruleCtx, placements, rule); allowed {
			return true, nil
		}
	}
	return false, nil
}

func (m *multiClusterAdmissionValidator) placementsAllowed(
	ruleCtx context.Context,
	placements []*multicluster_types.Placement,
	rule *multicluster_types.MultiClusterRoleSpec_Rule,
) bool {
	// Treat map as a set, ignore value.
	allowedPlacementsSet := map[*multicluster_types.Placement]interface{}{}
	for _, placement := range placements {
		if _, ok := allowedPlacementsSet[placement]; ok {
			continue
		}
		for _, rulePlacement := range rule.GetPlacements() {
			if m.matcher.Matches(ruleCtx, placement, rulePlacement) {
				allowedPlacementsSet[placement] = nil
				break
			}
		}
	}
	return len(allowedPlacementsSet) == len(placements)
}

func (m *multiClusterAdmissionValidator) GetMatchingMultiClusterRoleBindings(
	ctx context.Context,
	userInfo authv1.UserInfo,
) ([]*multicluster_v1alpha1.MultiClusterRoleBinding, error) {
	roleBindingList, err := m.roleBindingClient.ListMultiClusterRoleBinding(ctx)
	if err != nil {
		return nil, err
	}
	ctx = contextutils.WithLoggerValues(ctx,
		zap.String("user", userInfo.Username),
		zap.String("uid", userInfo.UID),
		zap.Strings("groups", userInfo.Groups),
	)
	var matching []*multicluster_v1alpha1.MultiClusterRoleBinding
	for _, roleBindingIter := range roleBindingList.Items {
		roleBinding := roleBindingIter
		ctx := contextutils.WithLoggerValues(ctx, zap.String("role_binding", roleBinding.GetName()))
		for _, subject := range roleBinding.Spec.GetSubjects() {
			logger := contextutils.LoggerFrom(contextutils.WithLoggerValues(ctx,
				zap.String("subject_kind", subject.GetKind().GetValue()),
				zap.String("subject_name", subject.GetName()),
				zap.String("subject_api_group", subject.GetApiGroup().GetValue()),
			))
			switch subject.GetKind().GetValue() {
			case "User":
				logger.Debug("found user binding")
				if userInfo.Username == subject.GetName() {
					matching = append(matching, &roleBinding)
				}
			case "Group":
				logger.Debug("found group binding")
				if stringutils.ContainsString(subject.GetName(), userInfo.Groups) {
					matching = append(matching, &roleBinding)
				}
			default:
				logger.Debugw("unsupported Kind for role binding subject", zap.Any("subject", subject))
				continue
			}
		}
	}
	return matching, nil
}
