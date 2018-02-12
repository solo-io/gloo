package main

import (
	"log"

	//register plugins
	_ "github.com/solo-io/glue/internal/install"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
