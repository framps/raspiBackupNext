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

// LsblkDisk -
type LsblkDisk struct {
	Name       string
	Partitions map[int]*LsblkPartition
}

func (d LsblkDisk) String() string {
	var result bytes.Buffer
	result.WriteString(fmt.Sprintf("Name: %s", d.Name))
	if len(d.Partitions) > 0 {
		result.WriteString(" - ")
		for i := range d.Partitions {
			result.WriteString(d.Partitions[i].String())
		}
	}
	return result.String()
}

// NewLsblkDisk -
func NewLsblkDisk(name string) *LsblkDisk {
	disk := LsblkDisk{Name: name}
	disk.Partitions = make(map[int]*LsblkPartition)
	return &disk
}

// LsblkPartition -
type LsblkPartition struct {
	Number     int
	MajMin     string
	Rm         string
	Size       int64
	Ro         string
	Type       string
	Mountpoint string
}

func (p LsblkPartition) String() string {
	var result bytes.Buffer
	result.WriteString(fmt.Sprintf("Number: %d - MajMin: %s - RM: %s - Size: %d - RO: %s - Type: %s - Mountpoint: %s", p.Number, p.MajMin, p.Rm, p.Size, p.Ro, p.Type, p.Mountpoint))
	return result.String()
}

// LsblkDisks -
type LsblkDisks struct {
	Disks map[string]*LsblkDisk
}

func (d LsblkDisks) String() string {
	var result bytes.Buffer
	result.WriteString(fmt.Sprintf("Disks:\n"))

	index := make([]*LsblkDisk, 0, len(d.Disks))

	for _, disk := range d.Disks {
		index = append(index, disk)
	}

	sort.Slice(index, func(i, j int) bool {
		return index[i].Name < index[j].Name
	})

	for i := range index {
		result.WriteString(fmt.Sprintf("%s\n", *(index[i])))
	}

	return result.String()
}

// NewLsblkDisks -
func NewLsblkDisks() (*LsblkDisks, error) {

	logger := tools.Log

	lsblkids := LsblkDisks{make(map[string]*LsblkDisk, 16)}

	command := NewCommand(TypeSudo, "lsblk", "-r", "-n", "-b")
	result, err := command.Execute()
	if err != nil {
		logger.Errorf("NewLsblkid failed: %s", err.Error())
		return nil, err
	}

	logger.Debug(zap.String("Lsblkid", string(*result)))

	rdr := strings.NewReader(string(*result))
	lsblkids.parse(rdr)

	return &lsblkids, nil
}

/*
sda 8:0 0 250059350016 0 disk
sda1 8:1 0 231054770176 0 part /
sdb 8:16 0 1000204886016 0 disk
sdb1 8:17 0 1000202241024 0 part
Backup-System 252:2 0 329231892480 0 lvm /backup/system
Backup-Home 252:5 0 670967005184 0 lvm /backup/home
sdc 8:32 0 2000398934016 0 disk
Second2-BigData 252:0 0 1073741824000 0 lvm /disks/bigdata
sdd 8:48 0 2000398934016 0 disk
data2-VMWare 252:1 0 429496729600 0 lvm /disks/VMware
data2-homeDisk 252:3 0 322122547200 0 lvm /disks/homeDisk
data2-swap 252:4 0 8589934592 0 lvm
*/

func (d *LsblkDisks) parse(reader io.Reader) *LsblkDisks {

	scanner := bufio.NewScanner(reader)

	re := regexp.MustCompile(`[^\d]+([\d]+)`) // sda or sda1

	var disk *LsblkDisk

	for scanner.Scan() {
		line := scanner.Text()
		elements := strings.Split(line, " ")

		if elements[5] == "disk" {
			if disk != nil {
				d.Disks[disk.Name] = disk
			}
			disk = NewLsblkDisk(elements[0])
			continue
		} else if elements[5] == "part" {
			matches := re.FindStringSubmatch(elements[0])
			partitionNumberString := matches[1]
			partitionNumber, _ := strconv.Atoi(partitionNumberString)
			partition := LsblkPartition{}
			size, _ := strconv.ParseInt(elements[3], 10, 64)
			partition.Number, partition.MajMin, partition.Rm, partition.Size, partition.Ro, partition.Type, partition.Mountpoint =
				partitionNumber, elements[1], elements[2], size, elements[4], elements[5], elements[6]
			disk.Partitions[partitionNumber] = &partition
		}
	}
	if disk != nil {
		d.Disks[disk.Name] = disk
	}

	return d
}

// NewLsblkFromFile -
func NewLsblkFromFile(fileName string) (*LsblkDisks, error) {

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	blkid := LsblkDisks{make(map[string]*LsblkDisk, 16)}

	rdr := strings.NewReader(string(b))
	blkid.parse(rdr)

	return &blkid, nil

}

// NewLsblkToFile -
func NewLsblkToFile(fileName string) error {

	b, err := NewLsblkDisks()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, []byte(b.String()), os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}
