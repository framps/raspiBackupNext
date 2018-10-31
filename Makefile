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
MYFILES=$(shell go list ./... | grep -v /vendor/ | grep -v tools | grep -v -E '/raspiBackupNext$$' |  grep -v -E "discover|model")

ifdef DEBUG
	DEBUG=-debug
endif

setup:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

update:
	dep ensure

test:
	go test ${MYFILES} -v

build: setup update test-verbose build-local build-raspi

build-local:
	go build -o ${BIN_DIR}/${TARGET} ${TARGET}.go

build-raspi:
	OOS=linux GOARCH=arm GOARM=6 go build -o ${BIN_DIR}/${TARGET}_arm ${TARGET}.go

run:
	go run ${TARGET}.go ${DEBUG}
