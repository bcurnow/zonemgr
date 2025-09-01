#!/usr/bin/make

SHELL := /bin/bash
binaryName := zonemgr
mocks := dns/normalizer.go plugin_manager/plugin_manager.go plugins/soa_values_normalizer.go plugins/validator.go plugins/zonemgr_plugin.go utils/serial_index_manager.go

zonemgr:
	go build -o bin/${binaryName}

zonemgr-a-record-comment-override-plugin:
	mkdir -p examples/bin
	go build -o examples/bin/zonemgr-a-record-comment-override-plugin examples/zonemgr-a-record-comment-override-plugin.go

run-test:
	go test ./...

run-test-with-coverage:
	go test -cover -coverprofile=coverage.out -coverpkg github.com/bcurnow/zonemgr/cmd,github.com/bcurnow/zonemgr/dns,github.com/bcurnow/zonemgr/internal/plugins/builtin,github.com/bcurnow/zonemgr/models,github.com/bcurnow/zonemgr/plugin_manager,github.com/bcurnow/zonemgr/plugins,github.com/bcurnow/zonemgr/plugins/grpc,github.com/bcurnow/zonemgr/utils ./...

html-coverage:
	go tool cover -html=coverage.out

format:
	gofmt -l -w -s .

tidy:
	go mod tidy

mocks:
	@mkdir -p internal/mocks
	@$(foreach mock, $(mocks), mockgen -source=$(mock) -package mocks -self_package "github.com/bcurnow/zonemgr/internal/mocks" >internal/mocks/`basename $(mock)`;)

proto:
	buf generate

setup: format tidy mocks proto

build:setup zonemgr

build-all: build zonemgr-a-record-comment-override-plugin

test: build-all run-test
	
coverage: build-all run-test-with-coverage html-coverage

.PHONY: run-with-plugins

run-with-plugins: zonemgr zonemgr-a-record-comment-override-plugin
	ZONEMGR_PLUGIN_DIR=examples/bin/ ./bin/zonemgr plugins --log-level trace
