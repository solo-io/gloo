package crd_test

import (
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/glue-storage/pkg/storage"
	"github.com/solo-io/glue-storage/pkg/storage/crd"
	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "ant"
)

var (
	kubecfg = path.Join(os.Getenv("HOME"), ".kube/config")
)

func TestUpstream(t *testing.T) {

	s, err := GetCrdStorage()
	Expect(err).NotTo(HaveOccurred())

	x := &gluev1.Upstream{Name: "uscrd", Type: "aws"}
	_, err = s.Create(x)
	Expect(err).NotTo(HaveOccurred())

	spec := make(map[string]interface{})
	spec["region"] = "us-east-1"
	spec["secret"] = "my secret"
	x.Spec = spec
	_, err = s.Update(x)
	Expect(err).NotTo(HaveOccurred())

	y, err := s.Get(&gluev1.Upstream{Name: "uscrd"}, nil)
	Expect(err).NotTo(HaveOccurred())

	t.Log(y)
	lst, err := s.List(&gluev1.Upstream{}, nil)
	Expect(err).NotTo(HaveOccurred())

	for _, z := range lst {
		y, ok := z.(*gluev1.Upstream)
		Expect(ok).To(BeTrue())
		t.Log(y)
	}
	err = s.Delete(y)
	Expect(err).NotTo(HaveOccurred())
}

func TestUpstreamWatch(t *testing.T) {
	s, err := GetCrdStorage()
	Expect(err).NotTo(HaveOccurred())
	c := make(chan bool)
	err = s.Watch(&gluev1.Upstream{}, nil, func(item storage.Item, op storage.WatchOperation) {
		t.Log(item, op)
		if op == storage.DeleteOp {
			c <- true
		}
	})
	Expect(err).NotTo(HaveOccurred())
	TestUpstream(t)
	<-c
	s.WatchStop()
}

func TestVHost(t *testing.T) {

	s, err := GetCrdStorage()
	Expect(err).NotTo(HaveOccurred())
	x := &gluev1.VirtualHost{Name: "vhcrd"}
	_, err = s.Create(x)
	Expect(err).NotTo(HaveOccurred())

	x.Domains = []string{"domain1", "domain2"}
	_, err = s.Update(x)
	Expect(err).NotTo(HaveOccurred())
	y, err := s.Get(&gluev1.VirtualHost{Name: "vhcrd"}, nil)
	Expect(err).NotTo(HaveOccurred())
	t.Log(y)
	lst, err := s.List(&gluev1.VirtualHost{}, nil)
	Expect(err).NotTo(HaveOccurred())

	for _, z := range lst {
		y, ok := z.(*gluev1.VirtualHost)
		Expect(ok).To(BeTrue())
		t.Log(y)
	}
	err = s.Delete(y)
	Expect(err).NotTo(HaveOccurred())
}

func TestVHostWatch(t *testing.T) {
	s, err := GetCrdStorage()
	Expect(err).NotTo(HaveOccurred())
	c := make(chan bool)
	err = s.Watch(&gluev1.VirtualHost{}, nil, func(item storage.Item, op storage.WatchOperation) {
		t.Log(item, op)
		if op == storage.DeleteOp {
			c <- true
		}
	})
	Expect(err).NotTo(HaveOccurred())
	TestVHost(t)
	<-c
	s.WatchStop()
}

func GetCrdStorage() (storage.Storage, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubecfg)
	if err != nil {
		return nil, err
	}
	s, err := crd.NewCrdStorage(cfg, namespace)
	if err != nil {
		return nil, err
	}
	err = s.Register(&gluev1.Upstream{})
	if err != nil {
		return nil, err
	}

	err = s.Register(&gluev1.VirtualHost{})
	if err != nil {
		return nil, err
	}
	return s, nil
}

func init() {
	RegisterFailHandler(Fail)
}
