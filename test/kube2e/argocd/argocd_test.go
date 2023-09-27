package argocd_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/k8s-utils/testutils/helper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const namespace = defaults.GlooSystem

var _ = Describe("Kube2e: ArgoCD", func() {

	var (
		testHelper *helper.SoloTestHelper
	)

	BeforeEach(func() {

		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		testHelper, err = helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
			defaults.RootDir = filepath.Join(cwd, "../../..")
			defaults.HelmChartName = "gloo"
			defaults.InstallNamespace = namespace
			defaults.Verbose = true
			return defaults
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("Tests the Lifecycle", func() {
		installGloo()
		checkGlooHealthy(testHelper)
		uninstallGloo()
	})
})

func installGloo() {
	uninstallGloo()

	repo := "http://helm-repo/"

	// argocd --core app create gloo \
	// --repo http://helm-repo/ --helm-chart gloo --revision $VERSION
	// --dest-namespace gloo-system --dest-server https://kubernetes.default.svc
	// --sync-option CreateNamespace=true --upsert --helm-set gatewayProxies.gatewayProxy.service.type=ClusterIP

	runAndCleanCommand("argocd", "--core", "app", "create", "gloo",
		"--repo", repo, "--helm-chart", "gloo", "--revision", version,
		"--dest-namespace", "gloo-system", "--dest-server", "https://kubernetes.default.svc",
		"--sync-option", "CreateNamespace=true", "--upsert", "--helm-set", "gatewayProxies.gatewayProxy.service.type=ClusterIP")

	// argocd --core app sync gloo
	runAndCleanCommand("argocd", "--core", "app", "sync", "gloo")
}

func uninstallGloo() {
	// argocd --core app delete gloo -y
	cmd := exec.Command("argocd", "--core", "app", "delete", "gloo", "-y")
	cmd.Output()
}

func runAndCleanCommand(name string, arg ...string) []byte {
	cmd := exec.Command(name, arg...)
	b, err := cmd.Output()
	// for debugging in Cloud Build
	if err != nil {
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
	deploymentNames := []string{"gloo", "discovery", "gateway-proxy"}
	for _, deploymentName := range deploymentNames {
		runAndCleanCommand("kubectl", "rollout", "status", "deployment", "-n", testHelper.InstallNamespace, deploymentName)
	}
	kube2e.GlooctlCheckEventuallyHealthy(2, testHelper, "90s")
}
