package pkgarchive

import (
	"io"

	"github.com/docker/docker/pkg/archive"
)

type tarArchiver struct {
	opts ArchiverOptions
}

func newTarArchiver(opts ArchiverOptions) *tarArchiver {
	return &tarArchiver{opts}
}

func (c *tarArchiver) Compress(dirPath string) (io.ReadCloser, error) {
	excludePatterns, err := parseExcludePatternsFiles(dirPath, c.opts.ExcludePatternsFiles)
	if err != nil {
		return nil, err
	}
	return archive.TarWithOptions(
		dirPath,
		&archive.TarOptions{
			IncludeFiles:    []string{"."},
			ExcludePatterns: excludePatterns,
			Compression:     archive.Uncompressed,
			NoLchown:        true,
		},
	)
}

func (c *tarArchiver) Decompress(reader io.Reader, dirPath string) error {
	return archive.Untar(
		reader,
		dirPath,
		&archive.TarOptions{
			NoLchown: true,
		},
	)
}
