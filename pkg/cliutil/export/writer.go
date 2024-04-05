package export

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var _ ArchiveWriter = new(localArchiveWriter)

type ArchiveWriter interface {
	Write(ctx context.Context, artifactDir string) error
}

func NewLocalArchiveWriter(targetPath string) ArchiveWriter {
	return &localArchiveWriter{
		targetPath: targetPath,
	}
}

type localArchiveWriter struct {
	// targetDir is the destination directory where the artifact will be written
	targetPath string
}

func (l localArchiveWriter) Write(ctx context.Context, artifactDir string) error {
	if err := os.MkdirAll(filepath.Dir(l.targetPath), os.ModePerm); err != nil {
		return err
	}

	return CreateTarFile(artifactDir, l.targetPath)
}

// CreateTarFile creates a gzipped tar file from srcDir and writes it to outPath.
func CreateTarFile(srcDir, outPath string) error {
	mw, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer mw.Close()
	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return filepath.Walk(srcDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.Mode().IsRegular() {
			return nil
		}
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.Replace(file, srcDir, "", -1), string(filepath.Separator))
		header.Size = fi.Size()
		header.Mode = int64(fi.Mode())
		header.ModTime = fi.ModTime()
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		return nil
	})
}
