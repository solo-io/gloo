package setup

import (
	"context"
	"os"

	"github.com/solo-io/go-utils/contextutils"
)

const NamespacedRbacEnvName = "RBAC_NAMESPACED"

type RbacNamespaced bool

func MustGetRbacNamespaced(ctx context.Context) RbacNamespaced {
	rbacNamespaced := os.Getenv(NamespacedRbacEnvName)
	if rbacNamespaced == "true" {
		contextutils.LoggerFrom(ctx).Infof("%s environment variable is set to `true`. Running server in namespaced mode", NamespacedRbacEnvName)
		return true
	}
	return false
}
