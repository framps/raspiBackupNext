package model

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

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
	return fmt.Sprintf("PartitionNumber: %d - Start: %d - End: %d - Size: %d - Type: %s - FileSystem: %s - Flags: %s "+
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

	index := make([]*Partition, 0, len(d.Partitions))
	for _, partition := range d.Partitions {
		index = append(index, partition)
	}

	sort.Slice(index, func(i, j int) bool {
		return index[i].Name < index[j].Name
	})

	for _, partition := range index {
		result.WriteString(fmt.Sprintf("%s", partition))
		result.WriteString("\n")
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
func NewSystem() (*System, error) {

	logger := tools.Log

	system := System{}

	// retrieve all known disks of system
	lsblkDisks, err := commands.NewLsblkDisks()
	tools.HandleError(err)

	blkidDisk, err := commands.NewBlkidDisks()
	tools.HandleError(err)

	fmt.Printf("*** %s\n", blkidDisk)

	for _, d := range lsblkDisks.Disks {

		logger.Debugf("Processing disk %s", d.Name)
		disk := Disk{Name: d.Name}

		partedDisk, err := commands.NewPartedDisk("/dev/" + disk.Name)
		tools.HandleError(err)

		copier.Copy(&disk, &partedDisk)
		system.Disks = append(system.Disks, &disk)

		disk.Partitions = make(map[int]*Partition, len(d.Partitions))

		for i, p := range partedDisk.Partitions {
			partition := Partition{}
			copier.Copy(&partition, partedDisk.Partitions[i])
			disk.Partitions[p.Number] = &partition
		}
	}

	return &system, nil
}

func (s System) String() string {
	var result bytes.Buffer

	if len(s.Disks) > 0 {
		index := make([]*Disk, 0, len(s.Disks))
		for _, disk := range s.Disks {
			index = append(index, disk)
		}

		sort.Slice(index, func(i, j int) bool {
			return index[i].Name < index[j].Name
		})

		for i := range index {
			result.WriteString(index[i].String())
			if i != len(index)-1 {
				result.WriteString("\n")
			}
		}
	}

	return result.String()
}

// NewSystemToJSON - -
func NewSystemToJSON(fileName string) error {

	b, err := NewSystem()
	if err != nil {
		return err
	}

	var j []byte
	if j, err = json.MarshalIndent(b, "", " "); err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, j, os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}

// NewSystemFromJson -
func NewSystemFromJson(fileName string) (*System, error) {

	j, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var s System
	if err := json.Unmarshal(j, &s); err != nil {
		return nil, err
	}

	return &s, nil

}
