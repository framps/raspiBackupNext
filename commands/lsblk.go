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
	"strings"

	"github.com/framps/raspiBackup/go/tools"
	"go.uber.org/zap"
)

// LsblkDisk -
type LsblkDisk struct {
	Name string
	Size string
}

func (d LsblkDisk) String() string {
	return fmt.Sprintf("Name: %s Size: %s", d.Name, d.Size)
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

	command := NewCommand(TypeSudo, "lsblk", "-P", "-d")
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
NAME="sda" MAJ:MIN="8:0" RM="0" SIZE="931.5G" RO="0" TYPE="disk" MOUNTPOINT=""
NAME="sdb" MAJ:MIN="8:16" RM="0" SIZE="931.5G" RO="0" TYPE="disk" MOUNTPOINT=""
NAME="loop0" MAJ:MIN="7:0" RM="0" SIZE="3.7G" RO="0" TYPE="loop" MOUNTPOINT=""
NAME="mmcblk0" MAJ:MIN="179:0" RM="0" SIZE="14.9G" RO="0" TYPE="disk" MOUNTPOINT=""
*/

func (d *LsblkDisks) parse(reader io.Reader) *LsblkDisks {

	scanner := bufio.NewScanner(reader)

	r := regexp.MustCompile(`NAME="([^"]*)".*SIZE="([^"]*)"`)

	for scanner.Scan() {
		line := scanner.Text()
		if matchGroup := r.FindAllStringSubmatch(line, -1); matchGroup != nil {
			m := matchGroup[0]
			name := "/dev/" + m[1]
			size := m[2]
			d.Disks[name] = &LsblkDisk{Name: name, Size: size}
		}
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
