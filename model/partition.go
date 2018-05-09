package model

import (
	"bytes"
	"fmt"

	"github.com/framps/raspiBackup/go/commands"
	"github.com/framps/raspiBackup/go/tools"
	"github.com/jinzhu/copier"
)

// Partition -
type Partition struct {

	// from parted
	Name       string // /dev/sda1 or /dev/mmcblk0p1 or /dev/loop1
	Number     int
	Start      int64
	End        int64
	Size       int64
	Type       string // ext4
	FileSystem string
	Flags      string // lba

	// from blkid
	Uuid     string
	Partuuid string
	Label    string
	Ptype    string
}

// Disk -
type Disk struct {
	Name               string // /dev/sda or /dev/mmcblk0p or /dev/loop
	Size               string // 1000GB
	SectorSizeLogical  int    // 512
	SectorSizePhysical int    // 512
	PartitionTableType string // msdos
	Partitions         map[int]*Partition
}

func (d Disk) String() string {
	return fmt.Sprintf("Name: %s", d.Name)
}

// System -
type System struct {
	Disks         []*Disk
	Bootpartition *Partition
	Rootpartition *Partition
}

// NewSystem -
func NewSystem() *System {

	system := System{}

	blkidDisks, err := commands.NewBlkidDisks()
	if err != nil {
		tools.HandleError(err)
	}

	for _, d := range blkidDisks.Disks {
		disk := Disk{}
		copier.Copy(&disk, &d)
		system.Disks = append(system.Disks, &disk)
	}

	return &system
}

func (s System) String() string {
	var result bytes.Buffer
	for i, d := range s.Disks {
		result.WriteString(d.String())
		if i != len(s.Disks)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
