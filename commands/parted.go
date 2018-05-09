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

// PartedPartition -
type PartedPartition struct {
	Name       string // /dev/sda1 or /dev/mmcblk0p1 or /dev/loop1
	Number     int
	Start      int64
	End        int64
	Size       int64
	Type       string // ext4
	FileSystem string
	Flags      string // lba
}

func (p PartedPartition) String() string {
	return fmt.Sprintf("Partition: %6s Partitionnumber: %d Start: %d End: %d Size: %d Type: %s", p.Name, p.Number, p.Start, p.End, p.Size, p.Type)
}

// PartedDisk -
type PartedDisk struct {
	Name               string // /dev/sda or /dev/mmcblk0p or /dev/loop
	Size               string // 1000GB
	SectorSizeLogical  int    // 512
	SectorSizePhysical int    // 512
	PartitionTableType string // msdos
	Partitions         map[int]*PartedPartition
}

func (d PartedDisk) String() string {
	var result bytes.Buffer

	result.WriteString(fmt.Sprintf("Disk: %s Size: %s Sector size: %d/%d PartitiontableType: %s\n", d.Name, d.Size, d.SectorSizeLogical, d.SectorSizePhysical, d.PartitionTableType))

	index := make([]*PartedPartition, 0, len(d.Partitions))
	for _, partition := range d.Partitions {
		index = append(index, partition)
	}

	sort.Slice(index, func(i, j int) bool {
		return index[i].Name < index[j].Name
	})

	for _, partition := range index {
		result.WriteString(fmt.Sprintf("%s\n", partition))
	}
	return result.String()
}

/*
BYT;
/dev/sde:15613952s:scsi:512:512:msdos:Generic STORAGE DEVICE:;
1:8192s:93813s:85622s:fat16::lba;
2:94208s:15613951s:15519744s:ext4::;
*/

func (d *PartedDisk) parse(reader io.Reader) *PartedDisk {

	logger := tools.Log

	scanner := bufio.NewScanner(reader)

	// /dev/sde:15613952s:scsi:512:512:msdos:Generic STORAGE DEVICE:;
	// /dev/mmcblk0:31116288s:sd/mmc:512:512:msdos:SD SL16G;
	r := regexp.MustCompile("^/dev/[^:]+:")
	for scanner.Scan() {
		line := scanner.Text()
		if r.MatchString(line) {
			parts := strings.Split(line, ":")
			d.Name, d.Size, d.PartitionTableType = parts[0], parts[1], parts[5]
			n, _ := strconv.Atoi(parts[3])
			d.SectorSizeLogical = n
			n, _ = strconv.Atoi(parts[4])
			d.SectorSizePhysical = n
			d.Partitions = make(map[int]*PartedPartition, 16)
			break
		}
	}

	// 1:8192s:93813s:85622s:fat16::lba;
	// 2:94208s:15613951s:15519744s:ext4::;
	r = regexp.MustCompile("^([0-9]+):")
	for scanner.Scan() {
		line := scanner.Text()
		if r.MatchString(line) {
			parts := strings.Split(line, ":")
			v, _ := strconv.Atoi(parts[0])
			pInfix := ""
			if tools.IsSpecialPartition(d.Name) {
				pInfix = "p"
			}
			name := fmt.Sprintf("%s%s%d", d.Name, pInfix, v)
			start, _ := strconv.ParseInt(parts[1][:len(parts[1])-1], 10, 64)
			end, _ := strconv.ParseInt(parts[2][:len(parts[2])-1], 10, 64)
			size, _ := strconv.ParseInt(parts[3][:len(parts[3])-1], 10, 64)
			partition := PartedPartition{Name: name,
				Number:     v - 1,
				Start:      start,
				End:        end,
				Size:       size,
				Type:       parts[4],
				FileSystem: parts[5],
				Flags:      parts[6][:len(parts[6])-1]}
			logger.Debug(zap.Any("Partition", partition))
			d.Partitions[v] = &partition
		} else {
			break
		}
	}

	return d
}

// NewPartedDisk -
func NewPartedDisk(diskDeviceName string) (*PartedDisk, error) {

	logger := tools.Log

	disk := PartedDisk{Partitions: make(map[int]*PartedPartition, 16)}

	command := NewCommand(TypeSudo, "parted", "-m", diskDeviceName, "unit", "s", "print")
	result, err := command.Execute()
	if err != nil {
		logger.Errorf("NewDisk failed for %s: %s", diskDeviceName, err.Error())
		return nil, err
	}

	logger.Debug(zap.String("Disk", string(*result)))

	rdr := strings.NewReader(string(*result))
	disk.parse(rdr)

	return &disk, nil
}

// NewPartedFromFile -
func NewPartedFromFile(fileName string) (*PartedDisk, error) {

	logger := tools.Log

	logger.Debugf("Filename: %s\n", fileName)

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Contents: %s\n", string(b))

	disk := PartedDisk{Partitions: make(map[int]*PartedPartition, 16)}

	rdr := strings.NewReader(string(b))
	disk.parse(rdr)

	return &disk, nil
}

// NewPartedToFile -
func NewPartedToFile(fileName, deviceName string) error {

	b, err := NewPartedDisk(deviceName)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, []byte(b.String()), os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}
