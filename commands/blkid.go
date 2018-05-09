package commands

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
	"go.uber.org/zap"
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

func (b BlkidDisks) String() string {
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

	logger := tools.Log

	blkid := BlkidDisks{make(map[string]*BlkidDisk, 16)}

	command := NewCommand(TypeSudo, "blkid")
	result, err := command.Execute()
	if err != nil {
		logger.Errorf("NewBlkid failed: %s", err.Error())
		return nil, err
	}

	logger.Debug(zap.String("Blkid", string(*result)))
	rdr := strings.NewReader(string(*result))

	blkid.parse(rdr)

	return &blkid, nil
}

/*
/dev/sda1: UUID="96ad35d1-85b1-45c8-941d-5c06e2ccc3c4" TYPE="ext4" PARTUUID="1de6ca19-01"
/dev/sdb1: UUID="dFRcFX-d3bX-y9Pp-Hts3-vPyK-hhcL-lAluXv" TYPE="LVM2_member" PARTUUID="6c96114a-01"
/dev/sdc: UUID="3s8MZp-Dxcd-UExv-JyGC-bXNG-Os8Q-mwqVpP" TYPE="LVM2_member"
/dev/sdd: UUID="CuaqcK-bAgL-GKC1-EiDy-RFyB-oOhd-ZmfRGX" TYPE="LVM2_member"
/dev/mapper/Second2-BigData: UUID="49d59d1d-fbd0-4f84-8b89-b8a0a6c20b74" TYPE="ext4"
/dev/mapper/data2-VMWare: UUID="785577ee-f43d-486d-a88a-3fd03eaac78e" TYPE="ext4"
/dev/mapper/data2-homeDisk: UUID="02d66411-a8c7-4ff1-be08-6f4a58f87255" TYPE="ext4"
/dev/mapper/data2-swap: UUID="a2065c48-447a-4c3d-8b49-10183805ede4" TYPE="swap"
/dev/mapper/Backup-System: UUID="df0a76d8-9810-4af2-a451-05dc17731445" TYPE="ext4"
/dev/mapper/Backup-Home: UUID="cface186-85ce-4a79-bbbf-a22f3aa7a838" TYPE="ext4
*/

func (b *BlkidDisks) parse(reader io.Reader) *BlkidDisks {

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {

		line := scanner.Text()

		r := regexp.MustCompile(`^(/dev/.+)([\d]*):`) // /dev/sda1: UUID="c6ccdbd5-12da-4b78-98c6-13cd63a733c7" TYPE="ext4" PARTUUID="000bee5a-01"

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
