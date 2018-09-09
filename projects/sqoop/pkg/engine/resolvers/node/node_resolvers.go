package node

import (
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/types/v1"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/exec"
)

func NewNodeResolver(resolver *v1.NodeJSResolver) (exec.RawResolver, error) {
	return nil, errors.Errorf("nodejs resolvers currently unsupported")
}
