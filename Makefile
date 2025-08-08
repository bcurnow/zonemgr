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
	ZONEMGR_PLUGINS=examples/bin/ ./bin/zonemgr generate --inputFile examples/zones.yaml -l trace