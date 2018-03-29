package detector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-function-discovery/internal/detector"
	"github.com/solo-io/gloo-function-discovery/pkg/resolver"
	"github.com/solo-io/gloo-testing/helpers"
)

var _ = Describe("Marker", func() {
	Context("happy path", func() {
		resolve := &resolver.Resolver{}
		marker := NewMarker([]Detector{
			&mockDetector{shouldErr: true},
			&mockDetector{shouldErr: false},
		}, resolve)
		us := helpers.NewTestUpstream2()
		err := marker.MarkFunctionalUpstream(us)
		Expect(err).NotTo(HaveOccurred())
		Expect(us.ServiceInfo).NotTo(BeNil())
	})
})

type mockDetector struct {
	shouldErr bool
}

func (d *mockDetector) DetectFunctionalService(addr string) (*v1.ServiceInfo, map[string]string, error) {
	if d.shouldErr {
		return nil, nil, errors.New("mock: failed detection")
	}
	return &v1.ServiceInfo{Type: "mock_service"}, nil, nil
}
