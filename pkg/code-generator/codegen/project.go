package codegen

// SOLO-KIT Descriptors from which code can be generated

type ProjectConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	//PackageName string `json:"package_name"`
	//GoPackage   string `json:"go_package"`
}

type Project struct {
	ProjectConfig
	GroupName string `json:"group_name"` // eg. gloo.solo.io

	Resources      []*Resource
	ResourceGroups []*ResourceGroup
}

type Resource struct {
	Name       string
	PluralName string
	ShortName  string

	HasData   bool
	HasStatus bool
	Fields    []*Field

	BelongsToResourceGroups []*ResourceGroup
	BelongsToProject        *Project
}

type Field struct {
	Name     string
	TypeName string
}

type ResourceGroup struct {
	Name             string
	BelongsToProject *Project
	Resources        []*Resource
}
