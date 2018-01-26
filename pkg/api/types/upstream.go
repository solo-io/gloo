package types

type Upstream struct {
	// globally unique
	Name string
	// used to determine the upstream module
	Type string
	// Config for the upstream, specific to the type
	Config Spec
}
