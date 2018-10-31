package commands

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/framps/raspiBackupNext/tools"
)

// BlkidPartition -
type BlkidPartition struct {
	Number   int
	Uuid     string
	Type     string
	Partuuid string
	Label    string
	Pttype   string
}

func (b BlkidPartition) String() string {
	return fmt.Sprintf("Partition: %d Uuid: %40s Type: %6s Partuuid: %6s Pttype: %6s Label: %6s", b.Number, b.Uuid, b.Type, b.Partuuid, b.Pttype, b.Label)
}

// BlkidDisk -
type BlkidDisk struct {
	Name       string // e.g. /dev/sda
	Partitions map[int]*BlkidPartition
}

func (b BlkidDisk) String() string {
	var result bytes.Buffer
	result.WriteString(fmt.Sprintf("Disk: %s\n", b.Name))

	index := make([]*BlkidPartition, 0, len(b.Partitions))

	for _, part := range b.Partitions {
		index = append(index, part)
	}

	sort.Slice(index, func(i, j int) bool {
		return index[i].Number < index[j].Number
	})

	for i := range index {
		result.WriteString(fmt.Sprintf("%s\n", *(index[i])))
	}

	return result.String()
}

// BlkidDisks -
type BlkidDisks struct {
	Disks map[string]*BlkidDisk
}

func (b *BlkidDisks) String() string {
	var result bytes.Buffer

	index := make([]*BlkidDisk, 0, len(b.Disks))
	for _, disk := range b.Disks {
		index = append(index, disk)
	}

	sort.Slice(index, func(i, j int) bool {
		return index[i].Name < index[j].Name
	})

	for _, disk := range index {
		result.WriteString(fmt.Sprintf("%s", disk))
	}
	return result.String()
}

// NewBlkidDisks -
func NewBlkidDisks() (*BlkidDisks, error) {

	logger := tools.Logger

	blkid := BlkidDisks{make(map[string]*BlkidDisk, 16)}

	command := NewCommand(TypeSudo, "blkid")
	result, err := command.Execute()
	if err != nil {
		logger.Errorf("NewBlkid failed: %s", err.Error())
		return nil, err
	}

	rdr := strings.NewReader(string(*result))

	blkid.parse(rdr)

	return &blkid, nil
}

/*
/dev/sda1: UUID="96ad35d1-85b1-45c8-941d-5c06e2ccc3c4" TYPE="ext4" PARTUUID="1de6ca19-01"
/dev/sdb1: UUID="dFRcFX-d3bX-y9Pp-Hts3-vPyK-hhcL-lAluXv" TYPE="LVM2_member" PARTUUID="6c96114a-01"
/dev/sdc: UUID="3s8MZp-Dxcd-UExv-JyGC-bXNG-Os8Q-mwqVpP" TYPE="LVM2_member"
/dev/sdd: UUID="CuaqcK-bAgL-GKC1-EiDy-RFyB-oOhd-ZmfRGX" TYPE="LVM2_member"
/dev/mmcblk0p2: UUID="64a5e86f-5ed3-4c9f-aab3-c4ae24bff95a" TYPE="ext4" LABEL="system"
/dev/mmcblk0p1: SEC_TYPE="msdos" UUID="3312-932F" TYPE="vfat"
/dev/sdb1: LABEL="black" UUID="76e7d7d4-f6e9-4867-83a7-03eaa3fc878d" TYPE="ext4"
/dev/sda1: LABEL="silver" UUID="8095dbdf-9b0a-4dda-9352-d56366af43c8" TYPE="ext4"
*/

func (b *BlkidDisks) parse(reader io.Reader) *BlkidDisks {

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {

		line := scanner.Text()

		r := regexp.MustCompile(`^(/dev/[a-z]+(?:[0-9]+p)?)([\d]+):`) // /dev/sda1: UUID="c6ccdbd5-12da-4b78-98c6-13cd63a733c7" TYPE="ext4" PARTUUID="000bee5a-01"

		var disk *BlkidDisk
		if matchGroup := r.FindAllStringSubmatch(line, -1); matchGroup != nil {

			d := matchGroup[0]

			diskname := d[1]

			if e, ok := b.Disks[diskname]; !ok {
				disk = &BlkidDisk{Name: diskname, Partitions: make(map[int]*BlkidPartition, 16)}
				b.Disks[disk.Name] = disk
			} else {
				disk = e
			}

			partitionNumber, _ := strconv.Atoi(d[2])
			parts := strings.Split(line, " ")

			partition := BlkidPartition{Number: partitionNumber}
			for i := range parts {
				e := strings.Split(parts[i], "=")
				switch e[0] {
				case "UUID":
					partition.Uuid = strings.Replace(e[1], `"`, "", -1)
				case "LABEL":
					partition.Label = strings.Replace(e[1], `"`, "", -1)
				case "TYPE":
					partition.Type = strings.Replace(e[1], `"`, "", -1)
				case "PARTUUID":
					partition.Partuuid = strings.Replace(e[1], `"`, "", -1)
				case "PTTYPE":
					partition.Pttype = strings.Replace(e[1], `"`, "", -1)
				}
				disk.Partitions[partitionNumber] = &partition
			}
		}
	}
	return b
}

// NewBlkidFromFile -
func NewBlkidFromFile(fileName string) (*BlkidDisks, error) {

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	blkid := BlkidDisks{make(map[string]*BlkidDisk, 16)}

	rdr := strings.NewReader(string(b))
	blkid.parse(rdr)

	return &blkid, nil

}

// NewBlkidToFile -
func NewBlkidToFile(fileName string) error {

	b, err := NewBlkidDisks()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, []byte(b.String()), os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}
