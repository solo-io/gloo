package crd_test

import (
	"os"
	"path"
	"testing"

	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/gluectl/pkg/storage"
	"github.com/solo-io/gluectl/pkg/storage/crd"
	"k8s.io/client-go/tools/clientcmd"
)

func TestUpstream(t *testing.T) {

	s, err := GetCrdStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	x := &gluev1.Upstream{Name: "uscrd", Type: "aws"}
	_, err = s.Create(x)
	if err != nil {
		t.Fatal("Create failed", err)
	}
	spec := make(map[string]interface{})
	spec["region"] = "us-east-1"
	spec["secret"] = "my secret"
	x.Spec = spec
	_, err = s.Update(x)
	if err != nil {
		t.Fatal("Update failed", err)
	}
	y, err := s.Get(&gluev1.Upstream{Name: "uscrd"}, nil)
	if err != nil {
		t.Fatal("Get failed", err)
	}
	t.Log(y)
	lst, err := s.List(&gluev1.Upstream{}, nil)
	for _, z := range lst {
		y, ok := z.(*gluev1.Upstream)
		if !ok {
			t.Fatal("List failed - type assertion")
		}
		t.Log(y)
	}
	err = s.Delete(y)
	if err != nil {
		t.Fatal("Delete failed", err)
	}
}

func TestUpstreamWatch(t *testing.T) {
	s, err := GetCrdStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	c := make(chan bool)
	err = s.Watch(&gluev1.Upstream{}, nil, func(item storage.Item, op storage.WatchOperation) {
		t.Log(item, op)
		if op == storage.Delete {
			c <- true
		}
	})
	if err != nil {
		t.Fatal("Watch failed", err)
	}
	TestUpstream(t)
	<-c
	s.WatchStop()
}

func TestVHost(t *testing.T) {

	s, err := GetCrdStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	x := &gluev1.VirtualHost{Name: "vhcrd"}
	_, err = s.Create(x)
	if err != nil {
		t.Fatal("Create failed", err)
	}

	x.Domains = []string{"domain1", "domain2"}
	_, err = s.Update(x)
	if err != nil {
		t.Fatal("Update failed", err)
	}
	y, err := s.Get(&gluev1.VirtualHost{Name: "vhcrd"}, nil)
	if err != nil {
		t.Fatal("Get failed", err)
	}
	t.Log(y)
	lst, err := s.List(&gluev1.VirtualHost{}, nil)
	for _, z := range lst {
		y, ok := z.(*gluev1.VirtualHost)
		if !ok {
			t.Fatal("List failed - type assertion")
		}
		t.Log(y)
	}
	err = s.Delete(y)
	if err != nil {
		t.Fatal("Delete failed", err)
	}
}

func TestVHostWatch(t *testing.T) {
	s, err := GetCrdStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	c := make(chan bool)
	err = s.Watch(&gluev1.VirtualHost{}, nil, func(item storage.Item, op storage.WatchOperation) {
		t.Log(item, op)
		if op == storage.Delete {
			c <- true
		}
	})
	if err != nil {
		t.Fatal("Watch failed", err)
	}
	TestVHost(t)
	<-c
	s.WatchStop()
}

func GetCrdStorage() (storage.Storage, error) {
	kubecfg := path.Join(os.Getenv("HOME"), ".kube/config")
	cfg, err := clientcmd.BuildConfigFromFlags("", kubecfg)
	if err != nil {
		return nil, err
	}
	s, err := crd.NewCrdStorage(cfg, "ant")
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
