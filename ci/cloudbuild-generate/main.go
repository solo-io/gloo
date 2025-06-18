package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

type TemplateData struct {
	Version string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <cloudbuild_version>")
		fmt.Println("Example: go run main.go 0.13.0")
		os.Exit(1)
	}

	version := os.Args[1]
	data := TemplateData{Version: version}

	templates := []string{
		"publish-artifacts.yaml.tmpl",
		"run-tests.yaml.tmpl",
	}

	for _, tmplName := range templates {
		tmplContent, err := templateFS.ReadFile("templates/" + tmplName)
		if err != nil {
			fmt.Printf("Error reading template %s: %v\n", tmplName, err)
			os.Exit(1)
		}

		tmpl, err := template.New(tmplName).Parse(string(tmplContent))
		if err != nil {
			fmt.Printf("Error parsing template %s: %v\n", tmplName, err)
			os.Exit(1)
		}

		outputDir := "../../ci/cloudbuild"
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", outputDir, err)
			os.Exit(1)
		}

		outputFile := filepath.Join(outputDir, tmplName[:len(tmplName)-5])
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", outputFile, err)
			os.Exit(1)
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
			fmt.Printf("Error executing template %s: %v\n", tmplName, err)
			os.Exit(1)
		}

		fmt.Printf("Generated %s\n", outputFile)
	}

	fmt.Printf("Successfully generated cloudbuild files with version %s\n", version)
}
