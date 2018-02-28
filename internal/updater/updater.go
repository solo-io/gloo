package updater

import (
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/pkg/functiontypes"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

func UpdateFunctionalUpstreams(gloo storage.Interface, upstreams []*v1.Upstream, secrets secretwatcher.SecretMap) error {
	// nothing to do
	if len(upstreams) == 0 {
		return nil
	}
	for _, us := range upstreams {
		switch functiontypes.GetFunctionType(us) {
		case functiontypes.FunctionTypeLambda:
			if len(secrets) == 0 {
				log.Warnf("lambda upstream detected, but no secrets have been read yet")
				continue
			}
			lambdaFuncs, err := getLambdaFuncs(us, secrets)
			if err != nil {
				return errors.Wrap(err, "updating lambda functions")
			}
			// todo: write these bad boys to the upstream
		}
	}
	return nil
}

func updateUpstreamWithFuncs(us *v1.Upstream, funcs []*v1.Function) {
	// sort funcs, make sure
}
