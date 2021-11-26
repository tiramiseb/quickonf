package output

import "fmt"

func (i *Instruction) Error(text string) {
	i.stopTicker()
	i.output = text
	i.status = instructionError
	i.group.preRedraw()
	redraw()
}

func (i *Instruction) Errorf(format string, a ...interface{}) {
	i.stopTicker()
	i.output = fmt.Sprintf(format, a...)
	i.status = instructionError
	i.group.preRedraw()
	redraw()
}
