package file_test

import (
	"os"
	"path"
	"testing"

	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/gluectl/pkg/storage"
	"github.com/solo-io/gluectl/pkg/storage/file"
)

func TestUpstream(t *testing.T) {

	s, err := GetFileStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	x := &gluev1.Upstream{Name: "testcrd", Type: "aws"}
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
	y, err := s.Get(&gluev1.Upstream{Name: "testcrd"}, nil)
	if err != nil {
		t.Fatal("Get failed", err)
	}
	t.Log(y)
	err = s.Delete(y)
	if err != nil {
		t.Fatal("Delete failed", err)
	}
}

func TestList(t *testing.T) {
	s, err := GetFileStorage()
	if err != nil {
		t.Fatal("GetClient failed", err)
	}
	x, err := s.List(&gluev1.Upstream{}, nil)
	for _, z := range x {
		y, ok := z.(*gluev1.Upstream)
		if !ok {
			t.Fatal("List failed - type assertion")
		}
		t.Log(y)
	}
}

func TestUpstreamWatch(t *testing.T) {
	s, err := GetFileStorage()
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

func GetFileStorage() (storage.Storage, error) {
	cfg := path.Join(os.Getenv("HOME"), "glueFile")

	return file.NewFileStorage(cfg, "ant")
}
