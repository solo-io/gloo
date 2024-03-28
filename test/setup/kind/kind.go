package kind

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cmd"
	create "sigs.k8s.io/kind/pkg/cmd/kind/create/cluster"
	get "sigs.k8s.io/kind/pkg/cmd/kind/get/clusters"
	load "sigs.k8s.io/kind/pkg/cmd/kind/load/docker-image"

	"github.com/solo-io/gloo/test/setup/helpers"
)

var (
	ErrNotFound      = errors.New("cluster not found")
	ErrAlreadyExists = func(name string) string {
		return fmt.Sprintf("failed to create cluster: node(s) already exist for a cluster with the name %q", name)
	}
)

func Get(cluster *v1alpha4.Cluster) error {
	if cluster == nil {
		return nil
	}

	buf := &bytes.Buffer{}
	cmd := get.NewCommand(cmd.NewLogger(), cmd.IOStreams{
		Out: buf,
	})

	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		return err
	}

	if !strings.Contains(buf.String(), cluster.Name) {
		return ErrNotFound
	}
	return nil
}

func Create(cluster *v1alpha4.Cluster) error {
	if cluster == nil {
		return nil
	}

	timerFn := helpers.TimerFunc(fmt.Sprintf("[%s] kind create", cluster.Name))
	defer timerFn()

	buf := bytes.Buffer{}
	if err := yaml.NewEncoder(&buf).Encode(cluster); err != nil {
		return err
	}

	cmd := create.NewCommand(cmd.NewLogger(), cmd.IOStreams{
		In:     &buf,
		Out:    io.Discard,
		ErrOut: io.Discard,
	})
	cmd.SetArgs([]string{"--config=-"})

	if err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}

func LoadImage(imageRef string, cluster string) error {
	if imageRef == "" {
		return nil
	}

	cmd := load.NewCommand(cmd.NewLogger(), cmd.StandardIOStreams())
	cmd.SetArgs([]string{imageRef, "--name", cluster})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	return cmd.Execute()
}
