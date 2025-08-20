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
	mockgen -source=plugins/zonemgr_plugin.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins">plugins/mock_zonemgr_plugin.go
	mockgen -source=plugins/validator.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins">plugins/mock_validator.go
	mockgen -source=plugins/plugin_manager.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins" >plugins/mock_plugin_manager.go
	mockgen -source=dns/normalizer.go -package dns -self_package "github.com/bcurnow/zonemgr/dns" >dns/mock_normalizer.go

proto:
	buf generate
