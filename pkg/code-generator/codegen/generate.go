package codegen

import (
	"bytes"

	"github.com/solo-io/solo-kit/pkg/code-generator/templates"
)

type File struct {
	Filename string
	Content  string
}

type Files []File

func GenerateFiles(project *Project) (Files, error) {
	files, err := generateFilesForProject(project)
	if err != nil {
		return nil, err
	}
	for _, res := range project.Resources {
		fs, err := generateFilesForResource(project, res)
		if err != nil {
			return nil, err
		}
		files = append(files, fs...)
	}
	for _, grp := range project.ResourceGroups {
		fs, err := generateFilesForResourceGroup(project, grp)
		if err != nil {
			return nil, err
		}
		files = append(files, fs...)
	}
	return files, nil
}

func generateFilesForResource(project *Project, resource *Resource) (Files, error) {
	var v Files

	return v, nil
}

func generateFilesForResourceGroup(project *Project, resource *ResourceGroup) (Files, error) {
	var v Files

	return v, nil
}

func generateFilesForProject(project *Project) (Files, error) {
	var v Files

	return v, nil
}

func generateResourceClient(resource *Resource) (string, error) {
	buf := &bytes.Buffer{}
	if err := templates.ResourceClientTemplate.Execute(buf, resource); err != nil {
		return "", err
	}
	return buf.String(), nil
}
