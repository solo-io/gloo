package kube_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path/filepath"

	"fmt"

	. "github.com/solo-io/glue/secrets/watcher/kube"
	. "github.com/solo-io/glue/test/helpers"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("KubeSecretWatcher", func() {

	var (
		masterUrl, kubeconfigPath string
		mkb                       *MinikubeInstance
		namespace                 string
	)
	BeforeSuite(func() {
		namespace = RandString(8)
		mkb = NewMinikube(false, namespace)
		err := mkb.Setup()
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl, err = mkb.Addr()
		Must(err)
	})
	AfterSuite(func() {
		mkb.Teardown()
	})
	Describe("controller", func() {
		It("watches kube secrets", func() {
			cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			watcher, err := NewSecretWatcher(masterUrl, kubeconfigPath, time.Second)
			Expect(err).NotTo(HaveOccurred())

			// add a secret
			kubeClient, err := kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			secret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "secret-",
					Namespace:    namespace,
				},
				Data: map[string][]byte{"username": []byte("me@example.com"), "password": []byte("foobar")},
			}

			createdSecret, err := kubeClient.CoreV1().Secrets(namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())

			// give controller time to register
			time.Sleep(time.Second * 2)

			go watcher.UpdateRefs([]string{createdSecret.Name})

			select {
			case <-time.After(time.Second * 5):
				Expect(fmt.Errorf("expected to have received resource event before 5s")).NotTo(HaveOccurred())
			case secrets := <-watcher.Secrets():
				Expect(len(secrets)).To(Equal(1))
				Expect(secrets[createdSecret.Name]["username"]).To(Equal([]byte("me@example.com")))
			case err := <-watcher.Error():
				Expect(err).To(BeNil())
			}
		})
	})
})
