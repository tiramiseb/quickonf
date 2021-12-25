package apply

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands"
)

type ChangeMsg struct {
	Gidx int
}

type commandOutputs struct {
	m *model
}

func (m *model) commandOutputs() *commandOutputs {
	return &commandOutputs{m}
}

func (c *commandOutputs) NewCommandOutput(name string) commands.Output {
	out := &commandOutput{c.m, name, commands.StatusInfo, "Running..."}
	c.m.outputs = append(c.m.outputs, out)
	return out
}

type commandOutput struct {
	m       *model
	name    string
	status  commands.Status
	message string
}

func (c *commandOutput) Info(message string) {
	c.status = commands.StatusInfo
	c.message = message
	c.m.messages <- ChangeMsg{c.m.idx}
}

func (c *commandOutput) Infof(format string, a ...interface{}) {
	c.status = commands.StatusInfo
	c.message = fmt.Sprintf(format, a...)
	c.m.messages <- ChangeMsg{c.m.idx}
}

func (c *commandOutput) Success(message string) {
	c.status = commands.StatusSuccess
	c.message = message
	c.m.messages <- ChangeMsg{c.m.idx}
}

func (c *commandOutput) Successf(format string, a ...interface{}) {
	c.status = commands.StatusSuccess
	c.message = fmt.Sprintf(format, a...)
	c.m.messages <- ChangeMsg{c.m.idx}
}

func (c *commandOutput) Error(message string) {
	c.status = commands.StatusError
	c.message = message
	c.m.messages <- ChangeMsg{c.m.idx}
}

func (c *commandOutput) Errorf(format string, a ...interface{}) {
	c.status = commands.StatusError
	c.message = fmt.Sprintf(format, a...)
	c.m.messages <- ChangeMsg{c.m.idx}
}
