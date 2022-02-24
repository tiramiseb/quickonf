package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func apply(i int) tea.Cmd {
	return func() tea.Msg {
		go global.DisplayedGroups[i].Apply()
		return nil
	}
}

func (m *model) listenSignal() tea.Msg {
	<-m.signalTarget
	return newSignal{}
}

type newSignal struct{}
