package updater

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/internal/updater/gcf"
	"github.com/solo-io/gloo-function-discovery/internal/updater/lambda"
	"github.com/solo-io/gloo-function-discovery/internal/updater/swagger"
	"github.com/solo-io/gloo-function-discovery/pkg/functiontypes"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

var localCache = make(map[string][]*v1.Function)

// if forceSync is set, ignore the local cache and poll for new function list anyway
// we want to forceSync on every refreshDuration
// on a config / secrets change, we don't want to force sync
// else we can get into an update loop
func UpdateFunctionalUpstreams(gloo storage.Interface, upstreams []*v1.Upstream, secrets secretwatcher.SecretMap, forceSync bool) error {
	// nothing to do
	if len(upstreams) == 0 {
		return nil
	}
	for _, us := range upstreams {
		if !forceSync && functionListsEqual(us.Functions, localCache[us.Name]) {
			// ignore upstreams whose function list matches our cache
			continue
		}
		var funcs []*v1.Function
		var err error
		switch functiontypes.GetFunctionType(us) {
		case functiontypes.FunctionTypeLambda:
			if len(secrets) == 0 {
				log.Warnf("lambda upstream detected, but no secrets have been read yet")
				continue
			}
			funcs, err = lambda.GetFuncs(us, secrets)
			if err != nil {
				return errors.Wrap(err, "updating lambda functions")
			}
		case functiontypes.FunctionTypeGfuncs:
			if len(secrets) == 0 {
				log.Warnf("google functions upstream detected, but no secrets have been read yet")
				continue
			}
			funcs, err = gcf.GetFuncs(us, secrets)
			if err != nil {
				return errors.Wrap(err, "updating google functions")
			}
		case functiontypes.FunctionTypeSwagger:
			funcs, err = swagger.GetFuncs(us)
			if err != nil {
				return errors.Wrap(err, "updating swagger functions")
			}
		}
		if err := updateUpstreamWithFuncs(gloo, us, funcs); err != nil {
			return errors.Wrap(err, "updating upstream object with new funcs")
		}
	}
	return nil
}

func updateUpstreamWithFuncs(gloo storage.Interface, us *v1.Upstream, funcs []*v1.Function) error {
	// sort funcs for idempotency
	sort.SliceStable(funcs, func(i, j int) bool {
		return funcs[i].Name < funcs[j].Name
	})
	// no update to do
	if functionListsEqual(us.Functions, funcs) {
		return nil
	}
	us.Functions = funcs
	localCache[us.Name] = funcs
	_, err := gloo.V1().Upstreams().Update(us)
	return err
}

func functionListsEqual(funcs1, funcs2 []*v1.Function) bool {
	if len(funcs1) != len(funcs2) {
		return false
	}
	for i := range funcs1 {
		fn1 := funcs1[i]
		fn2 := funcs2[i]
		if !fn1.Equal(fn2) {
			return false
		}
	}
	return true
}
