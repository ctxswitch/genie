PWD := $(shell pwd)
LDFLAGS ?= "-s -w -X main.Version=$(VERSION)"
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
	@CGO_ENABLED=0 GOOS=linux go build -trimpath --ldflags $(LDFLAGS) -o dynamo

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

# DELETE ME
.PHONY: run
run:
	$(eval POD := $(shell kubectl get pods -n strataviz -l name=strataviz-analytics -o=custom-columns=:metadata.name --no-headers))
	kubectl exec -n strataviz -it pod/$(POD) -- bash -c "go run main.go -s kafka.kind"
