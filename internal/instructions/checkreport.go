package instructions

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands"
)

// CheckReport is a single report after checking is something must be applied
type CheckReport struct {
	Name         string
	Status       commands.Status
	Message      string
	apply        commands.Apply
	signalTarget chan bool
}

func (c *CheckReport) HasApply() bool {
	return c.apply != nil
}

func (c *CheckReport) Apply() bool {
	if c.apply == nil {
		return true
	}

	return c.apply(c)
}

func (c *CheckReport) Info(message string) {
	c.Status = commands.StatusInfo
	c.Message = message
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}

func (c *CheckReport) Infof(format string, a ...interface{}) {
	c.Status = commands.StatusInfo
	c.Message = fmt.Sprintf(format, a...)
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}

func (c *CheckReport) Running(message string) {
	c.Status = commands.StatusRunning
	c.Message = message
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}

func (c *CheckReport) Runningf(format string, a ...interface{}) {
	c.Status = commands.StatusRunning
	c.Message = fmt.Sprintf(format, a...)
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}

func (c *CheckReport) Success(message string) {
	c.Status = commands.StatusSuccess
	c.Message = message
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}

func (c *CheckReport) Successf(format string, a ...interface{}) {
	c.Status = commands.StatusSuccess
	c.Message = fmt.Sprintf(format, a...)
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}

func (c *CheckReport) Error(message string) {
	c.Status = commands.StatusError
	c.Message = message
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}

func (c *CheckReport) Errorf(format string, a ...interface{}) {
	c.Status = commands.StatusError
	c.Message = fmt.Sprintf(format, a...)
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}
