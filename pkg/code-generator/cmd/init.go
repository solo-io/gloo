package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/solo-io/solo-kit/pkg/errors"
)

func InitProject(projectName, projectGopath string, resourceNames ...string) error {
	if projectName == "" {
		return errors.Errorf("project name cannot be empty")
	}
	if projectGopath == "" {
		return errors.Errorf("project gopath cannot be empty")
	}
	if len(resourceNames) == 0 {
		return errors.Errorf("must provide at least 1 resource name")
	}

	path := filepath.Join(projectName, "pkg", "api", "v1")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	// generate.sh
	err = genFile(path,
		"generate.sh",
		generate_sh,
		struct{ ProjectGopath string }{ProjectGopath: projectGopath},
	)
	if err != nil {
		return err
	}

	// generate.go
	err = genFile(path,
		"generate.go",
		generate_go,
		struct{ ProjectGopath string }{ProjectGopath: projectGopath},
	)
	if err != nil {
		return err
	}

	// resource.proto
	for _, resourceName := range resourceNames {
		if len(resourceName) < 3 {
			return errors.Errorf("resource name must be >= 3 characters")
		}
		shortName := strings.ToLower(resourceName[:2])
		pluralName := strings.ToLower(resourceName + "s")
		err := genFile(path,
			strcase.ToLowerCamel(resourceName)+".proto",
			resource_proto,
			struct {
				ProjectGopath string
				ProjectName   string
				ResourceName  string
				PluralName    string
				ShortName     string
			}{
				ProjectGopath: projectGopath,
				ProjectName:   projectName,
				ResourceName:  strcase.ToCamel(resourceName),
				PluralName:    pluralName,
				ShortName:     shortName,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func genFile(path, filename, contents string, data interface{}) error {
	tmpl, err := template.New(filename).Parse(contents)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(path, tmpl.Name()), buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

const generate_go = `package v1

` + `//go:generate {{ .ProjectGopath }}/pkg/api/v1/generate.sh
`

const generate_sh = `#!/usr/bin/env bash

GOGO_OUT_FLAG="--gogo_out=${GOPATH}/src/"
SOLO_KIT_FLAG="--plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out"

ROOT=${GOPATH}/src/{{ .ProjectGopath }}

OUT=${ROOT}/pkg/api/v1/
protoc -I=./ \
    -I=${GOPATH}/src/github.com/gogo/protobuf/ \
    -I=${ROOT} \
    -I=${GOPATH}/src \
    ${GOGO_OUT_FLAG} \
    ${SOLO_KIT_FLAG}=${OUT} \
    *.proto
`

const resource_proto = `syntax = "proto3";
package {{ .ProjectName }}.api.v1;
option go_package = "{{ .ProjectGopath }}/pkg/api/v1";

import "gogoproto/gogo.proto";
option (gogoproto.equal_all) = true;

import "github.com/solo-io/solo-kit/api/v1/metadata.proto";
import "github.com/solo-io/solo-kit/api/v1/status.proto";

/*
@solo-kit:resource
@solo-kit:resource.short_name={{ .ShortName }}
@solo-kit:resource.plural_name={{ .PluralName }}
@solo-kit:resource.group_name={{ .ProjectName }}.api.v1
@solo-kit:resource.version=v1

// TODO: place your comments here
 */
message {{ .ResourceName }} {
    // The Resource-specific config is called a spec.
    {{ .ResourceName }}Spec spec = 2;

    // Status indicates the validation status of the resource. Status is read-only by clients, and set by {{ .ProjectName }} during validation
    core.solo.io.Status status = 6 [(gogoproto.nullable) = false, (gogoproto.moretags) = "testdiff:\"ignore\""];

    // Metadata contains the object metadata for this resource
    core.solo.io.Metadata metadata = 7 [(gogoproto.nullable) = false];
}

// TODO: describe the {{ .ResourceName }}Spec
message {{ .ResourceName }}Spec {
	// TODO: add fields
}
`
