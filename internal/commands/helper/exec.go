package helper

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"os/user"
)

// Exec executes a command
func Exec(env []string, output io.Writer, cmd string, args ...string) error {
	return execCmd(env, output, cmd, args...)
}

// ExecAs executes a command as an user
func ExecAs(usr *user.User, env []string, output io.Writer, cmd string, args ...string) error {
	args = append([]string{"-u", usr.Username, "--", cmd}, args...)
	return execCmd(env, output, "runuser", args...)
}

func execCmd(env []string, output io.Writer, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Env = append(os.Environ(), "LANG=C")
	c.Env = append(c.Env, env...)
	var errBuf bytes.Buffer
	if output == nil {
		c.Stdout = &errBuf
		c.Stderr = &errBuf
	} else {
		c.Stdout = output
		c.Stderr = output
	}
	err := c.Run()
	if err != nil && output == nil {
		if ee, ok := err.(*exec.ExitError); ok {
			ee.Stderr = bytes.ReplaceAll(errBuf.Bytes(), []byte("\n"), []byte(" "))
		}
	}
	return err
}

func ExecErr(err error) string {
	if err == nil {
		return ""
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return string(exitErr.Stderr)
	}
	return err.Error()
}
