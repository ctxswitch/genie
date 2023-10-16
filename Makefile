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

deps:
	@echo "Checking dependencies"
	@(echo "Installing golangci-lint" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.52.2)

verify: vet test lint

lint:
	@GO111MODULE=on ${GOPATH}/bin/golangci-lint cache clean
	@GO111MODULE=on ${GOPATH}/bin/golangci-lint run --timeout=5m --config ./.golangci.yml

vet:
	@go vet ./...

test:
	@go test -race ./...

build: verify
	go build -trimpath --ldflags $(LDFLAGS) -o genie

# Temporary.  Will set up automated builds later.
publish:
	@docker build . -t strataviz/genie:0.1.1
	@docker tag strataviz/genie:0.1.1 strataviz/genie:latest
	@docker push strataviz/genie --all-tags

testcerts:
	@(cd tests/integration/nginx && ./gen-certs.sh)

testnginx:
	@docker run --name genie-nginx --rm \
		--mount type=bind,source=${PWD}/tests/integration/nginx/conf,target=/etc/nginx,readonly \
		--entrypoint="" \
		-p 8080:8080 -p 8443:8443 \
		nginx:latest \
		nginx -c /etc/nginx/nginx.conf

clean:
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@find . -name '*.zip' | xargs rm -fv
