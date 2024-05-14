package test

import (
	"fmt"

	. "github.com/onsi/gomega"
	"github.com/solo-io/k8s-utils/installutils/kuberesource"
	. "github.com/solo-io/k8s-utils/manifesttestutils"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//nolint:unparam // svcNamespace always receives "gloo-system"
func getService(testManifest TestManifest, svcNamespace string, svcName string) *corev1.Service {
	svcUns := testManifest.ExpectCustomResource("Service", svcNamespace, svcName)
	svcObj, err := kuberesource.ConvertUnstructured(svcUns)
	Expect(err).NotTo(HaveOccurred())
	Expect(svcObj).To(BeAssignableToTypeOf(&corev1.Service{}))
	return svcObj.(*corev1.Service)
}

//nolint:unparam // deploymentNamespace always receives "gloo-system"
func getDeployment(testManifest TestManifest, deploymentNamespace string, deploymentName string) *appsv1.Deployment {
	deployUns := testManifest.ExpectCustomResource("Deployment", deploymentNamespace, deploymentName)
	deployObj, err := kuberesource.ConvertUnstructured(deployUns)
	Expect(err).NotTo(HaveOccurred())
	Expect(deployObj).To(BeAssignableToTypeOf(&appsv1.Deployment{}))
	return deployObj.(*appsv1.Deployment)
}

//nolint:unparam // jobNamespace always receives "gloo-system"
func getJob(testManifest TestManifest, jobNamespace string, jobName string) *batchv1.Job {
	jobUns := testManifest.ExpectCustomResource("Job", jobNamespace, jobName)
	jobObj, err := kuberesource.ConvertUnstructured(jobUns)
	Expect(err).NotTo(HaveOccurred())
	Expect(jobObj).To(BeAssignableToTypeOf(&batchv1.Job{}))
	return jobObj.(*batchv1.Job)
}

//nolint:unparam // jobNamespace always receives "gloo-system"
func getConfigMap(testManifest TestManifest, namespace string, name string) *corev1.ConfigMap {
	configMapUns := testManifest.ExpectCustomResource("ConfigMap", namespace, name)
	configMapObj, err := kuberesource.ConvertUnstructured(configMapUns)
	Expect(err).NotTo(HaveOccurred())
	Expect(configMapObj).To(BeAssignableToTypeOf(&corev1.ConfigMap{}))
	return configMapObj.(*corev1.ConfigMap)
}

// verifies that the container contains an env var with the given name and value
func expectEnvVarExists(container corev1.Container, name string, value string) {
	foundName := false
	for _, envVar := range container.Env {
		if envVar.Name == name {
			Expect(envVar.Value).To(Equal(value), fmt.Sprintf("expected env var %s to have value %s", name, value))
			foundName = true
			break
		}
	}
	Expect(foundName).To(BeTrue(), fmt.Sprintf("env var with name %s should exist", name))
}

// verifies that the container does not contain an env var with the given name
func expectEnvVarDoesNotExist(container corev1.Container, name string) {
	for _, envVar := range container.Env {
		Expect(envVar.Name).NotTo(Equal(name), fmt.Sprintf("env var with name %s should not exist", name))
	}
}

func GetStructuredDeployment(t TestManifest, glooLabel string) *appsv1.Deployment {

	structuredDeployment := &appsv1.Deployment{}

	resources := t.SelectResources(func(u *unstructured.Unstructured) bool {
		if u.GetKind() == "Deployment" {
			if u.GetLabels()["gloo"] == glooLabel {
				deploymentObject, err := kuberesource.ConvertUnstructured(u)
				ExpectWithOffset(1, err).NotTo(HaveOccurred(), fmt.Sprintf("Deployment %s should be able to convert from unstructured", u.GetName()))
				var ok bool
				structuredDeployment, ok = deploymentObject.(*appsv1.Deployment)
				Expect(ok).To(BeTrue(), fmt.Sprintf("Deployment %+v should be able to cast to a structured deployment", u))
				return true
			}
		}
		return false
	})
	Expect(resources.NumResources()).To(Equal(1))

	return structuredDeployment
}
