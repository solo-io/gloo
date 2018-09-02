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
	for suffix, tmpl := range map[string]*template.Template{
		".go":             templates.ResourceTemplate,
		"_client.go":      templates.ResourceClientTemplate,
		"_client_test.go": templates.ResourceClientTestTemplate,
		"_reconciler.go":  templates.ResourceReconcilerTemplate,
	} {
		content, err := generateResourceFile(resource, tmpl)
		if err != nil {
			return nil, err
		}
		v = append(v, File{
			Filename: strcase.ToSnake(resource.Name) + suffix,
			Content:  content,
		})
	}
	return v, nil
}

func generateFilesForResourceGroup(rg *ResourceGroup) (Files, error) {
	var v Files
	for suffix, tmpl := range map[string]*template.Template{
		"_snapshot.go":              templates.ResourceGroupSnapshotTemplate,
		"_snapshot_emitter.go":      templates.ResourceGroupEmitterTemplate,
		"_snapshot_emitter_test.go": templates.ResourceGroupEmitterTestTemplate,
		"_event_loop.go":            templates.ResourceGroupEventLoopTemplate,
		"_event_loop_test.go":       templates.ResourceGroupEventLoopTestTemplate,
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
	for suffix, tmpl := range map[string]*template.Template{
		"_suite_test.go": templates.ProjectTestSuiteTemplate,
	} {
		content, err := generateProjectFile(project, tmpl)
		if err != nil {
			return nil, err
		}
		v = append(v, File{
			Filename: strcase.ToSnake(project.Name) + suffix,
			Content:  content,
		})
	}
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

func generateProjectFile(project *Project, tmpl *template.Template) (string, error) {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, project); err != nil {
		return "", err
	}
	return buf.String(), nil
}
