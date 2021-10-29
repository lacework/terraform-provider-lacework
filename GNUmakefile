# Tooling versions
GOLANGCILINTVERSION?=1.23.8
GOIMPORTSVERSION?=v0.1.2
GOXVERSION?=v1.0.1
GOTESTSUMVERSION?=v1.6.4

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
BINARY_PATH="registry.terraform.io/lacework/lacework/99.0.0/$$(go env GOOS)_$$(go env GOARCH)/terraform-provider-lacework_v99.0.0"
export GOFLAGS CGO_ENABLED

.PHONY: help
help:
	@echo "-------------------------------------------------------------------"
	@echo "Lacework terraform-provider-lacework Makefile helper:"
	@echo ""
	@grep -Fh "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/\\$$//' | sed -E 's/^([^:]*):.*##(.*)/ \1 -\2/'
	@echo "-------------------------------------------------------------------"

default: build

.PHONY: ci
ci: lint test fmtcheck imports-check ## *CI ONLY* Runs tests on CI pipeline

.PHONY: prepare
prepare: install-tools go-vendor ## Initialize the go environment

.PHONY: release
release: build-cross-platform ## *CI ONLY* Prepares a release of the Terraform provider
	scripts/release.sh prepare

.PHONY: deps
deps: ## Update a single dependency by providing the UPDATE_DEP environment variable
ifdef UPDATE_DEP
	@go get -u "$(UPDATE_DEP)"
endif
	@go mod vendor

.PHONY: alldeps
alldeps: ## Update all dependencies
	@go get -u
	@go mod vendor

PHONY: go-vendor
go-vendor: ## Runs go mod tidy, vendor and verify to cleanup, copy and verify dependencies
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: build
build: fmtcheck ## Runs fmtcheck and go install
	go install

.PHONY: build-cross-platform
build-cross-platform: ## Compiles the Terraform-Provider-Lacework for all supported platforms
	gox -output="bin/$(PACKAGENAME)_$(VERSION)_{{.OS}}_{{.Arch}}" \
            -os="linux windows freebsd" \
            -osarch="darwin/amd64 darwin/arm64 linux/arm linux/arm64 freebsd/arm freebsd/arm64" \
            -arch="amd64 386" \
            github.com/lacework/$(PACKAGENAME)

.PHONY: install
install: write-terraform-rc fmtcheck ## Updates the terraformrc to point to the BINARY_PATH. Installs the provider to the BINARY_PATH
	mkdir -vp $(DIR)
	go build -o $(DIR)/$(BINARY_PATH)

.PHONY: uninstall
uninstall: ## Removes installed provider package from BINARY_PATH
	@rm -vf $(DIR)/$(PACKAGENAME)

.PHONY: integration-test
integration-test: clean-test install ## Runs clean-test and install, then runs all integration tests
	gotestsum -f testname -- -v ./integration

.PHONY: test
test: fmtcheck ## Runs fmtcheck then runs all unit tests
	gotestsum -f testname -- -v -cover -coverprofile=$(COVERAGEOUT) $(TEST)

.PHONY: lint
lint: ## Runs go linter
	golangci-lint run

.PHONY: testacc
testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

.PHONY: vet
vet: ## Runs go vet
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: fmt
fmt: ## Runs go formatter
	gofmt -w $(GOFMT_FILES)

.PHONY: fmtcheck
fmtcheck: ## Runs formatting check
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: errcheck
errcheck: ## Runs error checking
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: test-compile
test-compile: ## Compile tests
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: website
website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: website-test
website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: install-tools
install-tools: ## Install go indirect dependencies
ifeq (, $(shell which golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCILINTVERSION)
endif
ifeq (, $(shell which goimports))
	GOFLAGS=-mod=readonly go install golang.org/x/tools/cmd/goimports@$(GOIMPORTSVERSION)
endif
ifeq (, $(shell which gox))
	GOFLAGS=-mod=readonly go install github.com/mitchellh/gox@$(GOXVERSION)
endif
ifeq (, $(shell which gotestsum))
	GOFLAGS=-mod=readonly go install gotest.tools/gotestsum@$(GOTESTSUMVERSION)
endif

.PHONY: write-terraform-rc
write-terraform-rc: ## Write to terraformrc file to mirror lacework/lacework to BINARY_PATH
	scripts/mirror-provider.sh

.PHONY: clean-test
clean-test: ## Find and remove any .terraform directories or tfstate files
	find . -name ".terraform*" -type f -exec rm -rf {} \;
	find . -name "terraform.tfstate*" -type f -exec rm -rf {} \;
