package discover

import (
	"bytes"
	"strings"

	"github.com/framps/raspiBackupNext/commands"
	"github.com/framps/raspiBackupNext/tools"
)

// System -
type System struct {
	SystemDevices *commands.SystemDevices
	BlkidDisks    *commands.BlkidDisks
	LsblkDisks    *commands.LsblkDisks
}

func NewSystem() *System {

	s := &System{}
	var err error

	s.SystemDevices, err = commands.NewSystemDevices()
	tools.HandleError(err)
	s.BlkidDisks, err = commands.NewBlkidDisks()
	tools.HandleError(err)
	s.LsblkDisks, err = commands.NewLsblkDisks()
	tools.HandleError(err)

	partedDisks := make([]*commands.PartedDisk, 0, len(s.LsblkDisks.Disks))
	for _, disk := range s.LsblkDisks.Disks {
		partedDisk, err := commands.NewPartedDisk(disk.Name)
		tools.HandleError(err)
		partedDisks = append(partedDisks, partedDisk)
	}

	return s
}

func (s System) String() string {
	var result bytes.Buffer

	sep := strings.Repeat("*", 30)
	result.WriteString(sep + "Systemdevices" + sep + "\n")
	result.WriteString(s.SystemDevices.String())
	result.WriteString(sep + "*** Blkid ***" + sep + "\n")
	result.WriteString(s.BlkidDisks.String())
	result.WriteString(sep + "*** Lsblk ***" + sep + "\n")
	result.WriteString(s.LsblkDisks.String())
	result.WriteString(sep + "*** Parted ***" + sep + "\n")
	for _, sd := range s.LsblkDisks.Disks {
		partedDisk, _ := commands.NewPartedDisk(sd.Name)
		result.WriteString(partedDisk.String())
	}

	return result.String()
}
