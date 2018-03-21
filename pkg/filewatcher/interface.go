package filewatcher

// map [ref] : contents
type Files map[string][]byte

// Interface is responsible for watching artifacts referenced by a config
type Interface interface {
	TrackFiles(fileRefs []string)

	// artifacts are pushed here whenever they are read
	Files() <-chan Files

	// should show valid if the most recent update passed, otherwise a useful error
	Error() <-chan error
}
