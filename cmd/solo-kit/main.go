package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/solo-io/solo-kit/pkg/code-generator/cmd"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	if i == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var (
		projectName   string
		resourceNames arrayFlags
	)
	flag.StringVar(&projectName, "name", "", "name of the project")
	flag.Var(&resourceNames, "resource", "name of a resource")
	flag.Parse()

	projectGopath, err := getGopath()
	if err != nil {
		log.Fatal(err)
	}

	projectGopath = filepath.Join(projectGopath, projectName)

	log.Printf("gopath: %v", projectGopath)

	err = cmd.InitProject(projectName, projectGopath, resourceNames...)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("success")
}

func getGopath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	gopath := filepath.Join(os.Getenv("GOPATH"), "src")
	i := strings.Index(wd, gopath)
	if i == -1 {
		return "", errors.Errorf("could not find gopath in %v", wd)
	}
	return wd[i+len(gopath)+1:], nil
}
