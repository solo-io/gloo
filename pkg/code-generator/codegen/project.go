package codegen

// SOLO-KIT Descriptors from which code can be generated

type ProjectConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	PackageName string `json:"package_name"`
	// GoPackage   string `json:"go_package"`
}

type Project struct {
	ProjectConfig
	GroupName string `json:"group_name"` // eg. gloo.solo.io

	Resources      []*Resource
	ResourceGroups []*ResourceGroup

	XDSResources []*XDSResource
}

type Resource struct {
	Name       string
	PluralName string
	ShortName  string

	HasStatus bool
	Fields    []*Field

	ResourceGroups []*ResourceGroup
	Project        *Project
}

type Field struct {
	Name     string
	TypeName string
}

type ResourceGroup struct {
	Name      string // eg. api.gloo.solo.io
	GoName    string // will be Api
	Project   *Project
	Resources []*Resource
}

type XDSResource struct {
	Name         string
	MessageType  string
	NameField    string
	NoReferences bool

	Project *Project
}
