package apply

import (
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/program/style"
)

func (m *model) updateGroupname() {
	if len(m.group.Name) > m.width-1 {
		m.groupName = m.group.Name[:m.width-1]
	} else {
		m.groupName = m.group.Name
	}
}

func (m *model) instructionLine(name string, status commands.Status, message string) string {
	prefix := "  " + name + "  "
	prefixSize := len(prefix)
	if prefixSize >= m.width {
		prefix = prefix[:m.width]
		return InstructionStyles[status].Render(prefix)
	}
	prefix = InstructionStyles[status].Render(prefix)
	message = "  " + message
	messageSize := len(message)
	widthForMessage := m.width - prefixSize
	if messageSize > widthForMessage {
		message = message[:widthForMessage]
	}
	return prefix + style.InstructionMessage.Render(message)
}
