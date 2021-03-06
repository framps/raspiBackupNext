package commands

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindmnt(t *testing.T) {

	cmds := []struct {
		DeviceName      string // /dev/sda1, /dev/mmcblk3p2
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

		t.Logf("Testing device %s\n", c.DeviceName)
		device, err := NewSystemDevice(c.DeviceName)

		if !c.Err {
			assert.NotNil(t, device)
			assert.NoError(t, err)
			assert.Equal(t, device.DeviceName, c.DeviceName)
			assert.Equal(t, device.Number, c.Number)
			assert.Equal(t, device.Disk, c.Disk)
			assert.Equal(t, device.PartitionName, c.PartitionName)
			assert.Equal(t, device.LocatedOnSDCard, c.LocatedOnSDCard)
		} else {
			assert.Error(t, err)
		}
	}
}
