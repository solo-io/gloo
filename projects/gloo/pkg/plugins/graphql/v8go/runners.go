package v8go

import (
	_ "embed"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/log"
	v1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/enterprise/graphql/v1"
	js "github.com/solo-io/solo-projects/projects/gloo/pkg/utils/js"
)

// to create the bundled js files please "make install-node-packages" or "build-stitching-bundles"

//go:embed stitching_bundled.js
var stitchingBundledJs string

//go:embed schema-diff_bundled.js
var schemaDiffJS string

// singleStitchRunner is the singleton for the StitchingScriptRunner, this is because the runner only needs to get cached
// once.
var singleStitchRunner *StitchingScriptRunner

// StitchingScriptRunner is the GraphQL stitching runner. This will take in the appropriate input for stitching and for
// schema diffs.
type StitchingScriptRunner struct {
	// stitchRunner is the v8go stitch program runner
	stitchRunner js.Runner
	// diffRunner is the v8go schema diff program runner
	diffRunner js.Runner
}

const (
	// Timeout string which sets a timeout on how long the stitching script runs
	// default is 10s
	StitchingScriptTimeoutVar = "GRAPHQL_STITCHING_SCRIPT_TIMEOUT"
)

// GetStitchingScriptRunner will create a new v8go runner for stitching and schema diff
// To set the timeout, use the GRAPHQL_STITCHING_SCRIPT_TIMEOUT environment variable. The default is 10 seconds.
func GetStitchingScriptRunner() (*StitchingScriptRunner, error) {
	if singleStitchRunner == nil {
		stitchRunner, err := js.NewV8RunnerInputOutput("stitching_bundled.js", stitchingBundledJs)
		if err != nil {
			return nil, err
		}
		diffRunner, err := js.NewV8RunnerInputOutput("schema-diff_bundled.js", schemaDiffJS)
		if err != nil {
			return nil, err
		}
		timeout := time.Second * 10
		if timeoutStr, ok := os.LookupEnv(StitchingScriptTimeoutVar); ok {
			timeout, err = time.ParseDuration(timeoutStr)
			if err != nil {
				return nil, eris.Wrapf(err, "unable to parse %s as a duration for the env var %s", timeoutStr, StitchingScriptTimeoutVar)
			}
		}
		stitchRunner.SetTimeout(timeout)
		diffRunner.SetTimeout(timeout)
		log.Debugf("graphql plugin: setting stitching execution timeout to %s", timeout.String())
		singleStitchRunner = &StitchingScriptRunner{
			stitchRunner: stitchRunner,
			diffRunner:   diffRunner,
		}
	}
	return singleStitchRunner, nil
}

// RunSchemaDiff will run the v8go schema diff program
func (r *StitchingScriptRunner) RunSchemaDiff(diffInput *v1.GraphQLInspectorDiffInput) (*v1.GraphQLInspectorDiffOutput, error) {
	inputBytes, err := proto.Marshal(diffInput)
	b64 := base64.StdEncoding.EncodeToString(inputBytes)
	if err != nil {
		return nil, eris.Wrap(err, "error encoding diff input")
	}
	output, err := r.diffRunner.Run(b64)
	if err != nil {
		return nil, eris.Wrap(err, "error running schema diff runner")
	}
	decodedStdOutString, err := base64.StdEncoding.DecodeString(output)
	if err != nil {
		return nil, eris.Wrapf(err, "error decoding %s from base64 protobuf: [%s]", output, diffInput)
	}
	stitchingInfoOut := &v1.GraphQLInspectorDiffOutput{}
	err = proto.Unmarshal(decodedStdOutString, stitchingInfoOut)
	if err != nil {
		return nil, eris.Wrapf(err, "unable to unmarshal graphql tools output to Go type: [%s]", string(decodedStdOutString))
	}
	return stitchingInfoOut, nil
}

// RunStitching will run the v8go sticthing program
func (r *StitchingScriptRunner) RunStitching(input []byte) (*v1.GraphQLToolsStitchingOutput, error) {
	output, err := r.stitchRunner.Run(base64.StdEncoding.EncodeToString(input))
	if err != nil {
		return nil, eris.Wrap(err, "error running stitching")
	}
	decodedStdOutString, err := base64.StdEncoding.DecodeString(output)
	if err != nil {
		return nil, eris.Wrapf(err, "error decoding %s from stitching base64 protobuf", output)
	}
	stitchingInfoOut := &v1.GraphQLToolsStitchingOutput{}
	err = proto.Unmarshal(decodedStdOutString, stitchingInfoOut)
	if err != nil {
		return nil, eris.Wrapf(err, "unable to unmarshal stitching graphql tools output to Go type: [%s]", string(decodedStdOutString))
	}
	return stitchingInfoOut, nil
}
