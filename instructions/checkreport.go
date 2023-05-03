package instructions

import (
	"fmt"
	"sync"
	"time"

	"github.com/tiramiseb/quickonf/commands"
)

var (
	spinnerSteps = [8]string{
		"▄ ",
		"█ ",
		"▀ ",
		"▀▀",
		" ▀",
		" █",
		" ▄",
		"▄▄",
	}
	nbSpinnerStep     = len(spinnerSteps)
	percentIndicators = [101]string{
		"     ",
		"▏    ", "▏    ", "▎    ", "▎    ", "▍    ",
		"▍    ", "▍    ", "▌    ", "▌    ", "▌    ",
		"▋    ", "▋    ", "▊    ", "▊    ", "▊    ",
		"▉    ", "▉    ", "█    ", "█    ", "█    ",

		"█▏   ", "█▏   ", "█▎   ", "█▎   ", "█▍   ",
		"█▍   ", "█▍   ", "█▌   ", "█▌   ", "█▌   ",
		"█▋   ", "█▋   ", "█▊   ", "█▊   ", "█▊   ",
		"█▉   ", "█▉   ", "██   ", "██   ", "██   ",

		"██▏  ", "██▏  ", "██▎  ", "██▎  ", "██▍  ",
		"██▍  ", "██▍  ", "██▌  ", "██▌  ", "██▌  ",
		"██▋  ", "██▋  ", "██▊  ", "██▊  ", "██▊  ",
		"██▉  ", "██▉  ", "███  ", "███  ", "███  ",

		"███▏ ", "███▏ ", "███▎ ", "███▎ ", "███▍ ",
		"███▍ ", "███▍ ", "███▌ ", "███▌ ", "███▌ ",
		"███▋ ", "███▋ ", "███▊ ", "███▊ ", "███▊ ",
		"███▉ ", "███▉ ", "████ ", "████ ", "████ ",

		"████▏", "████▏", "████▎", "████▎", "████▍",
		"████▍", "████▍", "████▌", "████▌", "████▌",
		"████▋", "████▋", "████▊", "████▊", "████▊",
		"████▉", "████▉", "█████", "█████", "█████",
	}
)

// CheckReport is a single report after checking is something must be applied
type CheckReport struct {
	Name         string
	level        int
	status       commands.Status
	message      string
	apply        commands.Apply
	signalTarget chan bool
	Before       string
	After        string

	runningStopper chan bool

	mu sync.Mutex
}

func (c *CheckReport) Apply() bool {
	if c.apply == nil {
		return true
	}

	return c.apply(c)
}

func (c *CheckReport) GetStatusAndMessage() (commands.Status, string, int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.status, c.message, c.level
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
	c.stopSpinner()
	c.setStatusAndMessage(commands.StatusInfo, message)
}

func (c *CheckReport) Infof(format string, a ...interface{}) {
	c.stopSpinner()
	c.setStatusAndMessage(commands.StatusInfo, fmt.Sprintf(format, a...))
}

func (c *CheckReport) Running(message string) {
	c.stopSpinner()
	c.spinRunning(message)
}

func (c *CheckReport) Runningf(format string, a ...interface{}) {
	c.stopSpinner()
	c.spinRunning(fmt.Sprintf(format, a...))
}

func (c *CheckReport) RunningPercent(message string) chan int {
	c.stopSpinner()
	return c.percentRunning(message)
}

func (c *CheckReport) RunningPercentf(format string, a ...interface{}) chan int {
	c.stopSpinner()
	return c.percentRunning(fmt.Sprintf(format, a...))
}

func (c *CheckReport) Success(message string) {
	c.stopSpinner()
	c.setStatusAndMessage(commands.StatusSuccess, message)
}

func (c *CheckReport) Successf(format string, a ...interface{}) {
	c.stopSpinner()
	c.setStatusAndMessage(commands.StatusSuccess, fmt.Sprintf(format, a...))
}

func (c *CheckReport) Error(message string) {
	c.stopSpinner()
	c.setStatusAndMessage(commands.StatusError, message)
}

func (c *CheckReport) Errorf(format string, a ...interface{}) {
	c.stopSpinner()
	c.setStatusAndMessage(commands.StatusError, fmt.Sprintf(format, a...))
}

func (c *CheckReport) spinRunning(message string) {
	c.runningStopper = make(chan bool)
	ticker := time.NewTicker(500 * time.Millisecond)
	var step int
	go func() {
		for {
			select {
			case <-c.runningStopper:
				ticker.Stop()
				return
			case <-ticker.C:
				c.setStatusAndMessage(commands.StatusRunning, fmt.Sprintf("▕%s▏%s", spinnerSteps[step], message))
				step++
				if step == nbSpinnerStep {
					step = 0
				}
			}
		}
	}()
}

func (c *CheckReport) percentRunning(message string) chan int {
	c.runningStopper = make(chan bool)
	percent := make(chan int)
	c.setStatusAndMessage(commands.StatusRunning, fmt.Sprintf("▕%s▏%s", percentIndicators[0], message))
	go func() {
		for {
			select {
			case <-c.runningStopper:
				close(percent)
				return
			case val := <-percent:
				if val < 0 {
					val = 0
				} else if val > 100 {
					val = 100
				}
				c.setStatusAndMessage(commands.StatusRunning, fmt.Sprintf("▕%s▏%s", percentIndicators[val], message))
			}
		}
	}()
	return percent
}

func (c *CheckReport) stopSpinner() {
	if c.runningStopper != nil {
		c.runningStopper <- true
		c.runningStopper = nil
	}
}
