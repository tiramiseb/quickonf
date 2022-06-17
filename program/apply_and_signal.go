package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/program/messages"
)

func (m *model) listenSignal() tea.Msg {
	<-m.signalTarget
	return messages.NewSignal{}
}
