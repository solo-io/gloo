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

site: docs/index.md
	mkdocs build

.PHONY: docker-site
docker-site: site
	docker build -t soloio/gloo-ee-docs:$(VERSION) .

.PHONY: docker-site-push
docker-site-push: docker-site
	docker push soloio/gloo-ee-docs:$(VERSION)

