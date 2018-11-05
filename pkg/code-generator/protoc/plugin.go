package protoc

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/solo-io/solo-kit/pkg/errors"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/solo-io/solo-kit/pkg/code-generator/codegen"
	"github.com/solo-io/solo-kit/pkg/utils/log"
)

// plugin is an implementation of protokit.Plugin
type Plugin struct {
	OutputDescriptors bool
}

func (p *Plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	log.DefaultOut = &bytes.Buffer{}
	if false {
		log.DefaultOut = os.Stderr
	}

	log.Printf("received request files: %v | params: %v", req.FileToGenerate, req.GetParameter())
	projectFile := req.GetParameter()
	if projectFile == "" {
		return nil, errors.Errorf(`must provide project file via --solo-kit_out=${PWD}/project.json:${OUT}`)
	}

	// if OutputDescriptors==true save request as a descriptors file and exit
	if p.OutputDescriptors {
		collectedDescriptorsBytes, err := proto.Marshal(req)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal %+v", req)
		}
		if err := ioutil.WriteFile(projectFile+".descriptors", collectedDescriptorsBytes, 0644); err != nil {
			return nil, errors.Wrapf(err, "failed to write %v", projectFile+".descriptors")
		}
	}

	project, err := codegen.ParseRequest(projectFile, req)
	if err != nil {
		return nil, err
	}

	code, err := codegen.GenerateFiles(project)
	if err != nil {
		return nil, err
	}

	log.Printf("%v", project)
	log.Printf("%v", code)

	resp := new(plugin_go.CodeGeneratorResponse)

	for _, file := range code {
		resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
			Name:    proto.String(file.Filename),
			Content: proto.String(file.Content),
		})
	}

	return resp, nil
}
