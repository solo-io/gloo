package main

import (
	"flag"

	"github.com/solo-io/solo-kit/pkg/utils/log"
)

//go:generate go run ./generate.go -o ../

func main() {
	outdir := flag.String("o", "../", "outdir")
	flag.Parse()
	if err := run(*outdir); err != nil {
		log.Fatalf("generate yaml samples err: %v", err)
	}
}

func run(outdir string) error {
	return nil
}
