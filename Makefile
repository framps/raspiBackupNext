.DEFAULT_GOAL := build
TARGET=raspiBackup
BIN_DIR=bin
MYFILES=$(shell go list ./... | grep -v /vendor/ | grep -v tools | grep -v -E '/go[^/]')
ifdef DEBUG
	DEBUG=-debug
endif

setup: build
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

update:
	dep ensure

test:
	go test ${MYFILES}

test-verbose:
		go test ${MYFILES} -v

build: test build-local build-raspi

build-local:
	go build -o ${BIN_DIR}/${TARGET} ${TARGET}.go

build-raspi:
	OOS=linux GOARCH=arm GOARM=6 go build -o ${BIN_DIR}/${TARGET}_arm ${TARGET}.go

run:
	go run ${TARGET}.go ${DEBUG}
