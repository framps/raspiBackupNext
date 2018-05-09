package tools

import "strings"

// IsSpecialPartition -
func IsSpecialPartition(deviceName string) bool {
	return strings.HasPrefix(deviceName, "/dev/mmcblk") || strings.HasPrefix(deviceName, "/dev/loop")
}
