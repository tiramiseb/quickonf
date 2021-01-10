package helper

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// SudoPassword must contain the sudo password, in order to user ExecSudo. This password may be set with the sudo-password instruction.
var SudoPassword string

// Exec executes a command as the current user, returning its stdout and stderr outputs combined.
//
// If command fails (errcode≠0), its combined output is in the error message.
func Exec(env []string, stdin string, cmd string, args ...string) ([]byte, error) {
	return execute(env, cmd, args, stdin)
}

// ExecSudo executes a command as root by using sudo, returning its stdout and stderr outputs combined.
//
// If command fails (errcode≠0), its combined output is in the error message.
func ExecSudo(env []string, stdin string, args ...string) ([]byte, error) {
	if SudoPassword == "" {
		return nil, errors.New("Sudo password is not set")
	}
	args = append([]string{"--reset-timestamp", "--prompt=", "--stdin"}, args...)
	if stdin == "" {
		stdin = SudoPassword
	} else {
		stdin = SudoPassword + "\n" + stdin
	}
	return execute(env, "sudo", args, stdin)
}

func execute(env []string, cmd string, args []string, stdin string) ([]byte, error) {
	command := exec.Command(cmd, args...)
	command.Env = append(os.Environ(), "LANG=C")
	command.Env = append(command.Env, env...)
	command.Stdin = strings.NewReader(stdin)
	cmdout, err := command.CombinedOutput()
	if err != nil {
		if len(cmdout) > 0 {
			return nil, errors.New(string(cmdout))
		}
		return nil, err
	}
	return cmdout, nil
}
