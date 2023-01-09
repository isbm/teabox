.DEFAULT_GOAL := build

GO_BUILD = go build -ldflags="-s -w"
PREFIX ?= /usr
BINDIR ?= ${PREFIX}/bin

fmt:
	go fmt ./...
.PHONY:fmt

vet: fmt
	go vet ./...
.PHONY:vet

build: vet \
	clean \
	teabox

.PHONY:build

teabox:
	cd cmd && $(GO_BUILD) -o teabox

clean:
	rm -f cmd/teabox

vendor:
	go mod tidy
	go mod vendor

package:
	go build -a -mod=vendor -buildmode=pie -ldflags="-s -w" -o teabox ./cmd/*go
