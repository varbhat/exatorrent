SHELL = /usr/bin/env sh
APP_NAME = exatorrent
PACKAGES ?= ./...
MAIN_SOURCE = exatorrent.go
.DEFAULT_GOAL := help

##help: Display list of commands
.PHONY: help
help: Makefile
	@printf "Options:\n"
	@sed -n 's|^##||p' $<

##web: Build the Web Client
.PHONY: web
web:
	cd internal/web && npm install && npm run build

##web: Build the Web Client for CI
.PHONY: web-ci
web-ci:
	cd internal/web && npm ci && npm run build

##app: Build the Application
.PHONY: app
app:
	env CGO_ENABLED=1 go build -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w"' -o  build/$(APP_NAME) $(MAIN_SOURCE)

##app-no-ui: Build the Application without UI
.PHONY: app-no-ui
app-no-ui:
	env CGO_ENABLED=1 go build -tags noui -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w"' -o  build/$(APP_NAME) $(MAIN_SOURCE)

##app-no-buildflags: Build the Application without any buildflags
.PHONY: app-no-buildflags
app-no-buildflags:
	env CGO_ENABLED=1 go build -o build/$(APP_NAME) $(MAIN_SOURCE)

##app-no-sl: Build the Application without -static build flag
.PHONY: app-no-sl
app-no-sl:
	env CGO_ENABLED=1 go build -trimpath -buildmode=pie -ldflags '-extldflags "-s -w"' -o  build/$(APP_NAME) $(MAIN_SOURCE)

##checksum: Generate sha256 checksums for the builds
.PHONY: checksum
checksum:
	cd build && sha256sum -b * > checksums_sha256.txt

##run: Runs the build
.PHONY: run
run:
	cd build && ./exatorrent*

##docker: Build the Docker image
.PHONY: docker
docker:
	docker build -t "exatorrent" .
