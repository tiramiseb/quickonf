package output

import (
	"fmt"
	"time"
)

const (
	instructionInfo = iota
	instructionSuccess
	instructionError
)

type Instruction struct {
	group  *Group
	name   string
	status int
	output string

	instructionPrefix string
	spinnerPos        int
	ticker            *time.Ticker
	tickStopper       chan bool
}

func (i *Instruction) stopTicker() {
	if i.ticker == nil {
		return
	}
	i.ticker.Stop()
	i.tickStopper <- true
}

func (i *Instruction) draw() string {
	name := i.name
	output := i.output
	totalLen := len(name) + 4 + len(output) // 4 = prefix + space + ":" + space
	if totalLen > width {
		removeFromEach := (totalLen - width) / 2
		name = name[:len(name)-removeFromEach]
		if len(name) > 6 {
			name = name[:len(name)-3] + "..."
		}
		output = output[:len(output)-removeFromEach]
		if len(output) > 6 {
			output = output[:len(output)-3] + "..."
		}
	}
	var prefix string
	var fg string
	switch i.status {
	case instructionInfo:
		prefix = prefixRunning
		fg = fgBlue
	case instructionSuccess:
		prefix = prefixSuccess
		fg = fgGreen
	case instructionError:
		prefix = prefixError
		fg = fgRed
	}
	return fmt.Sprintf("%s%s%s %s%s:%s %s", fg, bold, prefix, i.instructionPrefix, name, reset, output)
}

// func (i *InstructionOutput) Load(text string) {
// 	// Not  using the pterm spinner because it can't be put in an area
// 	ticker := time.NewTicker(500 * time.Millisecond)
// 	spinPos := 0
// 	go func() {
// 		for {
// 			select {
// 			case <-i.stopSpin:
// 				return
// 			case <-ticker.C:
// 				spinPos++
// 				if spinPos > 3 {
// 					spinPos = 0
// 				}
// 				i.content = pterm.Info.Sprintf("%s %s\n", customInstructionSpinner[spinPos], text)
// 				i.step.update()
// 			}
// 		}
// 	}()
// }

// func (i *InstructionOutput) StopLoad() {
// 	i.stopSpin <- true
// }
