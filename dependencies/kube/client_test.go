package kube_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"path/filepath"
	"time"

	. "github.com/solo-io/gloo-storage/dependencies/kube"
	restclient "k8s.io/client-go/rest"

	"unicode"

	"encoding/base64"

	"github.com/solo-io/gloo-storage/dependencies"
	. "github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Client", func() {

	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		masterUrl, kubeconfigPath string
		namespace                 string
		syncFreq                  = time.Minute
		cfg                       *restclient.Config
		client                    dependencies.FileStorage
		kube                      kubernetes.Interface
	)
	BeforeEach(func() {
		namespace = RandString(8)
		err := SetupKubeForTest(namespace)
		Must(err)
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		masterUrl = ""
		cfg, err = clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		client, err = NewFileStorage(cfg, namespace, syncFreq)
		Expect(err).NotTo(HaveOccurred())
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		TeardownKube(namespace)
	})
	Describe("create", func() {
		It("errors if the filename does not have a /", func() {
			file := &dependencies.File{
				Name:     "badfilename",
				Contents: []byte{},
			}
			_, err := client.Create(file)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid file ref for kubernetes"))
		})
		It("creates the config map", func() {
			file := &dependencies.File{
				Name:     "good/filename",
				Contents: []byte("hello"),
			}
			f, err := client.Create(file)
			Expect(err).NotTo(HaveOccurred())
			Expect(f).NotTo(BeNil())
			cm, err := kube.CoreV1().ConfigMaps(namespace).Get("good", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.Data).To(HaveLen(1))
			Expect(cm.Data).To(Equal(map[string]string{"filename": "hello"}))
		})
		It("creates the config map for a binary file", func() {
			cmName := "good"
			key := "filename"
			fileRef := cmName + "/" + key
			contents := []byte{1, 2, 3, unicode.MaxASCII + 1}
			file := &dependencies.File{
				Name:     fileRef,
				Contents: contents,
			}
			f, err := client.Create(file)
			Expect(err).NotTo(HaveOccurred())
			Expect(f).NotTo(BeNil())
			cm, err := kube.CoreV1().ConfigMaps(namespace).Get(cmName, v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.Data).To(HaveLen(1))
			Expect(cm.BinaryData).To(HaveLen(0))
			Expect(cm.Data).To(Equal(map[string]string{key: base64.StdEncoding.EncodeToString(contents)}))
		})
	})
})
