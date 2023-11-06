package kubeutils

import (
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/testutils"
)

const (
	defaultEventuallyTimeout         = 30 * time.Second
	defaultEventuallyPollingInterval = 1 * time.Second
)

// WaitForRollout waits for the specified deployment to be rolled out successfully.
func WaitForRollout(deploymentName string, deploymentNamespace string, intervals ...interface{}) {
	WaitForRolloutWithOffset(1, deploymentName, deploymentNamespace, intervals...)
}

// RolloutRestart restarts the specified deployment and waits for it to be rolled out successfully.
func RolloutRestart(deploymentName string, deploymentNamespace string, intervals ...interface{}) {
	RolloutRestartWithOffset(1, deploymentName, deploymentNamespace, intervals...)
}

func WaitForRolloutWithOffset(offset int, deploymentName string, deploymentNamespace string, intervals ...interface{}) {
	timeoutInterval, pollingInterval := getTimeoutAndPollingIntervalsOrDefault(intervals...)
	EventuallyWithOffset(offset+1, func() (bool, error) {
		out, err := testutils.KubectlOut("rollout", "status", "-n", deploymentNamespace, fmt.Sprintf("deployment/%s", deploymentName))
		fmt.Println(out)
		return strings.Contains(out, "successfully rolled out"), err
	}, timeoutInterval, pollingInterval).Should(BeTrue())
}

func RolloutRestartWithOffset(offset int, deploymentName string, deploymentNamespace string, intervals ...interface{}) {
	out, err := testutils.KubectlOut("rollout", "restart", "-n", deploymentNamespace, fmt.Sprintf("deployment/%s", deploymentName))
	fmt.Println(out)
	ExpectWithOffset(offset+1, err).ToNot(HaveOccurred())

	WaitForRolloutWithOffset(offset+1, deploymentName, deploymentNamespace, intervals...)
}

// this is copied from gloo https://github.com/solo-io/gloo/blob/067ce54aa9de3cdcce89c6465b101cbc599d9e4f/test/helpers/input_resources.go#L82
// should ideally be reused but for now it's being duplicated so we don't need to depend on a gloo release to expose/use it
func getTimeoutAndPollingIntervalsOrDefault(intervals ...interface{}) (interface{}, interface{}) {
	var timeoutInterval, pollingInterval interface{}

	timeoutInterval = defaultEventuallyTimeout
	pollingInterval = defaultEventuallyPollingInterval

	if len(intervals) > 0 {
		timeoutInterval = intervals[0]
	}
	if len(intervals) > 1 {
		pollingInterval = intervals[1]
	}

	return timeoutInterval, pollingInterval
}
