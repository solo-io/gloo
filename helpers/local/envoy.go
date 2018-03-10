package localhelpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"io/ioutil"

	"github.com/onsi/ginkgo"
)

const defualtEnvoyDockerImage = "soloio/envoy:v0.1.2"

func buildBootstrap(xdsPort uint32) []byte {

	return []byte(fmt.Sprintf(envoyConfigTemplate, xdsPort))
}

const envoyConfigTemplate = `
node:
 cluster: ingress
 id: testnode

static_resources:
  clusters:
  - name: xds_cluster
    connect_timeout: 5.000s
    hosts:
    - socket_address:
        address: localhost
        port_value: %d
    http2_protocol_options: {}
    type: STRICT_DNS

dynamic_resources:
  ads_config:
    api_type: GRPC
    cluster_names:
    - xds_cluster
  cds_config:
    ads: {}
  lds_config:
    ads: {}
  
admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000

`

type EnvoyFactory struct {
	envoypath string

	tmpdir string
}

func NewEnvoyFactory() (*EnvoyFactory, error) {
	envoypath := os.Getenv("ENVOY_BINARY")

	if envoypath != "" {
		return &EnvoyFactory{
			envoypath: envoypath,
		}, nil
	}

	// try to grab one form docker...
	tmpdir, err := ioutil.TempDir(os.Getenv("HELPER_TMP"), "envoy")
	if err != nil {
		return nil, err
	}

	bash := `
set -ex
CID=$(docker run -d  soloio/envoy:v0.1.2 /bin/bash -c exit)
docker cp $CID:/usr/local/bin/envoy .
docker rm $CID
    `
	scriptfile := filepath.Join(tmpdir, "getenvoy.sh")

	ioutil.WriteFile(scriptfile, []byte(bash), 0755)

	cmd := exec.Command("bash", scriptfile)
	cmd.Dir = tmpdir
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &EnvoyFactory{
		envoypath: filepath.Join(tmpdir, "envoy"),
		tmpdir:    tmpdir,
	}, nil
}

func (ef *EnvoyFactory) Clean() error {
	if ef.tmpdir != "" {
		os.RemoveAll(ef.tmpdir)

	}
	return nil
}

type EnvoyInstance struct {
	envoypath    string
	envoycfgpath string
	tmpdir       string
	cmd          *exec.Cmd
}

func (ef *EnvoyFactory) NewEnvoyInstance() (*EnvoyInstance, error) {
	// try to grab one form docker...
	tmpdir, err := ioutil.TempDir(os.Getenv("HELPER_TMP"), "envoy")
	if err != nil {
		return nil, err
	}

	envoyconfigyaml := filepath.Join(tmpdir, "envoyconfig.yaml")

	return &EnvoyInstance{
		envoypath:    ef.envoypath,
		envoycfgpath: envoyconfigyaml,
		tmpdir:       tmpdir,
	}, nil

}

func (ei *EnvoyInstance) Run() error {
	return ei.RunWithPort(8081)
}

func (ei *EnvoyInstance) RunWithPort(port uint32) error {
	err := ioutil.WriteFile(ei.envoycfgpath, buildBootstrap(port), 0644)

	if err != nil {
		return err
	}

	cmd := exec.Command(ei.envoypath, "-c", ei.envoycfgpath, "--v2-config-only")
	cmd.Dir = ei.tmpdir
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	err = cmd.Start()
	if err != nil {
		return err
	}
	ei.cmd = cmd
	return nil
}

func (ei *EnvoyInstance) Binary() string {
	return ei.envoypath
}

func (ei *EnvoyInstance) Clean() error {
	if ei.cmd != nil {
		ei.cmd.Process.Kill()
		ei.cmd.Wait()
	}
	if ei.tmpdir != "" {
		os.RemoveAll(ei.tmpdir)
	}
	return nil
}
