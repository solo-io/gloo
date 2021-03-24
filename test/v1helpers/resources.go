package v1helpers

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type InputResourceGetter func() (resources.InputResource, error)

func EventuallyResourceAccepted(getter InputResourceGetter) {
	EventuallyResourceAcceptedWithOffset(1, getter)
}

func EventuallyResourceAcceptedWithOffset(ginkgoOffset int, getter InputResourceGetter) {
	gomega.EventuallyWithOffset(ginkgoOffset+1, func() (core.Status, error) {
		resource, err := getter()
		if err != nil || resource.GetStatus() == nil {
			return core.Status{}, errors.Wrapf(err, "waiting for %v to be accepted, but status is %v", resource.GetMetadata().GetName(), resource.GetStatus())
		}

		return *resource.GetStatus(), nil
	}, "15s", "0.5s").Should(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"Reason": gomega.BeEmpty(),
		"State":  gomega.Equal(core.Status_Accepted),
	}))
}
