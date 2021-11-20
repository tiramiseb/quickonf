package helper

import (
	"io"
	"os"
	"os/exec"
)

// Exec executes a command
func Exec(env []string, output io.Writer, cmd string, args ...string) (wait func() error, err error) {
	command := exec.Command(cmd, args...)
	command.Env = append(os.Environ(), "LANG=C")
	command.Env = append(command.Env, env...)
	command.Stdout = output
	command.Stderr = output
	err = command.Start()
	return command.Wait, err
}
