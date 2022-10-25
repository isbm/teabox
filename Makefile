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
