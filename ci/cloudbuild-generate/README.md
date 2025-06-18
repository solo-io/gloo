# Cloudbuild Generator

This tool generates cloudbuild YAML files with the specified version.

## Usage

### Direct Go Run

```bash
go run main.go <cloudbuild_version>
```

Example:
```bash
go run main.go 0.13.0
```

### Using the Makefile Target

You can also use the provided Makefile target to generate the cloudbuild YAML files. This is the recommended way if you are working from the root of the repository.

```bash
make update-cloudbuild-version CLOUDBUILD_VERSION=<cloudbuild_version>
```

Example:
```bash
make update-cloudbuild-version CLOUDBUILD_VERSION=0.13.0
```

> **Note:** The `CLOUDBUILD_VERSION` environment variable is required. The Makefile target will run the generator with the specified version.

## Directory Structure

- `templates/`: Contains the template files for cloudbuild YAML files
- `main.go`: The main script that generates the files
- `README.md`: This file

## Templates

- `publish-artifacts.yaml.tmpl`: Template for publish-artifacts.yaml
- `run-tests.yaml.tmpl`: Template for run-tests.yaml 