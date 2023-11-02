package argocd_test

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-projects/test/kube2e"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var testHelper *helper.SoloTestHelper

const namespace = defaults.GlooSystem

var _ = Describe("Kube2e: ArgoCD", func() {

	BeforeEach(func() {
		var err error
		ctx := context.TODO()

		testHelper, err = kube2e.GetEnterpriseTestHelper(ctx, namespace)
		Expect(err).NotTo(HaveOccurred())
	})

	It("Tests the Lifecycle", func() {
		installGloo()
		checkGlooHealthy(testHelper)

		// Sync once more to simulate an upgrade
		syncGloo()
		checkGlooHealthy(testHelper)

		uninstallGloo()
	})
})

func installGloo() {
	uninstallGloo()

	repo := "http://helm-repo/"

	// argocd --core app create gloo \
	// --repo http://helm-repo/ --helm-chart gloo --revision $VERSION \
	// --dest-namespace gloo-system --dest-server https://kubernetes.default.svc \
	// --sync-option CreateNamespace=true --upsert --values-literal-file helm-override.yaml
	command := []string{"--core", "app", "create", "gloo-ee",
		"--repo", repo, "--helm-chart", "gloo-ee", "--revision", version, "--helm-set", "license_key=" + testHelper.LicenseKey,
		"--dest-namespace", "gloo-system", "--dest-server", "https://kubernetes.default.svc",
		"--sync-option", "CreateNamespace=true", "--upsert", "--values-literal-file", "helm-override.yaml"}
	fmt.Printf("Running argo command : %s\n", command)
	runAndCleanCommand("argocd", command...)

	syncGloo()
}

func syncGloo() {
	// argocd --core app sync gloo
	command := []string{"--core", "app", "sync", "gloo-ee"}
	fmt.Printf("Running argo command : %s\n", command)
	runAndCleanCommand("argocd", command...)
}

func uninstallGloo() {
	// argocd --core app delete gloo -y
	cmd := exec.Command("argocd", "--core", "app", "delete", "gloo-ee", "-y")
	cmd.Output()
}

func checkRolloutJobDeleted() {
	// Wait `gateway.rolloutJob.ttlSecondsAfterFinished` until the resource rollout job has been deleted
	fmt.Println("Waiting for the gloo-resource-rollout job to be cleaned up")
	time.Sleep(60 * time.Second)
	EventuallyWithOffset(1, func() string {
		cmd := exec.Command("kubectl", "-n", "gloo-system", "get", "jobs", "-A")
		b, err := cmd.Output()
		Expect(err).To(BeNil())
		return string(b)
	}, "60s", "10s").ShouldNot(
		And(ContainSubstring("gloo-resource-rollout "),
			ContainSubstring("gloo-ee-resource-rollout ")))
}

func checkGlooHealthyAndSyncedInArgo() {
	// Get the state of gloo
	// argocd app get gloo --hard-refresh -o json | jq '.status.health.status'
	fmt.Println("Checking if gloo is healthy")
	EventuallyWithOffset(1, func() string {
		command := "argocd app get gloo-ee --hard-refresh -o json | jq '.status.health.status'"
		cmd := exec.Command("bash", "-c", command)
		b, err := cmd.Output()
		Expect(err).To(BeNil())
		return string(b)
	}).Should(
		ContainSubstring("Healthy"))
	// argocd app get gloo --hard-refresh -o json | jq '.status.sync.status'
	fmt.Println("Checking if gloo is synced")
	EventuallyWithOffset(1, func() string {
		command := "argocd app get gloo-ee --hard-refresh -o json | jq '.status.sync.status'"
		cmd := exec.Command("bash", "-c", command)
		b, err := cmd.Output()
		Expect(err).To(BeNil())
		return string(b)
	}).Should(
		ContainSubstring("Synced"))
}

func runAndCleanCommand(name string, arg ...string) []byte {
	ctx, _ := context.WithTimeout(context.TODO(), 300*time.Second)
	cmd := exec.CommandContext(ctx, name, arg...)
	b, err := cmd.Output()
	// for debugging in Cloud Build
	if err != nil {
		fmt.Println(string(b))
		fmt.Println(err.Error())
		if v, ok := err.(*exec.ExitError); ok {
			fmt.Println("ExitError: ", string(v.Stderr))
		}
	}
	Expect(err).To(BeNil())
	cmd.Process.Kill()
	cmd.Process.Release()
	return b
}

func checkGlooHealthy(testHelper *helper.SoloTestHelper) {
	// Wait until the resource rollout job has been deleted to ensure that argo does not go out of sync
	checkRolloutJobDeleted()
	kube2e.CheckGlooHealthy(testHelper)
}
