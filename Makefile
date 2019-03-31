PROJECT_NAME := "github.com/tupass/tupass-backend"
PKG := "$(PROJECT_NAME)"
PKG_TESTING := "github.com/tupass/tupass-backend/testing"
PKG_LIST := $(shell go list ${PKG}/... | grep -v github.com/tupass/tupass-backend/testing)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)
PATH := $(GOPATH)/bin:$(PATH)

.PHONY: all gosec lint test test-api race dep dep-gosec dep-golint build build-local build-deb build-prod build-lib build-pam run-dev run-prod clean

all: run-dev

gosec: dep dep-gosec ## Check all packages with gosec
	@gosec -include=G102,G104 $(PKG)/...

lint: dep-golint ## Lint the files
	@golint -set_exit_status $(PKG)/...

test: dep ## Run unittests
	@go test -v -tags unit ${PKG_TESTING}

test-api: build ## Test that API respondes and returns a correct result
	./testing/test-api.sh

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

dep-gosec: ## Get dependencies required for gosec
	@go get github.com/securego/gosec/cmd/gosec/...

dep-golint: ## Get dependencies required for golint
	@go get -u golang.org/x/lint/golint

dep: ## Get dependencies required for Go implementation
	@go get github.com/GeertJohan/go.rice
	@go get github.com/GeertJohan/go.rice/rice
	@go get -v -d ./...

build: dep ## Build the binary file
	cd api && rice embed-go
	@go build -ldflags "-s -w" -o tupass-backend -i -v $(PKG)

build-local: dep ## Build bundled version of backend with frontend (and password list)
	cd api && rice embed-go
	cd web && ./buildFrontend.sh && rice embed-go
	@go build -ldflags="-X github.com/tupass/tupass-backend/web.localBuild=true -s -w" -i -v -o tupass $(PKG)
	GOOS=windows GOARCH=386 go build -ldflags="-X github.com/tupass/tupass-backend/web.localBuild=true -s -w" -o tupass.exe main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X github.com/tupass/tupass-backend/web.localBuild=true -s -w" -o tupass-darwin main.go

build-deb: ## Build .deb package for debian systems
	./build-deb.sh

## Build /usr/lib/libtupass.so (with bundled password list) and /usr/include/libtupass.h
build-lib: dep
	cd pam && ./buildLib.sh

## Build /lib/secruity/pam_tupass.so and copy PAM module config to /usr/share/pam-configs/tupass 
build-pam: build-lib
	cd pam && ./buildPam.sh

run-dev: build ## Run the backend for dev server
	APP_ENV=dev ./tupass-backend

run-prod: build ## Run the backend for staging/production server
	APP_ENV=prod ./tupass-backend

clean: ## Remove previous build
	@rm -rf api/rice-box.go pam/libtupass-test pam/libtupass.so pam/libtupass.h pam/pam_tupass.o tupass-1.0/usr/bin/tupass web/rice-box.go web/frontend web/tupass-frontend tupass-*.deb tupass tupass-backend tupass.exe

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
