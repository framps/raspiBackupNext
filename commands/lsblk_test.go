package commands

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/framps/raspiBackupNext/tools"
	"github.com/stretchr/testify/assert"
)

func TestLsblk(t *testing.T) {

	tools.NewLogger(false)

	files, err := filepath.Glob("lsblk_test/*.input")
	assert.NoErrorf(t, err, "Failed to retrieve testdata")

	t.Log("Starting TestLsblkid")

	for _, f := range files {
		t.Logf("Processing %s\n", f)

		parts := strings.Split(f, ".")
		inputFileName := parts[0] + ".input"
		expectedFileName := parts[0] + ".output"

		p, err := NewLsblkFromFile(inputFileName)
		assert.NoErrorf(t, err, "Unexpected error")

		expectedResult, err := ioutil.ReadFile(expectedFileName)
		assert.Equal(t, string(expectedResult), p.String())
		assert.NoError(t, err)
	}
}
