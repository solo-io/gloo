package checker

import (
	"context"
	"fmt"

	v1sets "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1/sets"
	"github.com/solo-io/go-utils/contextutils"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// Get a summary of deployments in the given namespace and cluster. To bypass the cluster check (e.g. for single-cluster
// use), pass in "" for the cluster.
func GetDeploymentsSummary(ctx context.Context, deployments v1sets.DeploymentSet, namespace, cluster string) *Summary {
	summary := &Summary{}

	for _, deploymentIter := range deployments.List() {
		deployment := deploymentIter

		if (cluster != "" && deployment.ClusterName != cluster) || deployment.Namespace != namespace {
			continue
		}

		summary.Total += 1

		errorMessage := getErrorMessage(ctx, deployment)
		if errorMessage != "" {
			summary.Errors = append(summary.Errors, getDeploymentError(deployment, errorMessage))
		}
	}

	SortLists(summary)
	return summary
}

func getErrorMessage(ctx context.Context, deployment *appsv1.Deployment) string {
	var output string

	// possible condition types listed at https://godoc.org/k8s.io/api/apps/v1#DeploymentConditionType
	// check for each condition independently because multiple conditions will be True and DeploymentReplicaFailure
	// tends to provide the most explicit error message.
	for _, condition := range deployment.Status.Conditions {
		message := getMessage(condition)
		if condition.Type == appsv1.DeploymentReplicaFailure && condition.Status == corev1.ConditionTrue {
			output = fmt.Sprintf("Deployment %s in namespace %s failed to create pods! %s", deployment.Name, deployment.Namespace, message)
		}
		if output != "" {
			return output
		}
	}

	for _, condition := range deployment.Status.Conditions {
		message := getMessage(condition)
		if condition.Type == appsv1.DeploymentProgressing && condition.Status != corev1.ConditionTrue {
			output = fmt.Sprintf("Deployment %s in namespace %s is not progressing! %s", deployment.Name, deployment.Namespace, message)
		}

		if output != "" {
			return output
		}
	}

	for _, condition := range deployment.Status.Conditions {
		message := getMessage(condition)
		if condition.Type == appsv1.DeploymentAvailable && condition.Status != corev1.ConditionTrue {
			output = fmt.Sprintf("Deployment %s in namespace %s is not available! %s", deployment.Name, deployment.Namespace, message)
		}

		if output != "" {
			return output
		}
	}

	for _, condition := range deployment.Status.Conditions {
		if condition.Type != appsv1.DeploymentAvailable &&
			condition.Type != appsv1.DeploymentReplicaFailure &&
			condition.Type != appsv1.DeploymentProgressing {
			contextutils.LoggerFrom(ctx).Debugw("Note: Unhandled deployment condition %s", condition.Type)
		}
	}
	return ""
}

func getDeploymentError(deployment *appsv1.Deployment, message string) *ResourceReport {
	return &ResourceReport{
		Ref: &v1.ObjectRef{
			Namespace: deployment.Namespace,
			Name:      deployment.Name,
		},
		Message: message,
	}
}

func getMessage(c appsv1.DeploymentCondition) string {
	if c.Message != "" {
		return fmt.Sprintf("Message: %s", c.Message)
	}
	return ""
}
