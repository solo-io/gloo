package nsselect

// All the cli-relevant resources keyed by namespace
type NsResourceMap map[string]*NsResource

// NsResource contains lists of the resources needed by the cli associated* with given namespace.
// *the association is by the namespace in which the CRD is installed, unless otherwise noted.
type NsResource struct {
	// keyed by namespace containing the CRD
	Secrets   []string
	Upstreams []string
}
