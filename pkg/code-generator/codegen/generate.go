package codegen

import (
	"bytes"

	"github.com/iancoleman/strcase"
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
		fs, err := generateFilesForResource(res)
		if err != nil {
			return nil, err
		}
		files = append(files, fs...)
	}
	for _, grp := range project.ResourceGroups {
		fs, err := generateFilesForResourceGroup(grp)
		if err != nil {
			return nil, err
		}
		files = append(files, fs...)
	}
	return files, nil
}

func generateFilesForResource(resource *Resource) (Files, error) {
	var v Files
	content, err := generateResourceExtensions(resource)
	if err != nil {
		return nil, err
	}
	v = append(v, File{
		Filename: strcase.ToSnake(resource.Name) + ".go",
		Content:  content,
	})
	content, err = generateResourceClient(resource)
	if err != nil {
		return nil, err
	}
	v = append(v, File{
		Filename: strcase.ToSnake(resource.Name) + "_client.go",
		Content:  content,
	})
	return v, nil
}

func generateFilesForResourceGroup(resource *ResourceGroup) (Files, error) {
	var v Files

	return v, nil
}

func generateFilesForProject(project *Project) (Files, error) {
	var v Files

	return v, nil
}

func generateResourceExtensions(resource *Resource) (string, error) {
	buf := &bytes.Buffer{}
	if err := templates.ResourceExtensionTemplate.Execute(buf, resource); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func generateResourceClient(resource *Resource) (string, error) {
	buf := &bytes.Buffer{}
	if err := templates.ResourceClientTemplate.Execute(buf, resource); err != nil {
		return "", err
	}
	return buf.String(), nil
}
