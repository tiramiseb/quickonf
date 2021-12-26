package apply

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/program/common/group"
)

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
	c.m.messages <- group.Msg{
		Gidx:  c.m.idx,
		Group: c.m.group,
		Type:  group.ApplyChange,
	}
}

func (c *commandOutput) Infof(format string, a ...interface{}) {
	c.status = commands.StatusInfo
	c.message = fmt.Sprintf(format, a...)
	c.m.messages <- group.Msg{
		Gidx:  c.m.idx,
		Group: c.m.group,
		Type:  group.ApplyChange,
	}
}

func (c *commandOutput) Success(message string) {
	c.status = commands.StatusSuccess
	c.message = message
	c.m.messages <- group.Msg{
		Gidx:  c.m.idx,
		Group: c.m.group,
		Type:  group.ApplyChange,
	}
}

func (c *commandOutput) Successf(format string, a ...interface{}) {
	c.status = commands.StatusSuccess
	c.message = fmt.Sprintf(format, a...)
	c.m.messages <- group.Msg{
		Gidx:  c.m.idx,
		Group: c.m.group,
		Type:  group.ApplyChange,
	}
}

func (c *commandOutput) Error(message string) {
	c.status = commands.StatusError
	c.message = message
	c.m.messages <- group.Msg{
		Gidx:  c.m.idx,
		Group: c.m.group,
		Type:  group.ApplyChange,
	}
}

func (c *commandOutput) Errorf(format string, a ...interface{}) {
	c.status = commands.StatusError
	c.message = fmt.Sprintf(format, a...)
	c.m.messages <- group.Msg{
		Gidx:  c.m.idx,
		Group: c.m.group,
		Type:  group.ApplyChange,
	}
}
