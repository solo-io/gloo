#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/site
SOURCES := $(shell find . -name "*.md" )
# TODO: use the above instead
VERSION ?= dev

all: site

#----------------------------------------------------------------------------------
# Index
#----------------------------------------------------------------------------------
docs/index.md: ../README.md
	cat ../README.md | sed 's@doc/docs/@@' > $@

#----------------------------------------------------------------------------------
# Site
#----------------------------------------------------------------------------------

.PHONY: generated-docs
generated-docs: docs/v1

docs/v1: $(shell find ../api)
	go run $(GOPATH)/src/github.com/solo-io/gloo/cmd/solo-kit-gen/main.go -r ../../gateway -i ../api
	go run $(GOPATH)/src/github.com/solo-io/gloo/cmd/solo-kit-gen/main.go -r ../
	# TODO ilackarms: make sure projects aren't creating docs for other projects
	# which gateway does for gloo here
	rm -rf docs/v1/github.com/

site: docs/index.md generated-docs
	mkdocs build

.PHONY: docker-site
docker-site: site
	docker build -t soloio/gloo-docs:$(VERSION) .

.PHONY: docker-site-push
docker-site-push: docker-site
	docker push soloio/gloo-docs:$(VERSION)
