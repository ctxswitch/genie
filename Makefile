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

TARGETDIR ?= $(PWD)/dist
GENIE_BIN ?= genie
GENIE_RELEASE_TARGET ?= $(TARGETDIR)/genie-$(SYSTEM)-$(ARCH).tar.gz

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

$(TARGETDIR):
	mkdir -p $(TARGETDIR)

.PHONY: build-release
build-release: $(TARGETDIR) $(GENIE_RELEASE_TARGET)

$(GENIE_RELEASE_TARGET):
	GOOS=$(SYSTEM) GOARCH=$(ARCH) CGO_ENABLED=0 go build -trimpath --ldflags "-s -w -X ctx.sh/genie/pkg/build.Version=$(VERSION)" -o $(TARGETDIR)/$(GENIE_BIN) main.go && \
		tar -C $(TARGETDIR) -zcvf $@ $(GENIE_BIN) && \
		rm -f $(TARGETDIR)/$(GENIE_BIN)

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
