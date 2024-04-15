package export

import (
	"context"
	"github.com/solo-io/gloo/pkg/utils/fileutils"
	"os"
	"path/filepath"
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

	return fileutils.CreateTarFile(artifactDir, l.targetPath)
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

	return fileutils.CopyDir(artifactDir, d.targetDir)
}
