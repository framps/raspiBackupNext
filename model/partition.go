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

	"github.com/framps/raspiBackupNext/commands"
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
func NewSystem() (*System, error) {

	system := System{}

	blkidDisks, err := commands.NewBlkidDisks()
	if err != nil {
		return nil, err
	}

	for _, d := range blkidDisks.Disks {
		disk := Disk{}
		copier.Copy(&disk, &d)
		system.Disks = append(system.Disks, &disk)
	}

	return &system, nil
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

// NewSystemToJson - -
func NewSystemToJson(fileName string) error {

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
