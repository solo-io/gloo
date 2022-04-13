package translation

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/golang/protobuf/proto"
	"github.com/rotisserie/eris"
	enterprisev1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/enterprise/graphql/v1"
)

func GetSchemaDiff(input *enterprisev1.GraphQLInspectorDiffInput) (*enterprisev1.GraphQLInspectorDiffOutput, error) {
	inputBytes, err := proto.Marshal(input)
	if err != nil {
		return nil, eris.Wrapf(err, "error marshaling to binary data")
	}
	jsPath := GetGraphqlJsRoot()
	cmd := exec.Command("node", jsPath+"schema-diff.js", base64.StdEncoding.EncodeToString(inputBytes))

	protoPath, err := GetGraphqlProtoRoot()
	if err != nil {
		return nil, err
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", GraphqlProtoRootEnvVar, protoPath))
	stdOutBuf := bytes.NewBufferString("")
	stdErrBuf := bytes.NewBufferString("")
	cmd.Stdout = stdOutBuf
	cmd.Stderr = stdErrBuf
	err = cmd.Run()
	stdOutString := stdOutBuf.String()
	if err != nil {
		return nil, eris.Wrapf(err, "error running schema diff script, stdout: %s, \nstderr: %s", stdOutString, stdErrBuf.String())
	}
	if len(stdOutString) == 0 {
		return nil, eris.Errorf("error running schema diff script, no schema diff generated, \nstderr: %s", stdErrBuf.String())
	}
	decodedStdOutString, err := base64.StdEncoding.DecodeString(stdOutBuf.String())
	if err != nil {
		return nil, eris.Wrapf(err, "error decoding %s from base64 protobuf", stdOutBuf.String())
	}

	output := &enterprisev1.GraphQLInspectorDiffOutput{}
	err = proto.Unmarshal(decodedStdOutString, output)
	if err != nil {
		return nil, eris.Wrap(err, "unable to unmarshal graphql inspector output to Go type")
	}
	return output, nil
}
