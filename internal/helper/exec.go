package helper

import (
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

// Exec executes a command
func Exec(env []string, output io.Writer, cmd string, args ...string) (wait func() error, err error) {
	command := execCmd(env, output, cmd, args...)
	err = command.Start()
	return command.Wait, err
}

func ExecAs(usr *user.User, env []string, output io.Writer, cmd string, args ...string) (wait func() error, err error) {
	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return nil, err
	}
	gid, err := strconv.Atoi(usr.Gid)
	if err != nil {
		return nil, err
	}
	command := execCmd(env, output, cmd, args...)
	command.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}
	err = command.Start()
	return command.Wait, err
}

func execCmd(env []string, output io.Writer, cmd string, args ...string) *exec.Cmd {
	command := exec.Command(cmd, args...)
	command.Env = append(os.Environ(), "LANG=C")
	command.Env = append(command.Env, env...)
	command.Stdout = output
	command.Stderr = output
	return command
}
