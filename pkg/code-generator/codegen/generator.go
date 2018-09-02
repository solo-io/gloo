package codegen

import (
	"bytes"
	"text/template"

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
	content, err := generateResourceFile(resource, templates.ResourceTemplate)
	if err != nil {
		return nil, err
	}
	v = append(v, File{
		Filename: strcase.ToSnake(resource.Name) + ".go",
		Content:  content,
	})
	content, err = generateResourceFile(resource, templates.ResourceClientTemplate)
	if err != nil {
		return nil, err
	}
	v = append(v, File{
		Filename: strcase.ToSnake(resource.Name) + "_client.go",
		Content:  content,
	})
	content, err = generateResourceFile(resource, templates.ResourceClientTestTemplate)
	if err != nil {
		return nil, err
	}
	v = append(v, File{
		Filename: strcase.ToSnake(resource.Name) + "_client_test.go",
		Content:  content,
	})
	content, err = generateResourceFile(resource, templates.ResourceReconcilerTemplate)
	if err != nil {
		return nil, err
	}
	v = append(v, File{
		Filename: strcase.ToSnake(resource.Name) + "_reconciler.go",
		Content:  content,
	})
	return v, nil
}

func generateFilesForResourceGroup(rg *ResourceGroup) (Files, error) {
	var v Files
	for suffix, tmpl := range map[string]*template.Template{
		"_snapshot.go":              templates.ResourceGroupSnapshotTemplate,
		"_snapshot_emitter.go":      templates.ResourceGroupEmitterTemplate,
		"_snapshot_emitter_test.go": templates.ResourceGroupEmitterTestTemplate,
		"_event_loop.go":            templates.ResourceGroupEventLoopTemplate,
	} {
		content, err := generateResourceGroupFile(rg, tmpl)
		if err != nil {
			return nil, err
		}
		v = append(v, File{
			Filename: strcase.ToSnake(rg.GoName) + suffix,
			Content:  content,
		})
	}
	return v, nil
}

func generateFilesForProject(project *Project) (Files, error) {
	var v Files

	return v, nil
}

func generateResourceFile(resource *Resource, tmpl *template.Template) (string, error) {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, resource); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func generateResourceGroupFile(rg *ResourceGroup, tmpl *template.Template) (string, error) {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, rg); err != nil {
		return "", err
	}
	return buf.String(), nil
}
