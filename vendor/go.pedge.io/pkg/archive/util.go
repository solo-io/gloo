package pkgarchive

import (
	"bufio"
	"os"
	"path/filepath"
)

func parseExcludePatternsFiles(dirPath string, excludePatternsFiles []string) ([]string, error) {
	var excludePatterns []string
	for _, excludePatternsFile := range excludePatternsFiles {
		subExcludePatterns, err := parseExcludePatternsFile(filepath.Join(dirPath, excludePatternsFile))
		if err != nil {
			return nil, err
		}
		if subExcludePatterns != nil && len(subExcludePatterns) > 0 {
			excludePatterns = append(excludePatterns, subExcludePatterns...)
		}
	}
	return excludePatterns, nil
}

func parseExcludePatternsFile(filePath string) (_ []string, retErr error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	var excludes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		excludes = append(excludes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return excludes, nil
}
