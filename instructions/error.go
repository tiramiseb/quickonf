package instructions

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands"
)

// Error is used when there is an error in the configuration file
type Error struct {
	Line    int
	Column  int
	Message string
}

func (e *Error) Name() string {
	return "error"
}

func (e *Error) RunCheck(vars *Variables, signalTarget chan bool, level int) ([]*CheckReport, bool) {
	return nil, false
}

func (e *Error) NotRunReports(level int) []*CheckReport {
	msg := e.description()
	return []*CheckReport{{
		Name:    "ERROR",
		level:   level,
		status:  commands.StatusError,
		message: msg.string(0),
	}}
}
func (e *Error) Reset() {}

func (e *Error) String() string {
	return e.indentedString(0)
}

func (e *Error) indentedString(level int) string {
	content := e.description()
	return content.string(level)
}
func (e *Error) description() stringBuilder {
	var content stringBuilder
	content.add(fmt.Sprintf("[%d:%d]", e.Line, e.Column))
	content.add(e.Message)
	return content
}

func (e *Error) hasConfigError() bool {
	return true
}
