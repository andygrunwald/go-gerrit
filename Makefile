.PHONY: fmt vet check-vendor lint check clean test build
PACKAGES = $(shell go list ./...)
PACKAGE_DIRS = $(shell go list -f '{{ .Dir }}' ./...)

check: test vet lint

test:
	go test -v ./...

vet:
	go vet $(PACKAGES) || (go clean $(PACKAGES); go vet $(PACKAGES))

lint:
	[ -f $(GOPATH)/bin/gometalinter ] || go get -u github.com/alecthomas/gometalinter
	gometalinter --config gometalinter.json ./...

fmt:
	[ -f $(GOPATH)/bin/goimports ] || go get golang.org/x/tools/cmd/goimports
	go fmt $(PACKAGES)
	goimports -w $(PACKAGE_DIRS)
