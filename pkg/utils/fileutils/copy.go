package fileutils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CopyDir copies all the files from srcDir into targetDir.
func CopyDir(srcDir, targetDir string) error {
	if err := filepath.Walk(srcDir, func(srcFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		destFile := filepath.Join(targetDir, strings.TrimPrefix(srcFile, srcDir))

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
