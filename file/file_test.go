package file_test

import (
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/glue-storage/pkg/storage"
	"github.com/solo-io/glue-storage/pkg/storage/file"
	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
)

const (
	namespace = "ant"
)

var (
	pathCfg = path.Join(os.Getenv("HOME"), "glueFolder")
)

func TestUpstream(t *testing.T) {

	s, err := GetFileStorage()
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
	for _, z := range lst {
		y, ok := z.(*gluev1.Upstream)
		Expect(ok).To(BeTrue())

		t.Log(y)
	}

	err = s.Delete(y)
	Expect(err).NotTo(HaveOccurred())

}

func TestUpstreamWatch(t *testing.T) {
	s, err := GetFileStorage()
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

	s, err := GetFileStorage()
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
	for _, z := range lst {
		y, ok := z.(*gluev1.VirtualHost)
		Expect(ok).To(BeTrue())

		t.Log(y)
	}

	err = s.Delete(y)
	Expect(err).NotTo(HaveOccurred())

}

func TestVHostWatch(t *testing.T) {
	s, err := GetFileStorage()
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

func GetFileStorage() (storage.Storage, error) {
	s, err := file.NewFileStorage(pathCfg, namespace)

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
