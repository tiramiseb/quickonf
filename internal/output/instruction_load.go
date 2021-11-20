package output

import (
	"fmt"
	"time"
)

var loadSpinner = []string{"|▀ | ", "| ▀| ", "| ▄| ", "|▄ | "}

func (i *Instruction) Load(text string) {
	i.stopTicker()
	i.output = text
	i.status = instructionInfo
	i.ShowLoadSpinner()
	i.group.preRedraw()
	redraw()
}

func (i *Instruction) Loadf(format string, a ...interface{}) {
	i.stopTicker()
	i.output = fmt.Sprintf(format, a...)
	i.status = instructionInfo
	i.ShowLoadSpinner()
	i.group.preRedraw()
	redraw()
}

func (i *Instruction) ShowLoadSpinner() {
	i.spinnerPos = 0
	i.instructionPrefix = loadSpinner[i.spinnerPos]
	i.ticker = time.NewTicker(200 * time.Millisecond)
	i.tickStopper = make(chan bool)
	go func() {
		for {
			select {
			case <-i.tickStopper:
				i.instructionPrefix = ""
				return
			case <-i.ticker.C:
				i.spinnerPos++
				if i.spinnerPos >= len(loadSpinner) {
					i.spinnerPos = 0
				}
				i.instructionPrefix = loadSpinner[i.spinnerPos]
				i.group.preRedraw()
				redraw()
			}
		}
	}()
}
