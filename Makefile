#######################################################################################################################
#
#    Next raspiBackup version written in go
#
#    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
#
########################################################################################################################

.DEFAULT_GOAL := build
TARGET=raspiBackup
BIN_DIR=bin
MYFILES=$(shell go list ./... | grep -v tools | grep -v -E '/raspiBackupNext$$' |  grep -v -E "discover|model")
export GO111MODULE=on

ifdef DEBUG
	DEBUG=-debug
endif

deps:
	go build -v ./...

test:
	go test ${MYFILES} -v

build: deps test build-local build-raspi

build-local:
	go build -o ${BIN_DIR}/${TARGET} ${TARGET}.go

build-raspi:
	OOS=linux GOARCH=arm GOARM=6 go build -o ${BIN_DIR}/${TARGET}_arm ${TARGET}.go

run:
	go run ${TARGET}.go ${DEBUG}
