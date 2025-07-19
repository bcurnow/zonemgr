#!/usr/bin/make

SHELL := /bin/bash
currentDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
binaryName := zonemgr

build: clear tidy format build

go-build:
	env GOARCH=amd64 GOOS=linux go build -o bin/${binaryName}

format:
	gofmt -l -w -s .

tidy:
	go mod tidy

clear:
	clear