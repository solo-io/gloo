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

var _ ArchiveWriter = new(localTarWriter)
var _ ArchiveWriter = new(localDirWriter)

type ArchiveWriter interface {
	Write(ctx context.Context, artifactDir string) error
}

func NewLocalTarWriter(targetPath string) ArchiveWriter {
	return &localTarWriter{
		targetPath: targetPath,
	}
}

type localTarWriter struct {
	// targetPath is the destination path where the tarball will be written
	targetPath string
}

func (l *localTarWriter) Write(ctx context.Context, artifactDir string) error {
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

type localDirWriter struct {
	// targetDir is the destination directory where the artifact will be written
	targetDir string
}

func NewLocalDirWriter(targetDir string) ArchiveWriter {
	return &localDirWriter{
		targetDir: targetDir,
	}
}

func (d *localDirWriter) Write(ctx context.Context, artifactDir string) error {
	if err := os.MkdirAll(d.targetDir, os.ModePerm); err != nil {
		return err
	}

	return d.copyFrom(artifactDir)
}

func (d *localDirWriter) copyFrom(srcDir string) error {
	if err := filepath.Walk(srcDir, func(srcFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		destFile := filepath.Join(d.targetDir, strings.TrimPrefix(srcFile, srcDir))

		// copy
		srcReader, err := os.Open(srcFile)
		if err != nil {
			return err
		}
		defer srcReader.Close()

		if err := os.MkdirAll(filepath.Dir(destFile), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.Create(destFile)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcReader)
		return err

	}); err != nil {
		return err
	}

	return nil
}
