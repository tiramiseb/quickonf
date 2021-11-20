package output

import "fmt"

// func (i *Instruction) Success(text string) {
// 	i.stopTicker()
// 	i.output = text
// 	i.status = instructionSuccess
// 	i.group.preRedraw()
// 	redraw()
// }

func (i *Instruction) Successf(format string, a ...interface{}) {
	i.stopTicker()
	i.output = fmt.Sprintf(format, a...)
	i.status = instructionSuccess
	i.group.preRedraw()
	redraw()
}
