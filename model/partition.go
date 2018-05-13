package model

import (
	"bytes"
	"fmt"

	"github.com/framps/raspiBackupNext/commands"
	"github.com/framps/raspiBackupNext/tools"
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

func (p Partition) String() string {
	return fmt.Sprintf("PartitionNumber: %d - StartBlock: %d - EndBlock: %d - Size: %d - Type: %s - FileSystem: %s - Flags: %s "+
		"UUid: %s - Partuuid: %s - Label: %s - PType: %s",
		p.Number, p.Start, p.End, p.Size, p.Type, p.FileSystem, p.Flags,
		p.Uuid, p.Partuuid, p.Label, p.Ptype)
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
	var result bytes.Buffer
	result.WriteString(fmt.Sprintf("Name: %s - Size: %s - LogicalSectorSize: %d - PhysicalSectorSize: %d - PartitionTableType: %s\n",
		d.Name, d.Size, d.SectorSizeLogical, d.SectorSizePhysical, d.PartitionTableType))

	for _, p := range d.Partitions {
		result.WriteString(p.String())
	}
	return result.String()
}

// System -
type System struct {
	Disks         []*Disk
	Bootpartition *Partition
	Rootpartition *Partition
}

// NewSystem -
func NewSystem() *System {

	logger := tools.Log

	system := System{}

	// retrieve all known disks of system
	blkidDisks, err := commands.NewBlkidDisks()
	tools.HandleError(err)

	for _, d := range blkidDisks.Disks {

		logger.Debugf("Processing disk %s", d.Name)
		disk := Disk{Name: d.Name}

		partedDisk, err := commands.NewPartedDisk(disk.Name)
		tools.HandleError(err)
		fmt.Printf("---> %s\n", partedDisk)
		copier.Copy(&disk, &partedDisk)
		system.Disks = append(system.Disks, &disk)

		disk.Partitions = make(map[int]*Partition, len(d.Partitions))

		for _, p := range d.Partitions {
			partition := Partition{}
			//copier.Copy(&partition, &p)
			fmt.Printf("%#v\n", *partedDisk.Partitions[p.Number])
			copier.Copy(&partition, partedDisk.Partitions[p.Number])
			disk.Partitions[p.Number] = &partition
		}
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
