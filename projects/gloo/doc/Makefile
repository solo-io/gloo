#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/site
SOURCES := $(shell find . -name "*.md" )

RELEASE := "true"
ifeq ($(TAGGED_VERSION),)
	# TAGGED_VERSION := $(shell git describe --tags)
	# This doesn't work in CI, need to find another way...
	TAGGED_VERSION := vdev
	RELEASE := "false"
endif
VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)

all: site

#----------------------------------------------------------------------------------
# Index
#----------------------------------------------------------------------------------
docs/index.md: ../README.md
	cat ../README.md | sed 's@doc/docs/@@' > $@

docs/cli/glooctl.md: $(shell find ../cli -name "*.go" | grep -v test.go | grep -v '\.\#*')
	go run gen_docs.go

#----------------------------------------------------------------------------------
# Site
#----------------------------------------------------------------------------------

site: docs/index.md docs/cli/glooctl.md
	mkdocs build

# To run this, install firebase CLI and first run firebase login (requires solo account)
# This is run by CI and is not intended to be run manually unless otherwise necessary
.PHONY: deploy-site
deploy-site: site
ifeq ($(RELEASE),"true")
	firebase deploy --only hosting:glooe-site
endif
