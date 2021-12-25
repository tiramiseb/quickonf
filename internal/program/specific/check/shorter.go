package check

import (
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
)

func (m *model) updateGroupname() {
	if len(m.group.Name) > m.width-3 {
		m.groupName = m.group.Name[:m.width-3]
	} else {
		m.groupName = m.group.Name
	}
}

func (m *model) instructionLine(report instructions.CheckReport) string {
	prefix := "  " + report.Name + "  "
	prefixSize := len(prefix)
	if prefixSize >= m.width {
		prefix = prefix[:m.width]
		return InstructionStyles[report.Status].Render(prefix)
	}
	prefix = InstructionStyles[report.Status].Render(prefix)
	message := "  " + report.Message
	messageSize := len(message)
	widthForMessage := m.width - prefixSize
	if messageSize > widthForMessage {
		message = message[:widthForMessage]
	}
	return prefix + style.InstructionMessage.Render(message)
}
