package types

type Upstream struct {
	// globally unique
	Name string
	// used to determine the upstream module
	Type string
	// Config for the upstream, specific to the type
	Spec
	// Functions specify the functions that live on this upstream
	Functions []Function
}

type Function struct {
	Name string
	// Config for the function, specific to the upstream type
	Spec
}
