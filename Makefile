VERSION:=$(shell cat VERSION | tr -d '\n')
CONTAINER_IMAGE:=ghcr.io/kakkoyun/readelf-sections:$(VERSION)

LDFLAGS="-X main.version=$(VERSION)"

.PHONY: build
build: readelf-sections

readelf-sections: go.mod main.go
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o $@ main.go

.PHONY: container
container: readelf-sections
	docker build -t $(CONTAINER_IMAGE) .
