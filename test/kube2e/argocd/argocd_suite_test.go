package argocd_test

import (
	"fmt"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/gloo/test/helpers"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-projects/test/kube2e"
)

func TestArgoCD(t *testing.T) {
	helpers.RegisterGlooDebugLogPrintHandlerAndClearLogs()
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "ArgoCD Suite")
}

var _ = BeforeSuite(func() {
	installArgoCD()
	deployHelmRepo()
})

var _ = AfterSuite(func() {
	cleanup()
})

func installArgoCD() {
	fmt.Println("Installing ArgoCD")
	// helm repo add argo https://argoproj.github.io/argo-helm
	kube2e.RunAndCleanCommand("helm", "repo", "add", "argo", "https://argoproj.github.io/argo-helm")

	// helm install argocd argo/argo-cd --wait
	kube2e.RunAndCleanCommand("helm", "install", "argocd", "argo/argo-cd", "--wait")
}

func cleanup() {
	fmt.Println("Cleanup")
	uninstallGloo()

	// "kubectl delete -f ./artifacts/helm-repo.yaml
	cmd := exec.Command("kubectl", "delete", "-f", "./artifacts/helm-repo.yaml")
	cmd.Output()
	// helm uninstall argocd --wait
	kube2e.RunAndCleanCommand("helm", "uninstall", "argocd", "--wait")
}

func deployHelmRepo() {
	fmt.Println("Deploying helm repo")
	// ./artifacts/deploy-helm-server.sh
	kube2e.RunAndCleanCommand("./artifacts/deploy-helm-server.sh")
}
