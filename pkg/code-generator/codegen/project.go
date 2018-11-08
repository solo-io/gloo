package codegen

// SOLO-KIT Descriptors from which code can be generated

type ProjectConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Project struct {
	ProjectConfig
	GroupName string

	Resources      []*Resource
	ResourceGroups []*ResourceGroup

	XDSResources []*XDSResource
}

type Resource struct {
	Name       string
	PluralName string
	ShortName  string
	GroupName  string // eg. gloo.solo.io
	// ImportPrefix will equal GroupName+"." if the resource does not belong to the project
	// else it will be empty string. used in event loop files
	ImportPrefix string
	// empty unless resource is external
	// format "github.com/solo-io/solo-kit/foo/bar"
	GoPackage string

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
	Imports   string // if this resource group contains any imports from other projects
	Project   *Project
	Resources []*Resource
}

type XDSResource struct {
	Name         string
	MessageType  string
	NameField    string
	NoReferences bool

	Project   *Project
	GroupName string // eg. gloo.solo.io
	Package   string // proto package for the message
}
