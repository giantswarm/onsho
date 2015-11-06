PROJECT=moa
ORGANIZATION=giantswarm

SOURCE := $(shell find . -name '*.go')
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif
ifndef GOARCH
	GOARCH := $(shell go env GOARCH)
endif

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

	docker run \
	    --rm \
		-v $(shell pwd):/usr/code \
	    -e GOOS=linux \
	    -e GOARCH=amd64 \
	    -e GOPATH=/usr/code/.gobuild \
	    -w /usr/code \
	    golang:1.5 \
	    go get github.com/$(ORGANIZATION)/$(PROJECT)

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

fmt:
	gofmt -l -w .
