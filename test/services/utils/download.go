package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type GetBinaryParams struct {
	Filename    string // the name of the binary on the $PATH or in the docker container
	DockerImage string // the docker image to use if Env or Local are not present
	DockerPath  string // the location of the binary in the docker container, including the filename
	EnvKey      string // the environment var to look at for a user-specified service binary
	TmpDir      string // the temp directory to store a downloaded binary if needed
}

// GetBinary uses the passed params structure to get a binary for the service in the first found of 3 locations:
//
// 1. specified via environment variable
//
// 2. matching binary already on path
//
// 3. download the hard-coded version via docker and extract the binary from that
func GetBinary(params GetBinaryParams) (string, error) {
	// first check if an environment variable was specified for the binary location
	envPath := os.Getenv(params.EnvKey)
	if envPath != "" {
		log.Printf("Using %s specified in environment variable %s: %s", params.Filename, params.EnvKey, envPath)
		return envPath, nil
	}

	// next check if we have a matching binary on $PATH
	localPath, err := exec.LookPath(params.Filename)
	if err == nil {
		log.Printf("Using %s from PATH: %s", params.Filename, localPath)
		return localPath, nil
	}

	// finally, try to grab one from docker
	return extractFileFromDocker(context.Background(), params.TmpDir, params)

}

// extractFileFromDocker by spinning up a container and copying the file out of it
// once its complete kill the container to leave us in a good state
func extractFileFromDocker(ctx context.Context, tmpdir string, params GetBinaryParams) (string, error) {
	log.Printf("Using %s from docker image: %s", params.Filename, params.DockerImage)
	dir := "."
	if tmpdir != "" {
		dir = tmpdir
	}
	destination := fmt.Sprintf("%s/%s", dir, params.Filename)

	// extract the envoy binary from a good docker image
	// to do this spin up an instance, copy out the binary, and then kill the instance
	cli, err := dockclient.NewClientWithOpts(dockclient.FromEnv)
	if err != nil {
		return "", errors.Wrapf(err, "should be able to create a docker client, check that docker is on your path")

	}

	// see https://docs.docker.com/engine/api/sdk/examples/
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: params.DockerImage,
			Cmd:   []string{"/bin/bash", "-c", "exit"},
		},
		nil, nil, nil, "",
	)
	if err != nil {
		return "", errors.Wrapf(err, "check to see that the specified image exists", params.DockerImage)
	}
	// we have to start the container to be able to copy out the envoy binary
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "something may be wrong with your system, inspect the container to see failure")
	}

	inspectionInfo, err := cli.DistributionInspect(ctx, resp.ID, "")
	if err != nil {
		return "", errors.Wrapf(err, "check to see if the image exists and is valid")
	}

	log.Printf("Using Envoy Image with Digest: %v", inspectionInfo.Descriptor.Digest)

	contentReader, stat, err := cli.CopyFromContainer(ctx, resp.ID, params.DockerPath)
	if err != nil {
		return "", errors.Wrapf(err, "check out your container, was it actually a valid path within your docker image?")
	}

	// https://github.com/docker/cli/blob/b1d27091e50595fecd8a2a4429557b70681395b2/cli/command/container/cp.go#L167C1-L172C3
	srcInfo := archive.CopyInfo{
		Path:   params.DockerPath,
		Exists: true,
		IsDir:  stat.Mode.IsDir(),
	}

	err = archive.CopyTo(contentReader, srcInfo, destination)
	if err != nil {
		return "", errors.Wrapf(err, "check to see if your permissions are out of whack")
	}

	// If we dont use force then we may have a race condition and get the following:
	// 	"Stop the container before attempting removal or force remove"
	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		return "", errors.Wrapf(err, "we should only get an error due to permissions or something else beating us to this step")
	}

	return filepath.Join(tmpdir, params.Filename), nil

}
