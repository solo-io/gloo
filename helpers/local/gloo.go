package localhelpers

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/file"

	"github.com/solo-io/gloo-testing/helpers"
)

type GlooInstance struct {
	srcpath  string
	gloopath string
	wasbuilt bool

	tmpdir string
	store  storage.Interface
	cmd    *exec.Cmd
}

func NewGlooInstance() (*GlooInstance, error) {

	gi, err := newGlooInstance2()

	if err != nil {
		return nil, err
	}
	gi.initStorage()

	if err != nil {
		return nil, err
	}
	return gi, nil
}

func newGlooInstance2() (*GlooInstance, error) {
	gloopath := os.Getenv("GLOO_BINARY")

	tmpdir, err := ioutil.TempDir(os.Getenv("HELPER_TMP"), "gloo")
	if err != nil {
		return nil, err
	}

	if gloopath != "" {
		return &GlooInstance{
			gloopath: gloopath,
			tmpdir:   tmpdir,
		}, nil
	}
	srcpath := filepath.Join(helpers.SoloDirectory(), "gloo")
	gloopath = filepath.Join(srcpath, "gloo")
	gi := &GlooInstance{
		srcpath:  srcpath,
		gloopath: gloopath,
		tmpdir:   tmpdir,
	}
	gi.build()
	return gi, nil

}

func (gi *GlooInstance) EnvoyPort() uint32 {
	return 8080
}

func (gi *GlooInstance) build() error {
	if gi.srcpath == "" {
		if gi.gloopath == "" {
			return errors.New("can't build gloo and none provided")
		}
		return nil
	}
	gi.wasbuilt = true

	cmd := exec.Command("go", "build", "-v", "-i", "-gcflags", "-N -l", "-o", "gloo", "main.go")
	cmd.Dir = gi.srcpath
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (gi *GlooInstance) AddUpstream(u *v1.Upstream) error {
	_, err := gi.store.V1().Upstreams().Create(u)
	return err
}

func (gi *GlooInstance) AddVhost(u *v1.VirtualHost) error {
	_, err := gi.store.V1().VirtualHosts().Create(u)
	return err
}

func (gi *GlooInstance) initStorage() error {

	dir := gi.tmpdir
	client, err := file.NewStorage(filepath.Join(dir, "_gloo_config"), time.Hour)
	if err != nil {
		return errors.New("failed to start file config watcher for directory " + dir)
	}
	err = client.V1().Register()
	if err != nil {
		return errors.New("failed to register file config watcher for directory " + dir)
	}
	gi.store = client
	return nil

}

func (gi *GlooInstance) Run() error {
	cmd := exec.Command(gi.gloopath,
		"--storage.type=file",
		"--storage.refreshrate=1s",
		"--secrets.type=file",
		"--secrets.refreshrate=1s",
		"--xds.port=8081",
	)

	cmd.Dir = gi.tmpdir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	gi.cmd = cmd
	return nil
}

func (gi *GlooInstance) Clean() error {
	if gi.cmd != nil {
		gi.cmd.Process.Kill()
		gi.cmd.Wait()
	}
	if gi.tmpdir != "" {
		defer os.RemoveAll(gi.tmpdir)

	}

	return nil
}
