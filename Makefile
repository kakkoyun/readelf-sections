VERSION ?= $(shell git describe --exact-match --tags $$(git log -n1 --pretty='%h') 2>/dev/null || echo "$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)")
CONTAINER_IMAGE := ghcr.io/kakkoyun/readelf-sections:$(VERSION)

LDFLAGS="-X main.version=$(VERSION)"

.PHONY: build
build: readelf-sections

.PHONY: install
install: readelf-sections
	go install .

readelf-sections: go.mod main.go
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o $@ main.go

.PHONY: container
container: readelf-sections
	docker build -t $(CONTAINER_IMAGE) .

.PHONY: push-container
push-container:
	docker push $(CONTAINER_IMAGE)
