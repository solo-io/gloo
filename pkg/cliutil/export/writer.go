package export

var _ ArchiveWriter = new(localArchiveWriter)

type ArchiveWriter interface {
	Write(artifactDir string) error
}

func NewLocalArchiveWriter(localDir string) ArchiveWriter {
	return &localArchiveWriter{
		dir: localDir,
	}
}

type localArchiveWriter struct {
	dir string
}

func (l localArchiveWriter) Write(artifactDir string) error {
	//TODO implement me
	panic("implement me")
}
