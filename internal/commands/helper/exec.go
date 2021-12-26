package helper

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/user"
)

func ExecOutAsLines(env []string, cmd string, args ...string) ([]string, error) {
	var buf bytes.Buffer
	wait, err := Exec(env, &buf, cmd, args...)
	if err != nil {
		return nil, err
	}
	if err := wait(); err != nil {
		return nil, err
	}
	var out []string
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		out = append(out, scanner.Text())
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Exec executes a command
func Exec(env []string, output io.Writer, cmd string, args ...string) (wait func() error, err error) {
	command := execCmd(env, output, cmd, args...)
	err = command.Start()
	return command.Wait, err
}

func ExecAs(usr *user.User, env []string, output io.Writer, cmd string, args ...string) (wait func() error, err error) {
	args = append([]string{"-u", usr.Username, "--", cmd}, args...)
	command := execCmd(env, output, "runuser", args...)
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
