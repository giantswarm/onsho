PROJECT=onsho
ORGANIZATION=giantswarm

SOURCE := $(shell find . -name '*.go')
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)

.PHONY: all clean run-tests deps bin install

all: deps $(PROJECT)

ci: clean all run-tests

clean:
	rm -rf $(GOPATH) $(PROJECT)

run-tests:
	docker run \
	    --rm \
		-v $(shell pwd):/usr/code \
	    -e GOOS=linux \
	    -e GOARCH=amd64 \
	    -e GOPATH=/usr/code/.gobuild \
	    -w /usr/code \
	    golang:1.5 \
	    go test

# deps
deps: .gobuild
.gobuild:
	mkdir -p $(PROJECT_PATH)
	cd $(PROJECT_PATH) && ln -s ../../../.. $(PROJECT)

	@GOPATH=$(GOPATH) builder go get github.com/spf13/cobra
	@GOPATH=$(GOPATH) builder go get github.com/mitchellh/go-homedir
	@GOPATH=$(GOPATH) builder go get github.com/satori/go.uuid
	@GOPATH=$(GOPATH) builder go get github.com/ryanuber/columnize
	@GOPATH=$(GOPATH) builder go get gopkg.in/yaml.v2

	# Fetch test packages
	@GOPATH=$(GOPATH) builder go get github.com/onsi/gomega
	@GOPATH=$(GOPATH) builder go get github.com/onsi/ginkgo

# build
$(PROJECT): $(SOURCE) VERSION
	@echo Building for $(GOOS)/$(GOARCH)
	docker run \
	    --rm \
	    -v $(shell pwd):/usr/code \
	    -e GOOS=linux \
	    -e GOARCH=amd64 \
	    -e GOPATH=/usr/code/.gobuild \
	    -w /usr/code \
	    golang:1.5 \
	    go build -a -ldflags "-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" -o $(PROJECT)

install: $(PROJECT)
	cp $(PROJECT) /usr/local/bin/
	mkdir -p ~/.giantswarm/onsho/images
	cp ipxe/ipxe.iso ~/.giantswarm/onsho/images

fmt:
	gofmt -l -w .
