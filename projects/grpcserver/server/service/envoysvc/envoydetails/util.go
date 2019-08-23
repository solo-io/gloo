package envoydetails

import kubev1 "k8s.io/api/core/v1"

func getName(pod kubev1.Pod) string {
	if id, ok := pod.Labels[GatewayProxyIdLabel]; ok {
		return id
	}
	return pod.Name
}
