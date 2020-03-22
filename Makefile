ARCH := darwin/amd64 linux/386 linux/amd64 windows/386 windows/amd64
VERSION := $(shell git describe --always --dirty)
LDFLAGS := -ldflags "-X main.version=${VERSION}"
GOFILES := $(shell go list -f '{{range .GoFiles}}{{.}} {{end}}' ./...)

remux: ${GOFILES}
	go build ${LDFLAGS}
test:
	go test

release: ${GOFILES}
	gox -output 'build/{{.Dir}}.{{.OS}}.{{.Arch}}' -osarch "${ARCH}" ${LDFLAGS}

fmt: ${GOFILES}
	go fmt ./...
tidy:
	go mod tidy

clean:
	rm -rf build remux

.PHONY: test release fmt tidy clean
