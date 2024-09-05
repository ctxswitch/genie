PWD := $(shell pwd)
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS ?= "-s -w -X build.Version=$(VERSION)"
TMPFILE := $(shell mktemp)
GOPATH := $(shell go env GOPATH)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

LOCALBIN ?= $(shell pwd)/bin

GOLANGCI_LINT_VERSION ?= v1.60.3
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint

deps: $(CONTROLLER_GEN) $(KUSTOMIZE) $(GOLANGCI_LINT) $(ADDLICENSE)

$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

$(GOLANGCI_LINT): $(LOCALBIN)
	@test -s $(GOLANGCI_LINT) || \
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}

.PHONY: verify
verify: vet test lint

.PHONY: test
test:
	go test -race ./... $(GO_VERBOSE) $(GO_COVERPROFILE)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run --fix

.PHONY: vet
vet:
	@go vet ./...

.PHONY: build
build: verify
	go build -trimpath --ldflags $(LDFLAGS) -o genie

.PHONY: testcerts
testcerts:
	@(cd tests/integration/nginx && ./gen-certs.sh)

.PHONY: testnginx
testnginx:
	@docker run --name genie-nginx --rm \
		--mount type=bind,source=${PWD}/tests/integration/nginx/conf,target=/etc/nginx,readonly \
		--entrypoint="" \
		-p 8080:8080 -p 8443:8443 \
		nginx:latest \
		nginx -c /etc/nginx/nginx.conf

.PHONY: clean
clean:
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@find . -name '*.zip' | xargs rm -fv
