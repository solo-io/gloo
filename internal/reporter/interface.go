package reporter

import (
	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/glue/pkg/api/types/v1"
)

type ObjectState string

const (
	ObjectStateAccepted ObjectState = "Accepted"
	ObjectStateRejected ObjectState = "Rejected"
)

type ConfigObjectReport struct {
	CfgObject v1.StorableConfigObject
	Err       *multierror.Error
}

type Interface interface {
	WriteReports(statuses []ConfigObjectReport) error
}
