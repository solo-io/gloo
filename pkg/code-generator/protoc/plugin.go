package protoc

import (
	"bytes"
	"encoding/json"
	"github.com/solo-io/solo-kit/pkg/errors"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/solo-io/solo-kit/pkg/code-generator/codegen"
	"github.com/solo-io/solo-kit/pkg/utils/log"
)

// plugin is an implementation of protokit.Plugin
type Plugin struct{}

func (p *Plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	log.DefaultOut = &bytes.Buffer{}
	if false {
		log.DefaultOut = os.Stderr
	}

	log.Printf("parsing request %v", req.FileToGenerate, req.GetParameter())
	paramString := req.GetParameter()
	if paramString == "" {
		return nil, errors.Errorf(`must provide project params via --solo-kit_out={"project_file": "${PWD}/project.json:${OUTDIR}"}`)
	}

	var params codegen.Params
	if err := json.Unmarshal([]byte(paramString), &params); err != nil {
		return nil, errors.Wrapf(err, "failed to parse  param string %v", paramString)
	}

	// append descriptors if they already exist
	collectedDescriptorsFile := params.ProjectFile + ".descriptors"
	descriptorBytes, err := ioutil.ReadFile(collectedDescriptorsFile)
	if err == nil {
		var previousDescriptors plugin_go.CodeGeneratorRequest
		if err := proto.Unmarshal(descriptorBytes, &previousDescriptors); err != nil {
			return nil, errors.Wrapf(err, "unmarshalling previous request")
		}
		// append all unique
		for _, prev := range previousDescriptors.ProtoFile {
			var duplicate bool
			for _, desc := range req.ProtoFile {
				if desc.GetName() == prev.GetName() {
					duplicate = true
					break
				}
			}
			if duplicate {
				continue
			}
			req.ProtoFile = append(req.ProtoFile, prev)
			for _, genFile := range previousDescriptors.FileToGenerate {
				if genFile == prev.GetName() {
					req.FileToGenerate = append(req.FileToGenerate, genFile)
					break
				}
			}
		}
	}

	if params.CollectionRun {
		collectedDescriptorsBytes, err := proto.Marshal(req)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(collectedDescriptorsFile, collectedDescriptorsBytes, 0644); err != nil {
			return nil, err
		}
		// only purpose in the collection run is to build up our request
		return nil, nil
	} else {
		os.Remove(collectedDescriptorsFile)
	}

	project, err := codegen.ParseRequest(params, req)
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
