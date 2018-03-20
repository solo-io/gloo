package kube_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path/filepath"

	"fmt"

	. "github.com/solo-io/gloo-testing/helpers"
	. "github.com/solo-io/gloo/pkg/artifactwatcher/kube"
	"github.com/solo-io/gloo/pkg/log"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("KubeArtifactWatcher", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		masterUrl, kubeconfigPath string
		namespace                 string
	)
	BeforeEach(func() {
		namespace = RandString(8)
		err := SetupKubeForTest(namespace)
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl = ""
	})
	AfterEach(func() {
		TeardownKube(namespace)
	})
	Describe("controller", func() {
		It("watches kube artifacts", func() {
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			watcher, err := NewArtifactWatcher(masterUrl, kubeconfigPath, time.Second, make(chan struct{}))
			Expect(err).NotTo(HaveOccurred())

			// add a artifact
			kubeClient, err := kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			artifact := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "artifact-",
					Namespace:    namespace,
				},
				Data: map[string]string{"username": "me@example.com", "password": "foobar"},
			}

			createdArtifact, err := kubeClient.CoreV1().ConfigMaps(namespace).Create(artifact)
			Expect(err).NotTo(HaveOccurred())

			// give controller time to register
			time.Sleep(time.Second * 2)

			go watcher.TrackArtifacts([]string{createdArtifact.Name})

			select {
			case <-time.After(time.Second * 5):
				Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
			case artifacts := <-watcher.Artifacts():
				Expect(len(artifacts)).To(Equal(1))
				Expect(artifacts[createdArtifact.Name]["username"]).To(Equal([]byte("me@example.com")))
			case err := <-watcher.Error():
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})
})
