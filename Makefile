#!/usr/bin/make

SHELL := /bin/bash
currentDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
binaryName := zonemgr

build: tidy format zonemgr

build-all: build zonemgr-a-record-comment-override-plugin

zonemgr:
	go build -o bin/${binaryName}

format:
	gofmt -l -w -s .

tidy:
	go mod tidy

zonemgr-a-record-comment-override-plugin:
	mkdir -p examples/bin
	go build -o examples/bin/zonemgr-a-record-comment-override-plugin examples/zonemgr-a-record-comment-override-plugin.go

.PHONY: run-with-plugins

run-with-plugins: zonemgr zonemgr-a-record-comment-override-plugin
	ZONEMGR_PLUGINS=examples/bin/ ./bin/zonemgr plugins

mocks:
	mockgen -source=dns/normalizer.go -package dns -self_package "github.com/bcurnow/zonemgr/dns" >dns/mock_normalizer.go
	mockgen -source=plugin_manager/plugin_manager.go -package plugin_manager -self_package "github.com/bcurnow/zonemgr/plugin_manager" >plugin_manager/mock_plugin_manager.go
	mockgen -source=plugins/soa_values_normalizer.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins">plugins/mock_soa_values_normalizer.go
	mockgen -source=plugins/validator.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins">plugins/mock_validator.go
	mockgen -source=plugins/zonemgr_plugin.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins">plugins/mock_zonemgr_plugin.go
	mockgen -source=utils/serial_index_manager.go -package utils -self_package "github.com/bcurnow/zonemgr/utils">utils/mock_serial_index_manager.go

proto:
	buf generate

test:
	go test -cover -coverprofile=coverage.out ./...

html-coverage:
	go tool cover -html=coverage.out

coverage: test html-coverage
