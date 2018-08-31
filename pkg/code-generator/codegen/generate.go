package codegen

type File struct {
	Filename string
	Content  string
}

type Files []File

func GenerateFiles(project *Project) Files {
	files := generateFilesForProject(project)
	for _, res := range project.Resources {
		files = append(files, generateFilesForResource(project, res)...)
	}
	for _, grp := range project.ResourceGroups {
		files = append(files, generateFilesForResourceGroup(project, grp)...)
	}
	return files
}

func generateFilesForResource(project *Project, resource *Resource) Files {
	var v Files

	return v
}

func generateFilesForResourceGroup(project *Project, resource *ResourceGroup) Files {
	var v Files

	return v
}

func generateFilesForProject(project *Project) Files {
	var v Files

	return v
}
