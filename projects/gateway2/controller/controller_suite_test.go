package controller_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/config"

	"github.com/solo-io/gloo/pkg/schemes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	"github.com/solo-io/gloo/projects/gateway2/controller"
	"github.com/solo-io/gloo/projects/gateway2/extensions"
	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	api "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc

	kubeconfig string

	gwClasses = sets.New(gatewayClassName, altGatewayClassName)
)

const (
	gatewayClassName      = "clsname"
	altGatewayClassName   = "clsname-alt"
	gatewayControllerName = "controller/name"
)

func getAssetsDir() string {
	assets := ""
	if os.Getenv("KUBEBUILDER_ASSETS") == "" {
		// set default if not user provided
		out, err := exec.Command("sh", "-c", "make -sC $(dirname $(go env GOMOD))/projects/gateway2 envtest-path").CombinedOutput()
		fmt.Fprintln(GinkgoWriter, "out:", string(out))
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		assets = strings.TrimSpace(string(out))
	}
	return assets
}

var _ = BeforeSuite(func() {
	log.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "crds"),
			filepath.Join("..", "..", "..", "install", "helm", "gloo", "crds"),
		},
		ErrorIfCRDPathMissing: true,
		// set assets dir so we can run without the makefile
		BinaryAssetsDirectory: getAssetsDir(),
		// web hook to add cluster ips to services

	}
	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())
	scheme := schemes.GatewayScheme()
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	webhookInstallOptions := &testEnv.WebhookInstallOptions
	mgrOpts := ctrl.Options{
		Scheme: scheme,
		WebhookServer: webhook.NewServer(webhook.Options{
			Host:    webhookInstallOptions.LocalServingHost,
			Port:    webhookInstallOptions.LocalServingPort,
			CertDir: webhookInstallOptions.LocalServingCertDir,
		}),
		Controller: config.Controller{
			// see https://github.com/kubernetes-sigs/controller-runtime/issues/2937
			// in short, our tests reuse the same name (reasonably so) and the controller-runtime
			// package does not reset the stack of controller names between tests, so we disable
			// the name validation here.
			SkipNameValidation: ptr.To(true),
		},
	}
	mgr, err := ctrl.NewManager(cfg, mgrOpts)
	Expect(err).ToNot(HaveOccurred())

	kubeconfig = generateKubeConfiguration(cfg)
	mgr.GetLogger().Info("starting manager", "kubeconfig", kubeconfig)

	exts, err := extensions.NewK8sGatewayExtensions(ctx, extensions.K8sGatewayExtensionsFactoryParameters{
		Mgr: mgr,
	})
	Expect(err).ToNot(HaveOccurred())
	cfg := controller.GatewayConfig{
		Mgr:            mgr,
		ControllerName: gatewayControllerName,
		GWClasses:      gwClasses,
		AutoProvision:  true,
		Kick:           func(ctx context.Context) { return },
		Extensions:     exts,
	}
	err = controller.NewBaseGatewayController(ctx, cfg)
	Expect(err).ToNot(HaveOccurred())

	err = cfg.Mgr.GetFieldIndexer().IndexField(ctx, &apixv1a1.XListenerSet{}, query.ListenerSetTargetField, query.IndexerByObjType)
	Expect(err).NotTo(HaveOccurred())

	for class := range gwClasses {
		err = k8sClient.Create(ctx, &api.GatewayClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: class,
			},
			Spec: api.GatewayClassSpec{
				ControllerName: api.GatewayController(gatewayControllerName),
				ParametersRef: &api.ParametersReference{
					Group:     api.Group(v1alpha1.GroupVersion.Group),
					Kind:      api.Kind("GatewayParameters"),
					Name:      wellknown.DefaultGatewayParametersName,
					Namespace: ptr.To(api.Namespace("default")),
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
	}

	err = k8sClient.Create(ctx, &v1alpha1.GatewayParameters{
		ObjectMeta: metav1.ObjectMeta{
			Name:      wellknown.DefaultGatewayParametersName,
			Namespace: "default",
		},
		Spec: v1alpha1.GatewayParametersSpec{
			Kube: &v1alpha1.KubernetesProxyConfig{
				Service: &v1alpha1.Service{
					Type: ptr.To(corev1.ServiceTypeLoadBalancer),
				},
				Istio: &v1alpha1.IstioIntegration{},
			},
		},
	})
	Expect(err).NotTo(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = mgr.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
	if kubeconfig != "" {
		os.Remove(kubeconfig)
	}
})

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

func generateKubeConfiguration(restconfig *rest.Config) string {
	clusters := make(map[string]*clientcmdapi.Cluster)
	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	contexts := make(map[string]*clientcmdapi.Context)

	clusterName := "cluster"
	clusters[clusterName] = &clientcmdapi.Cluster{
		Server:                   restconfig.Host,
		CertificateAuthorityData: restconfig.CAData,
	}
	authinfos[clusterName] = &clientcmdapi.AuthInfo{
		ClientKeyData:         restconfig.KeyData,
		ClientCertificateData: restconfig.CertData,
	}
	contexts[clusterName] = &clientcmdapi.Context{
		Cluster:   clusterName,
		Namespace: "default",
		AuthInfo:  clusterName,
	}

	clientConfig := clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters:   clusters,
		Contexts:   contexts,
		// current context must be mgmt cluster for now, as the api server doesn't have context configurable.
		CurrentContext: "cluster",
		AuthInfos:      authinfos,
	}
	// create temp file
	tmpfile, err := os.CreateTemp("", "ggii_envtest_*.kubeconfig")
	Expect(err).NotTo(HaveOccurred())
	tmpfile.Close()
	err = clientcmd.WriteToFile(clientConfig, tmpfile.Name())
	Expect(err).NotTo(HaveOccurred())
	return tmpfile.Name()
}
