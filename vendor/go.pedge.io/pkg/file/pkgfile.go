package pkgfile // import "go.pedge.io/pkg/file"

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// EditFile edits the file using f and saves it.
func EditFile(filePath string, f func(io.Reader, io.Writer) error) (retErr error) {
	tempFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("%s.%d", filepath.Base(filePath), os.Getpid()))
	if err := os.Rename(filePath, tempFilePath); err != nil {
		return err
	}
	var tempFile *os.File
	defer func() {
		if tempFile != nil {
			if err := tempFile.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}
		if err := os.Remove(tempFilePath); err != nil && retErr == nil {
			retErr = err
		}
	}()
	var err error
	tempFile, err = os.Open(tempFilePath)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return f(tempFile, file)
}

// StripPackageCommentsForFile strips the package comments for a file.
//
// TODO(pedge): not real verification that we are in package comments,
// just takes the first block comment in the file and eliminates it.
func StripPackageCommentsForFile(filePath string) error {
	return EditFile(filePath, editStripPackageCommentsForFile)
}

func editStripPackageCommentsForFile(reader io.Reader, writer io.Writer) error {
	bufioReader := bufio.NewReader(reader)
	inPackageComments := false
	for line, err := bufioReader.ReadString('\n'); err != io.EOF; line, err = bufioReader.ReadString('\n') {
		if err != nil {
			return err
		}
		if strings.HasPrefix(line, "/*") {
			inPackageComments = true
			continue
		}
		if inPackageComments && strings.HasPrefix(line, "*/") {
			if _, err := bufioReader.WriteTo(writer); err != nil {
				return err
			}
			return nil
		}
		if !inPackageComments {
			if _, err := writer.Write([]byte(line)); err != nil {
				return err
			}
		}
	}
	return nil
}
