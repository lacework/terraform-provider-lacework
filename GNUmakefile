# Tooling versions
GOLANGCILINTVERSION?=1.23.8
GOIMPORTSVERSION?=v0.1.2
GOXVERSION?=v1.0.1

TEST?=$$(go list ./... |grep -v 'vendor' | grep -v 'integration')
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=lacework
DIR=~/.terraform.d/plugins
GO_CLIENT_VERSION=master
COVERAGEOUT?=coverage.out
GOFLAGS=-mod=vendor
CGO_ENABLED?=0
PACKAGENAME?=terraform-provider-lacework
VERSION=$(shell cat VERSION)
export GOFLAGS CGO_ENABLED

default: build

ci: lint test fmtcheck imports-check

prepare: install-tools go-vendor

release: build-cross-platform
	scripts/release.sh prepare

deps:
ifdef UPDATE_DEP
	@go get -u "$(UPDATE_DEP)"
endif
	@go mod vendor

alldeps:
	@go get -u
	@go mod vendor

go-vendor:
	go mod tidy
	go mod vendor
	go mod verify

build: fmtcheck
	go install

build-cross-platform:
	gox -output="bin/$(PACKAGENAME)_$(VERSION)_{{.OS}}_{{.Arch}}" \
            -os="linux windows freebsd" \
            -osarch="darwin/amd64 darwin/arm64 linux/arm linux/arm64 freebsd/arm freebsd/arm64" \
            -arch="amd64 386" \
            github.com/lacework/$(PACKAGENAME)

install: fmtcheck
	mkdir -vp $(DIR)
	go build -o $(DIR)/$(PACKAGENAME)

uninstall:
	@rm -vf $(DIR)/$(PACKAGENAME)

integration-test:
	go test ./integration -v

test: fmtcheck
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

lint:
	golangci-lint run

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test

install-tools:
ifeq (, $(shell which golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCILINTVERSION)
endif
ifeq (, $(shell which goimports))
	go get golang.org/x/tools/cmd/goimports@$(GOIMPORTSVERSION)
endif
ifeq (, $(shell which gox))
	go get github.com/mitchellh/gox@$(GOXVERSION)
endif
