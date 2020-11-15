package helper

import (
	"errors"
	"os/exec"
	"strings"
)

// SudoPassword must contain the sudo password, in order to user ExecSudo. This password may be set with the sudo-password instruction.
var SudoPassword string

// Exec executes a command
func Exec(cmd string, args ...string) ([]byte, error) {
	cmdout, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil && len(cmdout) > 0 {
		return nil, errors.New(string(cmdout))
	}
	return cmdout, err
}

// ExecSudo executes a command as root by using sudo
func ExecSudo(args ...string) error {
	if SudoPassword == "" {
		return errors.New("Sudo password is not set")
	}
	args = append([]string{"--prompt=", "--stdin", "--reset-timestamp"}, args...)
	sudoCmd := exec.Command("sudo", args...)
	sudoCmd.Env = []string{"LANG=C"}
	sudoCmd.Stdin = strings.NewReader(SudoPassword)
	cmdout, err := sudoCmd.CombinedOutput()
	if err != nil && len(cmdout) > 0 {
		return errors.New(string(cmdout))
	}
	return err
}
