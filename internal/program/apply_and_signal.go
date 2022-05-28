package program

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) listenSignal() tea.Msg {
	<-m.signalTarget
	return newSignal{}
}

type newSignal struct{}
