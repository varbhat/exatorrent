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

##app-darwin-amd64: Build the Application for MacOS (amd64)
.PHONY: app-darwin-amd64
app-darwin-amd64:
	env CGO_ENABLED=1 GOOS="darwin" GOARCH="amd64" CC="o64-clang" CXX="o64-clang++" go build -trimpath -buildmode=pie -ldflags '-extldflags "-s -w"' -o  build/$(APP_NAME)-darwin-amd64 $(MAIN_SOURCE)

##app-darwin-arm64: Build the Application for MacOS (arm64)
.PHONY: app-darwin-arm64
app-darwin-arm64:
	env CGO_ENABLED=1 GOOS="darwin" GOARCH="arm64" CC="oa64-clang" CXX="oa64-clang++" go build -trimpath -buildmode=pie -ldflags '-extldflags "-s -w"' -o  build/$(APP_NAME)-darwin-arm64 $(MAIN_SOURCE)

##app-win-amd64: Build the Application for Windows (amd64)
.PHONY: app-win-amd64
app-win-amd64:
	env CGO_ENABLED=1 GOOS="windows" GOARCH="amd64" CC="x86_64-w64-mingw32-gcc" CXX="x86_64-w64-mingw32-g++" go build -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w"' -o  build/$(APP_NAME)-win-amd64.exe $(MAIN_SOURCE)

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
