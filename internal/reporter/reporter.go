package reporter

import (
	"github.com/pkg/errors"
	"github.com/solo-io/glue-storage"

	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
)

type reporter struct {
	store storage.Interface
}

func NewReporter(store storage.Interface) *reporter {
	return &reporter{store: store}
}

func (r *reporter) WriteReports(reports []ConfigObjectReport) error {
	for _, report := range reports {
		if err := r.writeReport(report); err != nil {
			return errors.Wrapf(err, "failed to write report for crd %v", report.CfgObject)
		}
		log.Debugf("wrote report for %v", report.CfgObject.GetName())
	}
	return nil
}

func (r *reporter) writeReport(report ConfigObjectReport) error {
	status := &v1.Status{
		State: v1.Status_Accepted,
	}
	if report.Err != nil {
		status.State = v1.Status_Rejected
		status.Reason = report.Err.Error()
	}
	name := report.CfgObject.GetName()
	switch report.CfgObject.(type) {
	case *v1.Upstream:
		us, err := r.store.V1().Upstreams().Get(report.CfgObject.GetName())
		if err != nil {
			return errors.Wrapf(err, "failed to find upstream %v", name)
		}
		us.Status = status
		if _, err := r.store.V1().Upstreams().Update(us); err != nil {
			return errors.Wrapf(err, "failed to update upstream store with status report")
		}
	case *v1.VirtualHost:
		virtualHost, err := r.store.V1().VirtualHosts().Get(name)
		if err != nil {
			return errors.Wrapf(err, "failed to find virtualhost %v", name)
		}
		virtualHost.Status = status
		if _, err := r.store.V1().VirtualHosts().Update(virtualHost); err != nil {
			return errors.Wrapf(err, "failed to update virtualhost store with status report")
		}
	}
	return nil
}
