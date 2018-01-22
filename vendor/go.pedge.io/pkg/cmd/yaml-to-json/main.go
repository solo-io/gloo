/*
Package main implements the yaml2json command-line tool, which takes
a YAML file as stdin, and outputs a JSON file to stdout.
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"go.pedge.io/pkg/yaml"
)

func main() {
	if err := do(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func do() error {
	var pretty bool
	flag.BoolVar(&pretty, "pretty", false, "Make the JSON output pretty.")
	flag.Parse()

	yamlData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	jsonData, err := pkgyaml.ToJSON(yamlData, pkgyaml.ToJSONOptions{Pretty: pretty})
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(jsonData)
	return err
}
