package main

import (
	"fmt"
	"os"

	"go.pedge.io/pkg/file"
)

func main() {
	if err := do(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func do() error {
	for _, filePath := range os.Args[1:] {
		if err := pkgfile.StripPackageCommentsForFile(filePath); err != nil {
			return err
		}
	}
	return nil
}
