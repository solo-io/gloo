package checker

import (
	"context"
	"fmt"

	v1sets "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/sets"
	"github.com/solo-io/go-utils/contextutils"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery/translator/summarize"
	corev1 "k8s.io/api/core/v1"
)

func GetPodsSummary(ctx context.Context, set v1sets.PodSet, namespace, cluster string) *types.GlooInstanceSpec_Check_Summary {
	summary := &types.GlooInstanceSpec_Check_Summary{}
	for _, podIter := range set.List() {
		pod := podIter
		if pod.ClusterName != cluster || pod.Namespace != namespace {
			continue
		}

		summary.Total += 1
		for _, condition := range pod.Status.Conditions {
			var message string

			if condition.Message != "" {
				message = fmt.Sprintf(" Message: %s", condition.Message)
			}

			// if condition is not met and the pod is not completed
			conditionNotMet := condition.Status != corev1.ConditionTrue && condition.Reason != "PodCompleted"

			// possible condition types listed at https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-conditions
			switch condition.Type {
			case corev1.PodScheduled:
				if conditionNotMet {
					summary.Errors = append(
						summary.Errors,
						&types.GlooInstanceSpec_Check_Summary_ResourceReport{
							Ref: &v1.ObjectRef{
								Namespace: pod.Namespace,
								Name:      pod.Name,
							},
							Message: fmt.Sprintf("Pod is not yet scheduled!%s\n", message),
						},
					)
				}
			case corev1.PodReady:
				if conditionNotMet {
					summary.Errors = append(
						summary.Errors,
						&types.GlooInstanceSpec_Check_Summary_ResourceReport{
							Ref: &v1.ObjectRef{
								Namespace: pod.Namespace,
								Name:      pod.Name,
							},
							Message: fmt.Sprintf("Pod is not ready!%s\n", message),
						},
					)
				}
			case corev1.PodInitialized:
				if conditionNotMet {
					summary.Errors = append(
						summary.Errors,
						&types.GlooInstanceSpec_Check_Summary_ResourceReport{
							Ref: &v1.ObjectRef{
								Namespace: pod.Namespace,
								Name:      pod.Name,
							},
							Message: fmt.Sprintf("Pod is not yet initialized!%s\n", message),
						},
					)
				}
			case corev1.PodReasonUnschedulable:
				if conditionNotMet {
					summary.Errors = append(
						summary.Errors,
						&types.GlooInstanceSpec_Check_Summary_ResourceReport{
							Ref: &v1.ObjectRef{
								Namespace: pod.Namespace,
								Name:      pod.Name,
							},
							Message: fmt.Sprintf("Pod is unschedulable!%s\n", message),
						},
					)
				}
			case corev1.ContainersReady:
				if conditionNotMet {
					summary.Errors = append(
						summary.Errors,
						&types.GlooInstanceSpec_Check_Summary_ResourceReport{
							Ref: &v1.ObjectRef{
								Namespace: pod.Namespace,
								Name:      pod.Name,
							},
							Message: fmt.Sprintf("Not all containers are ready!%s\n", message),
						},
					)
				}
			default:
				contextutils.LoggerFrom(ctx).Debugw("Note: Unhandled pod condition %s", condition.Type)
			}
		}
	}

	summarize.SortLists(summary)
	return summary
}
