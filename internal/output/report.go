package output

import (
	"strings"
)

// report returns a report to show
func report() string {
	b := strings.Builder{}
	for _, group := range groups {
		groupDisplayed := false
		for _, instr := range group.instructions {
			if instr.status == instructionError {
				if !groupDisplayed {
					b.WriteString(bgRed)
					b.WriteString("  ")
					b.WriteString(group.name)
					b.WriteString("  \n")
					b.WriteString(reset)
					groupDisplayed = true
				}
				b.WriteString(fgRed)
				b.WriteString("  ")
				b.WriteString(instr.name)
				b.WriteString(":")
				b.WriteString(reset)
				b.WriteString(" ")
				b.WriteString(instr.output)
				b.WriteString("\n")
			}
		}
	}
	return b.String()
}
