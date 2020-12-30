package helper

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

// SudoPassword must contain the sudo password, in order to user ExecSudo. This password may be set with the sudo-password instruction.
var SudoPassword string

// Exec executes a command as the current user, returning its stdin and stdout outputs combined.
//
// If command fails, its combined output is in the error message.
func Exec(env []string, cmd string, args ...string) ([]byte, error) {
	return execute(env, cmd, args, nil)
}

// ExecSudo executes a command as root by using sudo, returning its stdin and stdout outputs combined.
//
// If command fails, its combined output is in the error message.
func ExecSudo(env []string, args ...string) ([]byte, error) {
	if SudoPassword == "" {
		return nil, errors.New("Sudo password is not set")
	}
	args = append([]string{"--prompt=", "--stdin"}, args...)
	return execute(env, "sudo", args, strings.NewReader(SudoPassword))
}

func execute(env []string, cmd string, args []string, stdin io.Reader) ([]byte, error) {
	command := exec.Command(cmd, args...)
	command.Env = append(os.Environ(), "LANG=C")
	command.Env = append(command.Env, env...)
	command.Stdin = stdin
	cmdout, err := command.CombinedOutput()
	if err != nil {
		if len(cmdout) > 0 {
			return nil, errors.New(string(cmdout))
		}
		return nil, err
	}
	return cmdout, nil
}
