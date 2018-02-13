package crd

import (
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/glue/internal/pkg/kube/storage"
	"github.com/solo-io/glue/internal/reporter"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
	clientset "github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
)

type kubeReporter struct {
	client clientset.Interface
}

func NewKubeReporter(client clientset.Interface) *kubeReporter {
	return &kubeReporter{client: client}
}

func (r *kubeReporter) WriteReports(reports []reporter.ConfigObjectReport) error {
	for _, report := range reports {
		if err := r.writeReport(report); err != nil {
			return errors.Wrapf(err, "failed to write report for crd %v", report.CfgObject)
		}
		log.Debugf("wrote report for %v", report.CfgObject.GetStorageRef())
	}
	return nil
}

func (r *kubeReporter) writeReport(report reporter.ConfigObjectReport) error {
	ref := report.CfgObject.GetStorageRef()
	namespace, name, err := storage.ParseStorageRef(ref)
	if err != nil {
		return errors.Wrapf(err, "failed to parse kubernetes storage ref: %v", ref)
	}
	status := crdv1.CrdObjectStatus{
		State: reporter.ObjectStateAccepted,
	}
	if report.Err != nil {
		status.State = reporter.ObjectStateRejected
		status.Reason = report.Err.Error()
	}
	switch report.CfgObject.(type) {
	case *v1.Upstream:
		us, err := r.client.GlueV1().Upstreams(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to find crd for upstream %v", ref)
		}
		us.Status = status
		if _, err := r.client.GlueV1().Upstreams(namespace).Update(us); err != nil {
			return errors.Wrapf(err, "failed to update upstream crd with status report")
		}
	case *v1.VirtualHost:
		virtualHost, err := r.client.GlueV1().VirtualHosts(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to find crd for virtualhost %v", ref)
		}
		virtualHost.Status = status
		if _, err := r.client.GlueV1().VirtualHosts(namespace).Update(virtualHost); err != nil {
			return errors.Wrapf(err, "failed to update virtualhost crd with status report")
		}
	}
	return nil
}
