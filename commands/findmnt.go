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
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/framps/raspiBackupNext/tools"
)

const notFound = "Not found"

// SystemDevice -
type SystemDevice struct {
	FullName        string // /dev/sda1, /dev/mmcblk3p2
	Number          int    // 1, 2
	Disk            string // sda, mmcblk3
	PartitionName   string // sda, mmcblk3p
	LocatedOnSDCard bool
}

func (p SystemDevice) String() string {
	return fmt.Sprintf("FullName: %s Number: %d Disk: %s LocatedOnSDCard: %t\n", p.FullName, p.Number, p.Disk, p.LocatedOnSDCard)
}

// NewSystemDevice -
func NewSystemDevice(fullName string) (*SystemDevice, error) {

	logger := tools.Log

	p := &SystemDevice{FullName: fullName, LocatedOnSDCard: false}

	re := regexp.MustCompile("^/dev/([a-z]+)(([0-9]+)p)?([0-9]+)$")
	matches := re.FindStringSubmatch(fullName)

	if len(matches) != 5 {
		e := fmt.Errorf("Illegal systempartition %s", fullName)
		logger.Debug(e)
		return nil, e
	}

	// ["/dev/mmcblk3p1" "mmcblk" "3p" "3" "1"]
	// ["/dev/sda3" "sda" "" "" "3"]

	p.Number, _ = strconv.Atoi(matches[4])
	p.Disk = matches[1]
	p.PartitionName = matches[1]

	if matches[1] == "mmcblk" {
		p.Disk += matches[3]
		p.PartitionName += matches[2]
		p.LocatedOnSDCard = true
	}
	return p, nil
}

// SystemDevices -
type SystemDevices struct {
	Bootdevice *SystemDevice // /dev/sda1 or /dev/mmcblk0p1 or /dev/loop1
	Rootdevice *SystemDevice // /dev/root or /dev/sda1
}

func (d *SystemDevices) print(device *SystemDevice, buffer *bytes.Buffer, name string) *bytes.Buffer {
	buffer.WriteString(fmt.Sprintf("%s: %s", name, device))
	return buffer
}

func (d SystemDevices) String() string {
	var buffer bytes.Buffer
	d.print(d.Bootdevice, &buffer, "BootDevice")
	buffer.WriteString("\n")
	d.print(d.Rootdevice, &buffer, "RootDevice")
	return buffer.String()
}

// NewSystemDevices -
func NewSystemDevices() (*SystemDevices, error) {

	var (
		bootDevice string
		rootDevice string
	)

	logger := tools.Log

	command := NewCommand(TypeSudo, "findmnt", "/boot", "-o", "source", "-n")
	result, err := command.Execute()
	if err != nil {
		logger.Debugf("NewSystemDevices /boot failed: %s", err.Error())
	} else {
		rdr := strings.NewReader(string(*result))
		scanner := bufio.NewScanner(rdr)
		if scanner.Scan() {
			bootDevice = scanner.Text()
		} else {
			bootDevice = notFound
		}
	}

	command = NewCommand(TypeSudo, "findmnt", "/", "-o", "source", "-n")
	result, err = command.Execute()
	if err != nil {
		logger.Debugf("NewSystemDevices / failed: %s", err.Error())
	} else {
		rdr := strings.NewReader(string(*result))
		scanner := bufio.NewScanner(rdr)
		if scanner.Scan() {
			rootDevice = scanner.Text()
		} else {
			bootDevice = notFound
		}
	}

	systemDevice := SystemDevices{}
	bootSystemDevice, err := NewSystemDevice(bootDevice)
	if err == nil {
		systemDevice.Bootdevice = bootSystemDevice
	}
	rootSystemDevice, err := NewSystemDevice(rootDevice)
	if err == nil {
		systemDevice.Rootdevice = rootSystemDevice
	}

	return &systemDevice, nil
}

// NewSystemDevicesToFile -
func NewSystemDevicesToFile(fileName string) error {

	b, err := NewSystemDevices()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, []byte(b.String()), os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}
