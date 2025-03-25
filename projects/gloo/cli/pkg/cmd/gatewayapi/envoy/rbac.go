package envoy

import (
	"errors"
	"strconv"

	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
)

// reverseTranslateRbac converts an Envoy RBAC configuration back to Gloo RBAC policies
func reverseTranslateRbac(rbacPerRoute *envoyauthz.RBACPerRoute) (map[string]*rbac.Policy, error) {
	if rbacPerRoute == nil || rbacPerRoute.Rbac == nil || rbacPerRoute.Rbac.Rules == nil {
		return nil, nil
	}

	userPolicies := make(map[string]*rbac.Policy)

	// Only handle ALLOW action as that's what we generate in translateRbac
	if rbacPerRoute.Rbac.Rules.Action != envoycfgauthz.RBAC_ALLOW {
		return nil, errors.New("only ALLOW action is supported")
	}

	// Convert each Envoy policy back to a Gloo policy
	for name, envoyPolicy := range rbacPerRoute.Rbac.Rules.Policies {
		policy := &rbac.Policy{}

		// Handle permissions
		if len(envoyPolicy.Permissions) > 0 {
			permission := &rbac.Permissions{}

			// We only generate either a single permission or an AND of permissions
			mainPerm := envoyPolicy.Permissions[0]

			// Handle AND rules case
			if andRules := mainPerm.GetAndRules(); andRules != nil {
				for _, rule := range andRules.Rules {
					if err := extractPermissionRule(rule, permission); err != nil {
						return nil, err
					}
				}
			} else {
				// Handle single rule case
				if err := extractPermissionRule(mainPerm, permission); err != nil {
					return nil, err
				}
			}

			if !isEmpty(permission) {
				policy.Permissions = permission
			}
		}

		// Handle principals
		for _, envoyPrincipal := range envoyPolicy.Principals {
			principal, err := reverseTranslateJwtPrincipal(envoyPrincipal)
			if err != nil {
				return nil, err
			}
			if principal != nil {
				policy.Principals = append(policy.Principals, principal)
			}
		}

		userPolicies[name] = policy
	}

	return userPolicies, nil
}

func extractPermissionRule(rule *envoycfgauthz.Permission, permission *rbac.Permissions) error {
	if header := rule.GetHeader(); header != nil {
		switch header.Name {
		case ":path":
			if prefixMatch := header.GetPrefixMatch(); prefixMatch != "" {
				permission.PathPrefix = prefixMatch
			}
		case ":method":
			if exactMatch := header.GetExactMatch(); exactMatch != "" {
				permission.Methods = append(permission.Methods, exactMatch)
			}
		}
	} else if orRules := rule.GetOrRules(); orRules != nil {
		// Handle OR rules for methods
		for _, orRule := range orRules.Rules {
			if header := orRule.GetHeader(); header != nil && header.Name == ":method" {
				if exactMatch := header.GetExactMatch(); exactMatch != "" {
					permission.Methods = append(permission.Methods, exactMatch)
				}
			}
		}
	}
	return nil
}

func reverseTranslateJwtPrincipal(principal *envoycfgauthz.Principal) (*rbac.Principal, error) {
	if principal == nil {
		return nil, nil
	}

	// Handle AND rules case
	if andIds := principal.GetAndIds(); andIds != nil {
		// We expect all AND rules to be metadata matchers for the same JWT principal
		jwtPrincipal := &rbac.JWTPrincipal{
			Claims: make(map[string]string),
		}

		for _, id := range andIds.Ids {
			if metadata := id.GetMetadata(); metadata != nil {
				if err := extractJwtClaim(metadata, jwtPrincipal); err != nil {
					return nil, err
				}
			}
		}

		return &rbac.Principal{
			JwtPrincipal: jwtPrincipal,
		}, nil
	}

	// Handle single metadata matcher case
	if metadata := principal.GetMetadata(); metadata != nil {
		jwtPrincipal := &rbac.JWTPrincipal{
			Claims: make(map[string]string),
		}

		if err := extractJwtClaim(metadata, jwtPrincipal); err != nil {
			return nil, err
		}

		return &rbac.Principal{
			JwtPrincipal: jwtPrincipal,
		}, nil
	}

	return nil, nil
}

func extractJwtClaim(metadata *envoymatcher.MetadataMatcher, jwtPrincipal *rbac.JWTPrincipal) error {
	if len(metadata.Path) < 2 {
		return errors.New("invalid metadata path length")
	}

	// Extract claim name from the path
	claimName := metadata.Path[len(metadata.Path)-1].GetKey()

	// Extract value and matcher type
	if stringMatch := metadata.Value.GetStringMatch(); stringMatch != nil {
		if exact := stringMatch.GetExact(); exact != "" {
			jwtPrincipal.Claims[claimName] = exact
			jwtPrincipal.Matcher = rbac.JWTPrincipal_EXACT_STRING
		}
	} else if boolMatch := metadata.Value.GetBoolMatch(); boolMatch {
		jwtPrincipal.Claims[claimName] = strconv.FormatBool(boolMatch)
		jwtPrincipal.Matcher = rbac.JWTPrincipal_BOOLEAN
	} else if listMatch := metadata.Value.GetListMatch(); listMatch != nil {
		if oneOf := listMatch.GetOneOf(); oneOf != nil {
			if stringMatch := oneOf.GetStringMatch(); stringMatch != nil {
				if exact := stringMatch.GetExact(); exact != "" {
					jwtPrincipal.Claims[claimName] = exact
					jwtPrincipal.Matcher = rbac.JWTPrincipal_LIST_CONTAINS
				}
			}
		}
	}

	return nil
}
