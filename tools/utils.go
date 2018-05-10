package tools

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import "strings"

// IsSpecialPartition -
func IsSpecialPartition(deviceName string) bool {
	return strings.HasPrefix(deviceName, "/dev/mmcblk") || strings.HasPrefix(deviceName, "/dev/loop")
}
