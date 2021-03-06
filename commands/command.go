package commands

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/framps/raspiBackupNext/tools"
)

// Cmd -
type Cmd struct {
	*exec.Cmd
}

// String -
func (c Cmd) String() string {
	return fmt.Sprintf("%s %s", c.Path, c.Args)
}

// CommandType -
type CommandType int

const (
	// TypeNormal -
	TypeNormal CommandType = iota
	// TypeSudo -
	TypeSudo
	// TypeBash -
	TypeBash
)

// TypeStrings -
var TypeStrings = [...]string{"", "sudo", "bash"}

func (t CommandType) String() string {
	return TypeStrings[t]
}

// NewCommand -
func NewCommand(commandType CommandType, command string, args ...string) *Cmd {

	var result *Cmd
	switch commandType {
	case TypeNormal:
		result = &Cmd{exec.Command(command, args...)}
	case TypeSudo:
		result = &Cmd{exec.Command("sudo", append([]string{command}, args...)...)}
	case TypeBash:
		result = &Cmd{exec.Command("sh", append([]string{"-c"}, command, strings.Join(args, " "))...)}
	default:
		tools.HandleError(fmt.Errorf("Invalid command type %s", commandType))
	}

	return result
}

// Execute -
func (c *Cmd) Execute() (*[]byte, error) {
	tools.Logger.Debug("Executing command ", c.Path, c.Args)
	stdoutStderr, err := c.CombinedOutput()
	if err != nil {
		tools.Logger.Errorf("Command error: %s", err)
		return &stdoutStderr, err
	}
	tools.Logger.Debugf("Command result:\n%s", stdoutStderr)
	return &stdoutStderr, nil

}
