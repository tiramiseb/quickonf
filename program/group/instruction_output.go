package group

import (
	"fmt"

	"github.com/tiramiseb/quickonf/state"
)

type instructionStatus int

const (
	instructionInfo instructionStatus = iota
	instructionSuccess
	instructionError
)

type instructionOutput struct {
	m       *Model
	name    string
	Status  instructionStatus
	message string
}

func newInstructionOutput(model *Model, name string) *instructionOutput {
	out := &instructionOutput{m: model, name: name, message: "running..."}
	model.outputs = append(model.outputs, out)
	return out
}

func (i *instructionOutput) NewLine(name string) state.Output {
	return newInstructionOutput(i.m, name)
}

func (i *instructionOutput) Info(message string) {
	i.Status = instructionInfo
	i.message = message
	i.m.messages <- ChangeMessage{i.m.idx}
}

func (i *instructionOutput) Infof(format string, a ...interface{}) {
	i.Status = instructionInfo
	i.message = fmt.Sprintf(format, a...)
	i.m.messages <- ChangeMessage{i.m.idx}
}

func (i *instructionOutput) Success(message string) {
	i.Status = instructionSuccess
	i.message = message
	i.m.messages <- ChangeMessage{i.m.idx}
}

func (i *instructionOutput) Successf(format string, a ...interface{}) {
	i.Status = instructionSuccess
	i.message = fmt.Sprintf(format, a...)
	i.m.messages <- ChangeMessage{i.m.idx}
}

func (i *instructionOutput) Error(message string) {
	i.Status = instructionError
	i.message = message
	i.m.messages <- ChangeMessage{i.m.idx}
}

func (i *instructionOutput) Errorf(format string, a ...interface{}) {
	i.Status = instructionError
	i.message = fmt.Sprintf(format, a...)
	i.m.messages <- ChangeMessage{i.m.idx}
}
