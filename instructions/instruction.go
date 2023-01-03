package instructions

import "strings"

// Instruction is a single instruction
type Instruction interface {
	// Name returns the instruction name
	Name() string
	// RunCheck runs the instruction's check part and return the check reports
	RunCheck(vars *Variables, signalTarget chan bool, level int) ([]*CheckReport, bool)
	// Reset everything, to have it as it has never run
	Reset()

	// Return a textual representation of the instruction
	String() string
	indentedString(level int) string
}

type stringBuilder struct {
	content []string
}

func (s *stringBuilder) add(elem string) {
	elem = strings.ReplaceAll(elem, "\"", "\\\"")
	if strings.Contains(elem, " ") || strings.Contains(elem, "\n") || strings.Contains(elem, "\t") {
		elem = "\"" + elem + "\""
	}
	s.content = append(s.content, elem)
}

func (s *stringBuilder) string(level int) string {
	return strings.Repeat("  ", level) + strings.Join(s.content, " ")
}
