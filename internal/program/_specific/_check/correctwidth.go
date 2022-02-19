package check

import (
	"fmt"
	"strings"
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
