package upstreamgroupsvc

import (
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	FailedToReadUpstreamGroupError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to read upstream group %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToListUpstreamGroupsError = func(err error, namespace string) error {
		return errors.Wrapf(err, "Failed to list upstream groups in %v", namespace)
	}

	FailedToCreateUpstreamGroupError = func(err error, namespace, name string) error {
		return errors.Wrapf(err, "Failed to create upstream group %v.%v", namespace, name)
	}

	FailedToUpdateUpstreamGroupError = func(err error, namespace, name string) error {
		return errors.Wrapf(err, "Failed to update upstream group %v.%v", namespace, name)
	}

	FailedToParseUpstreamGroupFromYamlError = func(err error, namespace, name string) error {
		return errors.Wrapf(err, "Failed to parse upstream group %s.%s from YAML", namespace, name)
	}

	FailedToDeleteUpstreamGroupError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to delete upstream group %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToCheckIsUpstreamGroupReferencedError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to verify that upstream group %v.%v is not referenced in a virtual service", ref.GetNamespace(), ref.GetName())
	}

	CannotDeleteReferencedUpstreamGroupError = func(upstreamGroupRef *core.ResourceRef, virtualServiceRefs []*core.ResourceRef) error {
		return errors.Errorf("UpstreamGroup %v.%v is referenced in virtual service(s) %v",
			upstreamGroupRef.GetNamespace(),
			upstreamGroupRef.GetName(),
			virtualServiceRefs,
		)
	}
)
