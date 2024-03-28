package edge_classic_e2e

import (
	"context"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/gloo/test/snapshot/utils"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestEdgeGateway(t *testing.T) {
	helpers.RegisterGlooDebugLogPrintHandlerAndClearLogs()
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gloo Gateway Suite")
}

var (
	ctx    context.Context
	cancel context.CancelFunc

	clusterName, kubeCtx string
	clientScheme         *runtime.Scheme
	resourceClientset    *kube2e.KubeResourceClientSet
	kubeClient           client.Client
	runner               snapshot.TestRunner
)

var _ = BeforeSuite(StartTestHelper)

func StartTestHelper() {
	ctx, cancel = context.WithCancel(context.Background())

	clusterName = os.Getenv("CLUSTER_NAME")
	kubeCtx = fmt.Sprintf("kind-%s", clusterName)

	// set up resource client and kubeclient
	var err error
	resourceClientset, err = kube2e.NewDefaultKubeResourceClientSet(ctx)
	Expect(err).NotTo(HaveOccurred(), "can create kube resource client set")

	clientScheme, err = utils.BuildClientScheme()
	Expect(err).NotTo(HaveOccurred(), "can build client scheme")

	kubeClient, err = utils.GetClient(kubeCtx, clientScheme)
	Expect(err).NotTo(HaveOccurred(), "can create client")

	runner = snapshot.TestRunner{
		Name:      "classic-apis",
		Scheme:    clientScheme,
		ClientSet: resourceClientset,
		Client:    kubeClient,
	}
}

var _ = AfterSuite(func() {
	defer cancel()
})
