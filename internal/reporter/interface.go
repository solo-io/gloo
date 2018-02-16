package reporter

import (
	"github.com/solo-io/glue/pkg/api/types/v1"
)

type ConfigObjectReport struct {
	CfgObject v1.ConfigObject
	Err       error
}

type Interface interface {
	WriteReports(statuses []ConfigObjectReport) error
}
