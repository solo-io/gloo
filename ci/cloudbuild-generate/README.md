# Cloudbuild Generator

This tool generates cloudbuild YAML files with the specified version.

## Usage

```bash
go run main.go <cloudbuild_version>
```

Example:
```bash
go run main.go 0.13.0
```

## Directory Structure

- `templates/`: Contains the template files for cloudbuild YAML files
- `main.go`: The main script that generates the files
- `README.md`: This file

## Templates

- `publish-artifacts.yaml.tmpl`: Template for publish-artifacts.yaml
- `run-tests.yaml.tmpl`: Template for run-tests.yaml 