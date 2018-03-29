package updater

import (
	"sort"

	"reflect"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/internal/detector"
	"github.com/solo-io/gloo-function-discovery/internal/updater/gcf"
	"github.com/solo-io/gloo-function-discovery/internal/updater/lambda"
	"github.com/solo-io/gloo-function-discovery/internal/updater/swagger"
	"github.com/solo-io/gloo-function-discovery/pkg/functiontypes"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

func GetSecretRefsToWatch(upstreams []*v1.Upstream) []string {
	var refs []string
	for _, us := range upstreams {
		switch functiontypes.GetFunctionType(us) {
		case functiontypes.FunctionTypeLambda:
			ref, err := lambda.GetSecretRef(us)
			if err != nil {
				continue
			}
			refs = append(refs, ref)
		case functiontypes.FunctionTypeGfuncs:
			ref, err := gcf.GetSecretRef(us)
			if err != nil {
				continue
			}
			refs = append(refs, ref)
		}
	}
	return refs
}

// if forceSync is set, ignore the local cache and poll for new function list anyway
// we want to forceSync on every refreshDuration
// on a config / secrets change, we don't want to force sync
// else we can get into an update loop
func UpdateFunctions(gloo storage.Interface, upstreamName string, secrets secretwatcher.SecretMap) error {
	us, err := gloo.V1().Upstreams().Get(upstreamName)
	if err != nil {
		return errors.Wrapf(err, "failed to get existing upstream with name %v", upstreamName)
	}

	var funcs []*v1.Function
	switch functiontypes.GetFunctionType(us) {
	case functiontypes.FunctionTypeLambda:
		if len(secrets) == 0 {
			log.Warnf("lambda upstream detected, but no secrets have been read yet")
			return nil
		}
		funcs, err = lambda.GetFuncs(us, secrets)
		if err != nil {
			return errors.Wrap(err, "updating lambda functions")
		}
	case functiontypes.FunctionTypeGfuncs:
		if len(secrets) == 0 {
			log.Warnf("google functions upstream detected, but no secrets have been read yet")
			return nil
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
	default:
		return errors.Errorf("unknown function type")
	}

	if err := updateUpstreamWithFuncs(gloo, us.Name, funcs); err != nil {
		return errors.Wrap(err, "updating upstream object with new funcs")
	}
	return nil
}

func updateUpstreamWithFuncs(gloo storage.Interface, upstreamName string, funcs []*v1.Function) error {
	// sort funcs for idempotency
	sort.SliceStable(funcs, func(i, j int) bool {
		return funcs[i].Name < funcs[j].Name
	})

	usToUpdate, err := gloo.V1().Upstreams().Get(upstreamName)
	if err != nil {
		return errors.Wrapf(err, "failed to get existing upstream with name %v", upstreamName)
	}

	// no update to do
	if functionListsEqual(usToUpdate.Functions, funcs) {
		return nil
	}

	usToUpdate.Functions = mergeFuncs(usToUpdate.Functions, funcs)

	_, err = gloo.V1().Upstreams().Update(usToUpdate)
	if err != nil {
		return err
	}
	return nil
}

// get the unique set of funcs between two lists
// if conflict, new wins
func mergeFuncs(oldFuncs, newFuncs []*v1.Function) []*v1.Function {
	var notReplaced []*v1.Function
	for _, oldFunc := range oldFuncs {
		var replace bool
		for _, newFunc := range newFuncs {
			if newFunc.Name == oldFunc.Name {
				replace = true
				break
			}
		}
		if replace {
			continue
		}
		notReplaced = append(notReplaced, oldFunc)
	}
	return append(notReplaced, newFuncs...)
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

// update the upstream with service info and annotations
func UpdateServiceInfo(gloo storage.Interface,
	upstreamName string,
	marker *detector.Marker) error {
	usToUpdate, err := gloo.V1().Upstreams().Get(upstreamName)
	if err != nil {
		return errors.Wrapf(err, "failed to get existing upstream with name %v", upstreamName)
	}

	svcInfo, annotations, err := marker.DetectFunctionalUpstream(usToUpdate)
	if err != nil {
		return errors.Wrapf(err, "failed to discover whether %v is a functional upstream", usToUpdate.Name)
	}

	// no update to do
	if annotationsEqual(usToUpdate, annotations) && svcInfoEqual(usToUpdate, svcInfo) {
		return nil
	}

	usToUpdate.Metadata.Annotations = mergeAnnotations(usToUpdate.Metadata.Annotations, annotations)
	usToUpdate.ServiceInfo = svcInfo

	_, err = gloo.V1().Upstreams().Update(usToUpdate)
	if err != nil {
		return err
	}
	return nil
}

// get the unique set of funcs between two lists
// if conflict, new wins
func mergeAnnotations(oldAnnotations, newAnnotations map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range oldAnnotations {
		merged[k] = v
	}
	for k, v := range newAnnotations {
		merged[k] = v
	}
	return merged
}

func annotationsEqual(us *v1.Upstream, annotations map[string]string) bool {
	if us.Metadata == nil {
		if len(annotations) == 0 {
			return true
		}
		return false
	}
	if len(us.Metadata.Annotations) != len(annotations) {
		return false
	}
	for k, v := range us.Metadata.Annotations {
		if annotations[k] != v {
			return false
		}
	}
	return true
}

func svcInfoEqual(us *v1.Upstream, svcInfo *v1.ServiceInfo) bool {
	return reflect.DeepEqual(us.ServiceInfo, svcInfo)
}
