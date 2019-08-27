package helpers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/status"
)

type stubResource struct {
	statusCode core.Status_State
	reason     string
}

func (s stubResource) GetMetadata() core.Metadata {
	return core.Metadata{
		Name:      "n",
		Namespace: "ns",
	}
}

func (s stubResource) SetMetadata(meta core.Metadata) {
	panic("implement me")
}

func (s stubResource) Equal(that interface{}) bool {
	panic("implement me")
}

func (s stubResource) GetStatus() core.Status {
	return core.Status{
		State:  s.statusCode,
		Reason: s.reason,
	}
}

func (s stubResource) SetStatus(status core.Status) {
	panic("implement me")
}

var _ resources.InputResource = stubResource{}

var _ = Describe("InputResourceStatusGetter Test", func() {
	var (
		statusConverter = status.NewInputResourceStatusGetter()
	)

	Describe("GetApiStatusFromResource", func() {
		It("translates the status correctly when the resource is healthy", func() {
			healthyResource := &stubResource{
				statusCode: core.Status_Accepted,
			}

			convertedStatus := statusConverter.GetApiStatusFromResource(healthyResource)

			Expect(convertedStatus.Code).To(Equal(v1.Status_OK))
			Expect(convertedStatus.Message).To(Equal(""))
		})

		It("translates the status correctly when the resource is unhealthy", func() {
			unhealthyResource := &stubResource{
				statusCode: core.Status_Rejected,
				reason:     "test-reason",
			}

			convertedStatus := statusConverter.GetApiStatusFromResource(unhealthyResource)

			Expect(convertedStatus.Code).To(Equal(v1.Status_ERROR))
			Expect(convertedStatus.Message).To(Equal(status.ResourceRejected("ns", "n", "test-reason")))
		})

		It("gives a helpful message with a missed switch case", func() {
			unknownStatusResource := &stubResource{
				statusCode: -1337,
			}

			convertedStatus := statusConverter.GetApiStatusFromResource(unknownStatusResource)

			Expect(convertedStatus.Code).To(Equal(v1.Status_ERROR))
			Expect(convertedStatus.Message).To(Equal(status.UnknownFailure("ns", "n", -1337)))
		})

		It("translates the status correctly when the resource is pending", func() {
			pendingResource := &stubResource{
				statusCode: core.Status_Pending,
			}

			convertedStatus := statusConverter.GetApiStatusFromResource(pendingResource)

			Expect(convertedStatus.Code).To(Equal(v1.Status_WARNING))
			Expect(convertedStatus.Message).To(Equal(status.ResourcePending("ns", "n")))
		})
	})
})
