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

##app: Build the Application
.PHONY: app
app:
	env CGO_ENABLED=1 go build -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w"' -o  build/$(APP_NAME) $(MAIN_SOURCE)

##app-linux-amd64: Build the Application for linux (amd64)
.PHONY: app-linux-amd64
app-linux-amd64:
	env CGO_ENABLED=1 GOOS="linux" GOARCH="amd64" CC="x86_64-linux-musl-gcc" CXX="x86_64-linux-musl-g++" go build -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w"' -o  build/$(APP_NAME)-linux-amd64 $(MAIN_SOURCE)

##app-linux-arm64: Build the Application for linux (arm64)
.PHONY: app-linux-arm64
app-linux-arm64:
	env CGO_ENABLED=1 GOOS="linux" GOARCH="arm64" CC="aarch64-linux-musl-gcc" CXX="aarch64-linux-musl-g++" go build -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w"' -o  build/$(APP_NAME)-linux-arm64 $(MAIN_SOURCE)

##app-windows-amd64: Build the Application for windows (amd64)
.PHONY: app-windows-amd64
app-windows-amd64:
	env CGO_ENABLED=1 GOOS="windows" GOARCH="amd64" CC="x86_64-w64-mingw32-gcc" CXX="x86_64-w64-mingw32-g++" go build -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w"' -o  build/$(APP_NAME)-windows-amd64 $(MAIN_SOURCE)

##checksum: Generate sha256 checksums for the builds
.PHONY: checksum
checksum:
	cd build && sha256sum -b * > checksums_sha256.txt

##run: Runs the build
.PHONY: run
run:
	cd build && ./exatorrent*
