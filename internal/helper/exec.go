package helper

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// SudoPassword must contain the sudo password, in order to user ExecSudo. This password may be set with the sudo-password instruction.
var SudoPassword string

// Exec executes a command
func Exec(env []string, cmd string, args ...string) ([]byte, error) {
	cmdObj := exec.Command(cmd, args...)
	cmdObj.Env = append(os.Environ(), "LANG=C")
	cmdObj.Env = append(cmdObj.Env, env...)
	cmdout, err := cmdObj.CombinedOutput()
	if err != nil && len(cmdout) > 0 {
		return nil, errors.New(string(cmdout))
	}
	return cmdout, err
}

// ExecSudo executes a command as root by using sudo
func ExecSudo(env []string, args ...string) ([]byte, error) {
	if SudoPassword == "" {
		return nil, errors.New("Sudo password is not set")
	}
	args = append([]string{"--prompt=", "--stdin"}, args...)
	sudoCmd := exec.Command("sudo", args...)
	sudoCmd.Env = append(os.Environ(), "LANG=C")
	sudoCmd.Env = append(sudoCmd.Env, env...)
	sudoCmd.Stdin = strings.NewReader(SudoPassword)
	cmdout, err := sudoCmd.CombinedOutput()
	if err != nil && len(cmdout) > 0 {
		return nil, errors.New(string(cmdout))
	}
	return cmdout, err
}
