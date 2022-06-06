package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/internal/program/messages"
)

func (m *model) listenSignal() tea.Msg {
	<-m.signalTarget
	return messages.NewSignal{}
}
