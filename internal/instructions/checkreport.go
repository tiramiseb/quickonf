package instructions

import (
	"fmt"
	"sync"

	"github.com/tiramiseb/quickonf/internal/commands"
)

// CheckReport is a single report after checking is something must be applied
type CheckReport struct {
	Name         string
	status       commands.Status
	message      string
	apply        commands.Apply
	signalTarget chan bool
	Before       string
	After        string

	mu sync.Mutex
}

func (c *CheckReport) Apply() bool {
	if c.apply == nil {
		return true
	}

	return c.apply(c)
}

func (c *CheckReport) GetStatusAndMessage() (commands.Status, string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status, c.message
}

func (c *CheckReport) setStatusAndMessage(status commands.Status, message string) {
	c.mu.Lock()
	c.status = status
	c.message = message
	c.mu.Unlock()
	if c.signalTarget != nil {
		c.signalTarget <- true
	}
}

func (c *CheckReport) Info(message string) {
	c.setStatusAndMessage(commands.StatusInfo, message)
}

func (c *CheckReport) Infof(format string, a ...interface{}) {
	c.setStatusAndMessage(commands.StatusInfo, fmt.Sprintf(format, a...))
}

func (c *CheckReport) Running(message string) {
	c.setStatusAndMessage(commands.StatusRunning, message)
}

func (c *CheckReport) Runningf(format string, a ...interface{}) {
	c.setStatusAndMessage(commands.StatusRunning, fmt.Sprintf(format, a...))
}

func (c *CheckReport) Success(message string) {
	c.setStatusAndMessage(commands.StatusSuccess, message)
}

func (c *CheckReport) Successf(format string, a ...interface{}) {
	c.setStatusAndMessage(commands.StatusSuccess, fmt.Sprintf(format, a...))
}

func (c *CheckReport) Error(message string) {
	c.setStatusAndMessage(commands.StatusError, message)
}

func (c *CheckReport) Errorf(format string, a ...interface{}) {
	c.setStatusAndMessage(commands.StatusError, fmt.Sprintf(format, a...))
}
