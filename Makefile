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
	mkdir -p test
	mockgen -source=plugins/interface.go -package test -self_package "github.com/bcurnow/zonemgr/test">test/plugins_zonemgrplugin.go
	mockgen -source=plugins/manager/interface.go -package test -self_package "github.com/bcurnow/zonemgr/test" >test/manager_pluginmanager.go
	mockgen -source=normalize/interface.go -package test -self_package "github.com/bcurnow/zonemgr/test" >test/normalize_normalizer.go

proto:
	buf generate
