package output

func (i *Instruction) Info(text string) {
	i.stopTicker()
	i.output = text
	i.status = instructionInfo
	i.group.preRedraw()
	redraw()
}

// func (i *Instruction) Infof(format string, a ...interface{}) {
// 	i.stopTicker()
// 	i.output = fmt.Sprintf(format, a...)
// 	i.status = instructionInfo
// 	i.group.preRedraw()
// 	redraw()
// }
