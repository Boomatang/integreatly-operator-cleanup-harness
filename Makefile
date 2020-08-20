DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
OUT_FILE := "$(DIR)integreatly-operator-test-harness"

build:
	go mod vendor
	CGO_ENABLED=0 go test -v -c
