SHELL := /bin/bash
BASEDIR = $(shell pwd)

export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOPRIVATE=*.gitlab.com
export GOSUMDB=off

fmt:
	gofmt -w .
mod:
	go mod tidy
lint:
	golangci-lint run
.PHONY: test
test: mod
	go test -gcflags=-l -coverpkg=./... -coverprofile=coverage.data ./...
.PHONY: mysql
mysql:
	sh scripts/mysql.sh
help:
	@echo "fmt - format the source code"
	@echo "mod - go mod tidy"
	@echo "lint - run golangci-lint"
	@echo "test - unit test"
	@echo "mysql - launch a docker mysql"