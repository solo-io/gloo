package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"

	"os"

	"strings"

	"github.com/pseudomuto/protoc-gen-doc"
)

func main() {
	f := flag.String("f", os.Getenv("GOPATH")+"/src/github.com/solo-io/gloo-api/docs/"+"api.json", "input json file")
	tmplFile := flag.String("t", os.Getenv("GOPATH")+"/src/github.com/solo-io/gloo-api/docs/markdown.tmpl", "template to build from")
	outDir := flag.String("o", os.Getenv("GOPATH")+"/src/github.com/solo-io/gloo-api/docs/v1/", "output dir")
	flag.Parse()
	if err := run(*f, *tmplFile, *outDir); err != nil {
		log.Fatal(err)
	}
}

func run(file, tmplFile, outDir string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var protoDescriptor gendoc.Template
	err = json.Unmarshal(data, &protoDescriptor)
	if err != nil {
		return err
	}

	inputTemplate, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		return err
	}

	for _, protoFile := range protoDescriptor.Files {
		protoFile.Name = strings.TrimSuffix(protoFile.Name, ".proto")
		log.Printf(protoFile.Name)
		tmpl, err := template.New("Proto Doc Template").Funcs(map[string]interface{}{
			"p":        gendoc.PFilter,
			"para":     gendoc.ParaFilter,
			"nobr":     gendoc.NoBrFilter,
			"yamlType": yamlType,
		}).Parse(string(inputTemplate))
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, protoFile)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(outDir+"/"+protoFile.Name+".md", buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func yamlType(longType, label string) string {
	yamlType := func() string {
		switch longType {
		case "string":
			fallthrough
		case "uint32":
			fallthrough
		case "bool":
			fallthrough
		case "int32":
			return longType
		case "Status":
			fallthrough
		case "Metadata":
			return "(read only)"
		}
		return "{" + longType + "}"
	}()
	if label == "repeated" {
		yamlType = "[" + yamlType + "]"
	}
	return yamlType
}
