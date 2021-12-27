package check

import (
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
)

func (m *model) updateGroupname() {
	if len(m.group.Name) > m.width-3 {
		m.groupName = m.group.Name[:m.width-3]
		return
	}
	if m.group.Priority != 0 {
		priorityTag := fmt.Sprintf(" [!%d] ", m.group.Priority)
		nameAndTag := len(m.group.Name) + len(priorityTag)
		if nameAndTag > m.width-3 {
			m.groupName = m.group.Name
			return
		}
		spaces := m.width - 3 - nameAndTag
		m.groupName = m.group.Name + strings.Repeat(" ", spaces) + priorityTag
		return
	}
	m.groupName = m.group.Name
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
