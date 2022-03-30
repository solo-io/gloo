package stitching

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/golang/protobuf/proto"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
)

var (
	// Env var that has the path to the index.js file that runs the stitching code
	StitchingIndexFilePathEnvVar = "STITCHING_PATH"
	// Env var that has the path to the proto dependencies that the stitching index.js file requires
	StitchingProtoDependenciesPathEnvVar = "STITCHING_PROTO_DIR"
)

func ProcessStitchingInfo(pathToStitchingJsFile string, schemas *v1alpha1.GraphQLToolsStitchingInput) (*v1alpha1.GraphQlToolsStitchingOutput, error) {
	schemasBytes, err := proto.Marshal(schemas)
	if err != nil {
		return nil, eris.Wrapf(err, "error marshaling to binary data")
	}
	// This is the default path
	var stitchingPath = pathToStitchingJsFile
	// Used for local testing and unit/e2e tests
	if path := os.Getenv(StitchingIndexFilePathEnvVar); path != "" {
		stitchingPath = path
	}
	cmd := exec.Command("node", stitchingPath, base64.StdEncoding.EncodeToString(schemasBytes))
	// Set the environment variable STITCHING_PROTO_IMPORT_PATH for the node file to know where to import dependencies from.
	protoDirPath := "/usr/local/bin/js/proto/github.com/solo-io/solo-apis/api/gloo/gloo"
	if protoDir := os.Getenv(StitchingProtoDependenciesPathEnvVar); protoDir != "" {
		// JS needs the absolute path to the stitching proto dir, so we join path to current dir + path from repository root
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, eris.Wrap(err, "unable to get current directory path for running stitching script")
		}
		protoDirPath = path.Join(currentDir, os.Getenv(StitchingProtoDependenciesPathEnvVar))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", StitchingProtoDependenciesPathEnvVar, protoDirPath))
	stdOutBuf := bytes.NewBufferString("")
	stdErrBuf := bytes.NewBufferString("")
	cmd.Stdout = stdOutBuf
	cmd.Stderr = stdErrBuf
	err = cmd.Run()
	stdOutString := stdOutBuf.String()
	if err != nil {
		return nil, eris.Wrapf(err, "error running stitching info generation, stdout: %s, \nstderr: %s", stdOutString, stdErrBuf.String())
	}
	if len(stdOutString) == 0 {
		return nil, eris.Errorf("error running stitching info generation, no stitching info generated, \nstderr: %s", stdErrBuf.String())
	}
	decodedStdOutString, err := base64.StdEncoding.DecodeString(stdOutBuf.String())
	if err != nil {
		return nil, eris.Wrapf(err, "error decoding %s from base64 protobuf", stdOutBuf.String())
	}
	stitchingInfoOut := &v1alpha1.GraphQlToolsStitchingOutput{}
	err = proto.Unmarshal(decodedStdOutString, stitchingInfoOut)
	if err != nil {
		return nil, eris.Wrap(err, "unable to unmarshal graphql tools output to Go type")
	}
	return stitchingInfoOut, nil
}
