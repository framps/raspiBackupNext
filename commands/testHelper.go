package commands

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/framps/raspiBackupNext/tools"
	"github.com/stretchr/testify/assert"
)

// TestCommandType -
type TestCommandType int

const (
	// Blkid -
	Blkid TestCommandType = iota
	// Lsblkid -
	Lsblkid
	// Parted -
	Parted
)

// CommandFromFile -
func CommandFromFile(t TestCommandType, fileName string) (string, error) {

	switch t {
	case Blkid:
		r, e := NewBlkidFromFile(fileName)
		return r.String(), e
	case Lsblkid:
		r, e := NewLsblkFromFile(fileName)
		return r.String(), e
	case Parted:
		r, e := NewPartedFromFile(fileName)
		return r.String(), e
	}
	return "", nil
}

// Command -
func Command(t *testing.T, command TestCommandType, testName string) {

	tools.NewLogger(false)

	files, err := filepath.Glob("testData/" + testName + "_test/*.input")
	assert.NoErrorf(t, err, "Failed to retrieve testdata")

	var tests int
	for _, f := range files {
		t.Logf("Processing %s\n", f)

		parts := strings.Split(f, ".")
		inputFileName := parts[0] + ".input"
		expectedFileName := parts[0] + ".output"

		p, err := CommandFromFile(command, inputFileName)
		assert.NoErrorf(t, err, "Unexpected error")

		expectedResult, err := ioutil.ReadFile(expectedFileName)
		assert.Equal(t, string(expectedResult), p)
		assert.NoError(t, err)
		tests++
	}
	assert.NotZerof(t, tests, "No tests found")
}
