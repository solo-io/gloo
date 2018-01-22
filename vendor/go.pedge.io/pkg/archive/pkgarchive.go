package pkgarchive // import "go.pedge.io/pkg/archive"

import "io"

// Archiver compresses and decompresses a directory.
type Archiver interface {
	Compress(dirPath string) (io.ReadCloser, error)
	Decompress(reader io.Reader, dirPath string) error
}

// ArchiverOptions are options to the construction of a Archiver.
type ArchiverOptions struct {
	// If not set, no files are excluded.
	ExcludePatternsFiles []string
}

// NewTarArchiver returns a new Archiver for tar.
func NewTarArchiver(opts ArchiverOptions) Archiver {
	return newTarArchiver(opts)
}
