package helper

import (
	"errors"
	"os/exec"
	"strings"
)

var SudoPassword string

func Exec(cmd string, args ...string) ([]byte, error) {
	cmdout, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil && len(cmdout) > 0 {
		return nil, errors.New(string(cmdout))
	}
	return cmdout, err
}

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
