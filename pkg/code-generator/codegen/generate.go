package codegen

type file struct {
	filename string
	content  string
}

type files []file

func generateFiles(project *Project) files {
	files := generateFilesForProject(project)
	for _, res := range project.Resources {
		files = append(files, generateFilesForResource(project, res)...)
	}
	for _, grp := range project.ResourceGroups {
		files = append(files, generateFilesForResourceGroup(project, grp)...)
	}
	return files
}

func generateFilesForResource(project *Project, resource *Resource) files {
	var v files

	return v
}

func generateFilesForResourceGroup(project *Project, resource *ResourceGroup) files {
	var v files

	return v
}

func generateFilesForProject(project *Project) files {
	var v files

	return v
}
