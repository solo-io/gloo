package artifactwatcher

// map [ref] : map[filename] : contents
type Artifacts map[string]map[string][]byte

// Interface is responsible for watching artifacts referenced by a config
type Interface interface {
	TrackArtifacts(artifactRefs []string)

	// artifacts are pushed here whenever they are read
	Artifacts() <-chan Artifacts

	// should show valid if the most recent update passed, otherwise a useful error
	Error() <-chan error
}
