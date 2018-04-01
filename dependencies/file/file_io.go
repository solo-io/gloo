package file

import (
	"io/ioutil"

	"path/filepath"

	"os"

	"strings"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-storage/dependencies"
)

func writeFile(dir string, file *dependencies.File) error {
	if strings.Contains(file.Name, "/") {
		return errors.Errorf("file name cannot contain '/': %v", file.Name)
	}
	return ioutil.WriteFile(filepath.Join(dir, file.Name), file.Contents, 0644)
}

func readFile(dir, filename string) (*dependencies.File, error) {
	data, err := ioutil.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return nil, errors.Errorf("error reading file: %v", err)
	}
	return &dependencies.File{
		Name:     filename,
		Contents: data,
	}, nil
}

func deleteFile(dir, filename string) error {
	return os.Remove(filepath.Join(dir, filename))
}
