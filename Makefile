#!/usr/bin/make

SHELL := /bin/bash
binaryName := zonemgr

zonemgr:
	go build -o bin/${binaryName}

zonemgr-a-record-comment-override-plugin:
	mkdir -p examples/bin/comment-override
	go build -o examples/bin/comment-override/zonemgr-a-record-comment-override-plugin examples/comment-override/zonemgr-a-record-comment-override-plugin.go

zonemgr-a-record-not-implemented-plugin:
	mkdir -p examples/bin/not-implemented
	go build -o examples/bin/not-implemented/zonemgr-a-record-not-implemented-plugin examples/not-implemented/zonemgr-a-record-not-implemented-plugin.go

run-test:
	go test ./...

run-test-with-coverage:
	$(eval COV_PKGS=$(shell go list ./... | grep -v examples | grep -v proto | tr '\n' ','))
	go test ./... -cover -coverprofile=coverage.out -coverpkg $(COV_PKGS)
	./exclude-from-coverage.sh
	go tool cover -func coverage.out | grep "^total:"

html-coverage:
	go tool cover -html=coverage.out

format:
	gofmt -l -w -s .

tidy:
	go mod tidy

mocks: proto mocks-gen

mocks-gen:
	mockgen -source=dns/normalizer.go -package dns -self_package "github.com/bcurnow/zonemgr/dns">dns/mock_normalizer.go
	mockgen -source=dns/serial/serial_manager.go -package serial -self_package "github.com/bcurnow/zonemgr/dns/serial">dns/serial/mock_serial_manager.go
	mockgen -source=dns/serial/serial_number_generator.go -package serial -self_package "github.com/bcurnow/zonemgr/dns/serial">dns/serial/mock_serial_number_generator.go
	mockgen -source=plugins/plugin_manager/plugin_manager.go -package plugin_manager -self_package "github.com/bcurnow/zonemgr/plugins/plugin_manager">plugins/plugin_manager/mock_plugin_manager.go
	mockgen -source=plugins/proto/zonemgrplugin_grpc.pb.go -package proto -self_package "github.com/bcurnow/zonemgr/plugins/proto">plugins/proto/mock_zonemgr_plugin_client.go
	mockgen -source=plugins/soa_values_normalizer.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins">plugins/mock_soa_values_normalizer.go
	mockgen -source=plugins/validator.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins">plugins/mock_validator.go
	mockgen -source=plugins/zonemgr_plugin.go -package plugins -self_package "github.com/bcurnow/zonemgr/plugins">plugins/mock_zonemgr_plugin.go
	mockgen -source=utils/filesystem.go  -package utils -self_package "github.com/bcurnow/zonemgr/utils">utils/mock_filesystem.go

proto:
	buf generate

setup: format tidy mocks

build:setup zonemgr

build-all: build zonemgr-a-record-comment-override-plugin zonemgr-a-record-not-implemented-plugin


test: build-all run-test
	
coverage: build-all run-test-with-coverage html-coverage

.PHONY: run-with-plugins

run-with-plugins: zonemgr zonemgr-a-record-comment-override-plugin
	ZONEMGR_PLUGIN_DIR=examples/bin/comment-override ./bin/zonemgr plugins --log-level trace
