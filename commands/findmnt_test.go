package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindmnt(t *testing.T) {

	cmds := []struct {
		FullName        string // /dev/sda1, /dev/mmcblk3p2
		Number          int    // 1, 2
		Disk            string // sda, mmcblk3
		PartitionName   string // sda, mmcblk3p
		LocatedOnSDCard bool
		Err             bool
	}{
		// success
		{"/dev/sda1", 1, "sda", "sda", false, false},
		{"/dev/mmcblk1p3", 3, "mmcblk1", "mmcblk1p", true, false},
		{"/dev/loop7", 7, "loop", "loop", false, false},
		// failures
		{"loop7", 3, "loop", "loop", true, true},
		{"/dev/578", 3, "loop", "loop7", true, true},
		{"dev/sda", 3, "loop", "loop7", true, true},
	}

	for _, c := range cmds {

		t.Logf("Testing device %s\n", c.FullName)
		device, err := NewSystemDevice(c.FullName)

		if !c.Err {
			assert.NotNil(t, device)
			assert.NoError(t, err)
			assert.Equal(t, device.FullName, c.FullName)
			assert.Equal(t, device.Number, c.Number)
			assert.Equal(t, device.Disk, c.Disk)
			assert.Equal(t, device.PartitionName, c.PartitionName)
			assert.Equal(t, device.LocatedOnSDCard, c.LocatedOnSDCard)
		} else {
			assert.Error(t, err)
		}
	}
}
