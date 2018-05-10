package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {

	cmds := []struct {
		Command   string
		Type      CommandType
		Args      []string
		StdoutErr string
		Err       bool
	}{
		{"echo", TypeNormal, []string{"Hello world"}, "Hello world\n", false},
		{"grep", TypeNormal, []string{"package command.go"}, "exit status 1", true},
		{"rdlprmpf", TypeNormal, []string{}, "executable file not found in $PATH", true},
		{"pwd", TypeNormal, []string{}, "commands", false},
		// {"file", TypeNormal, []string{"x"}, "cannot open", false},
		{"id", TypeNormal, []string{}, "uid=", false},

		{"sudo", TypeNormal, []string{"cat", "/etc/passwd"}, ":0:0:root:/root:/bin/bash", false},
		{"cat", TypeSudo, []string{"/etc/passwd"}, ":0:0:root:/root:/bin/bash", false},
		{"blkid", TypeSudo, []string{}, "/dev/sda", false},

		{"sh", TypeNormal, []string{"-c", "ls", "-la", "*"}, "command_test.go", false},
		{"ls", TypeBash, []string{"-la", "*"}, "command_test.go", false},
	}

	for _, c := range cmds {

		ts := TypeStrings[c.Type]
		if len(ts) > 0 {
			ts += " "
		}
		t.Logf("Testing command '%s%s %v'\n", ts, c.Command, c.Args)
		command := NewCommand(c.Type, c.Command, c.Args...)
		result, err := command.Execute()
		if err == nil {
			// t.Logf("Result: %s\n", string(*result))
		}

		if !c.Err {
			assert.NotNil(t, result)
			assert.NoError(t, err)
			assert.Contains(t, string(*result), c.StdoutErr)
		} else {
			assert.Error(t, err)
			if result != nil {
				assert.Contains(t, err.Error(), c.StdoutErr)
			}
		}
	}
}
