package discover

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import (
	"bytes"
	"strings"
	"sync"

	"github.com/framps/raspiBackupNext/commands"
	"github.com/framps/raspiBackupNext/tools"
)

// System -
type System struct {
	SystemDevices *commands.SystemDevices
	BlkidDisks    *commands.BlkidDisks
	LsblkDisks    *commands.LsblkDisks
}

func NewSystem(parallelExecution bool) *System {

	s := System{}
	var err error

	if parallelExecution {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			s.SystemDevices, err = commands.NewSystemDevices()
			tools.HandleError(err)
			wg.Done()
		}()
		go func() {
			s.BlkidDisks, err = commands.NewBlkidDisks()
			tools.HandleError(err)
			wg.Done()
		}()
		go func() {
			s.LsblkDisks, err = commands.NewLsblkDisks()
			tools.HandleError(err)
			wg.Done()
		}()
		wg.Wait()
	} else {
		s.SystemDevices, err = commands.NewSystemDevices()
		s.BlkidDisks, err = commands.NewBlkidDisks()
		s.LsblkDisks, err = commands.NewLsblkDisks()
	}

	partedDisks := make([]*commands.PartedDisk, 0, len(s.LsblkDisks.Disks))
	for _, disk := range s.LsblkDisks.Disks {
		partedDisk, err := commands.NewPartedDisk("/dev/" + disk.Name)
		tools.HandleError(err)
		partedDisks = append(partedDisks, partedDisk)
	}

	return &s
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
		partedDisk, _ := commands.NewPartedDisk("/dev/" + sd.Name)
		result.WriteString(partedDisk.String())
	}
	return result.String()
}
